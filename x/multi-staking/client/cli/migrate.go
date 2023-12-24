package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	v0 "github.com/realio-tech/multi-staking-module/x/multi-staking/types/v0"

	"github.com/spf13/cobra"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"
)

type AppMap map[string]json.RawMessage

const flagGenesisTime = "genesis-time"

var (
	prefix                 = "realio"
	newBondedTokenDenom    = "stake"
	bondedPoolAddress, _   = sdk.Bech32ifyAddressBytes(prefix, authtypes.NewModuleAddress(stakingtypes.BondedPoolName))
	unbondedPoolAddress, _ = sdk.Bech32ifyAddressBytes(prefix, authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName))
)

// MigrateGenesisCmd returns a command to execute genesis state migration.
func MigrateStakingGenesisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-staking [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			var err error
			// Read genesis state
			genesisPath := args[0]
			genDoc, err := validateGenDoc(genesisPath)
			if err != nil {
				return err
			}

			var initialState AppMap
			if err := json.Unmarshal(genDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			// migrate
			newGenState, err := migrate(initialState, clientCtx)
			if err != nil {
				return err
			}

			genDoc.AppState, err = json.Marshal(newGenState)
			if err != nil {
				return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
			}

			// add genesis time
			genesisTime, _ := cmd.Flags().GetString(flagGenesisTime)
			if genesisTime != "" {
				var t time.Time

				err := t.UnmarshalText([]byte(genesisTime))
				if err != nil {
					return errors.Wrap(err, "failed to unmarshal genesis time")
				}

				genDoc.GenesisTime = t
			}

			// add chain-id
			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			if chainID != "" {
				genDoc.ChainID = chainID
			}

			bz, err := tmjson.Marshal(genDoc)
			if err != nil {
				return errors.Wrap(err, "failed to marshal genesis doc")
			}

			sortedBz, err := sdk.SortJSON(bz)
			if err != nil {
				return errors.Wrap(err, "failed to sort JSON genesis doc")
			}

			cmd.Println(string(sortedBz))

			// save new genesis
			genDoc.SaveAs(genesisPath)

			return nil
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "override genesis_time with this flag")
	cmd.Flags().String(flags.FlagChainID, "", "override chain_id with this flag")

	return cmd
}

// validateGenDoc reads a genesis file and validates that it is a correct
// Tendermint GenesisDoc. This function does not do any cosmos-related
// validation.
func validateGenDoc(importGenesisFile string) (*tmtypes.GenesisDoc, error) {
	genDoc, err := tmtypes.GenesisDocFromFile(importGenesisFile)
	if err != nil {
		return nil, fmt.Errorf("%s. Make sure that"+
			" you have correctly migrated all Tendermint consensus params",
			err.Error(),
		)
	}

	return genDoc, nil
}

func migrate(genesisState AppMap, ctx client.Context) (AppMap, error) {
	newGenState, err := migrateBank(genesisState)
	if err != nil {
		return nil, err
	}
	newGenState, err = migrateStaking(newGenState)
	if err != nil {
		return nil, err
	}
	newGenState, err = migrateDistribution(newGenState)
	if err != nil {
		return nil, err
	}
	return newGenState, nil
}

func migrateBank(genesisState AppMap) (AppMap, error) {
	rawData := genesisState[stakingtypes.ModuleName]
	var oldStakingState v0.GenesisState
	err := json.Unmarshal(rawData, &oldStakingState)
	if err != nil {
		return nil, err
	}

	rawData = genesisState[banktypes.ModuleName]
	var oldBankState banktypes.GenesisState
	err = json.Unmarshal(rawData, &oldBankState)
	if err != nil {
		return nil, err
	}

	newbalances, newSupply, err := convertBankState(oldBankState.Balances, oldBankState.Supply, oldStakingState)
	if err != nil {
		return nil, err
	}

	// new bank genesis
	var newBankState banktypes.GenesisState

	newBankState.Params = oldBankState.Params
	newBankState.Balances = newbalances
	newBankState.Supply = newSupply
	newBankState.DenomMetadata = oldBankState.DenomMetadata

	// replace to genesis state
	newData, err := json.Marshal(&newBankState)
	if err != nil {
		return nil, err
	}

	genesisState[banktypes.ModuleName] = newData

	return genesisState, nil
}

func convertBankState(
	oldBalances []banktypes.Balance,
	oldSupply sdk.Coins,
	oldStakingState v0.GenesisState,
) ([]banktypes.Balance, sdk.Coins, error) {
	// validator map
	validatorMap := make(map[string]v0.Validator, 0)
	for _, validator := range oldStakingState.Validators {
		validatorMap[validator.OperatorAddress] = validator
	}

	var newBalances []banktypes.Balance
	var bondedPoolBalance sdk.Coins
	var ubdPoolBalance sdk.Coins
	var totalNewBondedTokenAmount sdk.Coins
	var balancesIndexMap = make(map[string]uint64, 0)

	for _, delegation := range oldStakingState.Delegations {
		var delegationLockAmount sdk.Coins

		// change the delegation from staking bonded(unbonded) pool to intermedary account
		DelAddr := sdk.AccAddress(delegation.DelegatorAddress)
		intermediaryAccount := types.IntermediaryAccount(DelAddr)
		intermediaryBech32Addr, _ := sdk.Bech32ifyAddressBytes(prefix, intermediaryAccount)

		validator, ok := validatorMap[delegation.ValidatorAddress]
		if !ok {
			return nil, nil, fmt.Errorf("Error validator not found delegation %s", delegation.ValidatorAddress)
		}

		// calculate issued token
		val, tokenAmount := tokenAmountFromShares(validator, delegation.Shares)
		delegationAmount := sdk.NewCoins(sdk.NewCoin(validator.BondDenom, tokenAmount))

		delegationLockAmount = delegationLockAmount.Add(delegationAmount...)

		// move delegate token
		if validator.Status == stakingtypes.BondStatusBonded {
			// Caculate new bonded token amount in bondedPoolAddress
			bondedPoolBalance = bondedPoolBalance.Add(sdk.NewCoin(newBondedTokenDenom, tokenAmount)) // TODO: need to add ratio. Current ratio is 1
		} else {
			// Caculate new bonded token amount in unbondedPoolAddress
			ubdPoolBalance = ubdPoolBalance.Add(sdk.NewCoin(newBondedTokenDenom, tokenAmount)) // TODO: need to add ratio. Current ratio is 1
		}
		totalNewBondedTokenAmount = totalNewBondedTokenAmount.Add(sdk.NewCoin(newBondedTokenDenom, tokenAmount)) // TODO: need to add ratio. Current ratio is 1

		// update validator
		validatorMap[delegation.ValidatorAddress] = val

		index, found := balancesIndexMap[intermediaryBech32Addr]
		if !found {
			index := len(newBalances)
			balancesIndexMap[intermediaryBech32Addr] = uint64(index)

			balance := banktypes.Balance{
				Address: intermediaryBech32Addr,
				Coins:   delegationLockAmount,
			}

			newBalances = append(newBalances, balance)
		} else {
			delegationLockAmount = delegationLockAmount.Add(newBalances[index].Coins...)

			balance := banktypes.Balance{
				Address: intermediaryBech32Addr,
				Coins:   delegationLockAmount,
			}

			newBalances[index] = balance
		}

	}

	for _, ubdDelegation := range oldStakingState.UnbondingDelegations {
		var delegationStakeAmount sdk.Coins
		// change the delegation from staking unbonded pool to intermedary account
		DelAddr := sdk.AccAddress(ubdDelegation.DelegatorAddress)
		intermediaryAccount := types.IntermediaryAccount(DelAddr)
		intermediaryBech32Addr, _ := sdk.Bech32ifyAddressBytes(prefix, intermediaryAccount)

		for _, entry := range ubdDelegation.Entries {
			// move balance in entry from unbondedPoolAddress to intermediary address
			delegationStakeAmount = delegationStakeAmount.Add(entry.Balance)
			// Caculate new bonded token amount in unbondedPoolAddress
			ubdPoolBalance = ubdPoolBalance.Add(sdk.NewCoin(newBondedTokenDenom, entry.Balance.Amount))                       // TODO: need to add ratio. Current ratio is 1
			totalNewBondedTokenAmount = totalNewBondedTokenAmount.Add(sdk.NewCoin(newBondedTokenDenom, entry.Balance.Amount)) // TODO: need to add ratio. Current ratio is 1
		}

		index, found := balancesIndexMap[intermediaryBech32Addr]
		if !found {
			index := len(newBalances)
			balancesIndexMap[intermediaryBech32Addr] = uint64(index)

			balance := banktypes.Balance{
				Address: intermediaryBech32Addr,
				Coins:   delegationStakeAmount,
			}

			newBalances = append(newBalances, balance)
		} else {
			delegationStakeAmount = delegationStakeAmount.Add(newBalances[index].Coins...)

			balance := banktypes.Balance{
				Address: intermediaryBech32Addr,
				Coins:   delegationStakeAmount,
			}

			newBalances[index] = balance
		}
	}

	// append with oldBalances
	for _, balance := range oldBalances {
		// assign new bondedPool and ubdPool balances
		if balance.Address == bondedPoolAddress {
			balance.Coins = bondedPoolBalance
		}

		if balance.Address == unbondedPoolAddress {
			balance.Coins = ubdPoolBalance
		}

		newBalances = append(newBalances, balance)
	}

	// new supply
	newSupply := oldSupply.Add(totalNewBondedTokenAmount...)

	return newBalances, newSupply, nil
}

func tokenAmountFromShares(v v0.Validator, delShares sdk.Dec) (v0.Validator, math.Int) {
	remainingShares := v.DelegatorShares.Sub(delShares)

	var amount math.Int
	if remainingShares.IsZero() {
		// last delegation share gets any trimmings
		amount = v.Tokens
		v.Tokens = math.ZeroInt()
	} else {
		// leave excess tokens in the validator
		// however fully use all the delegator shares
		amount = tokensFromShares(v, delShares).TruncateInt()
		v.Tokens = v.Tokens.Sub(amount)

		if v.Tokens.IsNegative() {
			panic("attempting to remove more tokens than available in validator")
		}
	}

	v.DelegatorShares = remainingShares

	return v, amount
}

func tokensFromShares(v v0.Validator, shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).Quo(v.DelegatorShares)
}

func migrateStaking(genesisState AppMap) (AppMap, error) {

	rawData := genesisState[stakingtypes.ModuleName]
	var oldState v0.GenesisState
	err := json.Unmarshal(rawData, &oldState)
	fmt.Println("old state", oldState)
	if err != nil {
		return nil, err
	}

	newState := types.GenesisState{}
	// Migrate state.StakingGenesisState
	stakingGenesisState := stakingtypes.GenesisState{}

	unbondingTime, err := time.ParseDuration(oldState.Params.UnbondingTime)
	if err != nil {
		return nil, err
	}
	stakingGenesisState.Params = stakingtypes.Params{
		UnbondingTime:     unbondingTime,
		MaxValidators:     oldState.Params.MaxValidators,
		MaxEntries:        oldState.Params.MaxEntries,
		HistoricalEntries: oldState.Params.HistoricalEntries,
		BondDenom:         newBondedTokenDenom,
		MinCommissionRate: oldState.Params.MinCommissionRate,
	}
	stakingGenesisState.LastTotalPower = oldState.LastTotalPower
	stakingGenesisState.Validators = convertValidators(oldState.Validators)
	stakingGenesisState.Delegations = convertDelegations(oldState.Delegations)
	stakingGenesisState.UnbondingDelegations = convertUnbondings(oldState.UnbondingDelegations)
	stakingGenesisState.Redelegations = convertRedelegations(oldState.Redelegations)
	stakingGenesisState.Exported = oldState.Exported

	newState.StakingGenesisState = &stakingGenesisState
	// Migrate state.ValidatorAllowedToken field
	newState.ValidatorAllowedToken = make([]types.ValidatorAllowedToken, 0)

	for _, val := range oldState.Validators {
		allowedToken := types.ValidatorAllowedToken{
			ValAddr:    val.OperatorAddress,
			TokenDenom: val.BondDenom,
		}
		newState.ValidatorAllowedToken = append(newState.ValidatorAllowedToken, allowedToken)
	}
	// Migrate state.MultiStakingLock field
	newState.MultiStakingLocks = make([]types.MultiStakingLock, 0)

	for _, val := range stakingGenesisState.Validators {
		for _, del := range oldState.Delegations {
			if del.ValidatorAddress == val.OperatorAddress {
				lock := types.MultiStakingLock{
					ConversionRatio: sdk.OneDec(),
					DelAddr:         del.DelegatorAddress,
					ValAddr:         del.ValidatorAddress,
					LockedAmount:    val.TokensFromShares(del.Shares).TruncateInt(),
				}
				newState.MultiStakingLocks = append(newState.MultiStakingLocks, lock)
			}

		}
	}

	err = newState.Validate()
	if err != nil {
		return nil, err
	}

	newData, err := json.Marshal(&newState)
	if err != nil {
		return nil, err
	}

	genesisState[types.ModuleName] = newData

	return genesisState, nil
}

func convertValidators(validators []v0.Validator) []stakingtypes.Validator {
	newValidators := make([]stakingtypes.Validator, 0)
	for _, val := range validators {
		newVal := stakingtypes.Validator{
			OperatorAddress: val.OperatorAddress,
			ConsensusPubkey: val.ConsensusPubkey,
			Jailed:          val.Jailed,
			Status:          stakingtypes.BondStatus(stakingtypes.BondStatus_value[val.Status]),
			Tokens:          val.Tokens,
			DelegatorShares: val.DelegatorShares,
			Description:     stakingtypes.Description(val.Description),
			UnbondingHeight: val.UnbondingHeight,
			UnbondingTime:   val.UnbondingTime,
			Commission: stakingtypes.Commission{
				CommissionRates: stakingtypes.CommissionRates(val.Commission.CommissionRates),
				UpdateTime:      val.Commission.UpdateTime,
			},
			MinSelfDelegation: val.MinSelfDelegation,
		}
		newValidators = append(newValidators, newVal)
	}
	return newValidators
}

func convertDelegations(delegations []v0.Delegation) []stakingtypes.Delegation {
	newDelegations := make([]stakingtypes.Delegation, 0)
	for _, del := range delegations {
		newDel := stakingtypes.Delegation(del)
		newDelegations = append(newDelegations, newDel)
	}
	return newDelegations
}

func convertUnbondings(ubds []v0.UnbondingDelegation) []stakingtypes.UnbondingDelegation {
	newUbds := make([]stakingtypes.UnbondingDelegation, 0)
	for _, ubd := range ubds {
		newEntries := make([]stakingtypes.UnbondingDelegationEntry, 0)
		for _, entry := range ubd.Entries {
			newEntry := stakingtypes.UnbondingDelegationEntry{
				CreationHeight: entry.CreationHeight,
				CompletionTime: entry.CompletionTime,
				InitialBalance: entry.InitialBalance.Amount,
				Balance:        entry.Balance.Amount,
			}
			newEntries = append(newEntries, newEntry)
		}
		newUbd := stakingtypes.UnbondingDelegation{
			DelegatorAddress: ubd.DelegatorAddress,
			ValidatorAddress: ubd.ValidatorAddress,
			Entries:          newEntries,
		}
		newUbds = append(newUbds, newUbd)
	}
	return newUbds
}

func convertRedelegations(reDels []v0.Redelegation) []stakingtypes.Redelegation {
	newRedels := make([]stakingtypes.Redelegation, 0)
	for _, reDel := range reDels {
		newEntries := make([]stakingtypes.RedelegationEntry, 0)
		for _, entry := range reDel.Entries {
			newEntry := stakingtypes.RedelegationEntry{
				CreationHeight: entry.CreationHeight,
				CompletionTime: entry.CompletionTime,
				InitialBalance: entry.InitialBalance.Amount,
				SharesDst:      entry.SharesDst,
			}
			newEntries = append(newEntries, newEntry)
		}
		newRedel := stakingtypes.Redelegation{
			DelegatorAddress:    reDel.DelegatorAddress,
			ValidatorSrcAddress: reDel.ValidatorSrcAddress,
			ValidatorDstAddress: reDel.ValidatorDstAddress,
			Entries:             newEntries,
		}
		newRedels = append(newRedels, newRedel)
	}
	return newRedels
}

func migrateDistribution(genesisState AppMap) (AppMap, error) {
	rawData := genesisState[distrtypes.ModuleName]
	var oldState v0.DistrGenesisState

	err := json.Unmarshal(rawData, &oldState)
	if err != nil {
		return nil, err
	}

	// Migrate state.DistributionGenesisState
	newDelegatorStartingInfos := make([]distrtypes.DelegatorStartingInfoRecord, 0)

	fmt.Println(len(oldState.DelegatorStartingInfos))

	for _, info := range oldState.DelegatorStartingInfos {
		delAddr := sdk.AccAddress(info.DelegatorAddress)

		oldPeriod, _ := strconv.ParseUint(info.StartingInfo.PreviousPeriod, 10, 64)
		oldHeight, _ := strconv.ParseUint(info.StartingInfo.Height, 10, 64)

		startingInfo := distrtypes.DelegatorStartingInfo{
			PreviousPeriod: oldPeriod,
			Height:         oldHeight,
			Stake:          info.StartingInfo.Stake,
		}
		intermediaryAccount := types.IntermediaryAccount(delAddr)
		fmt.Println(delAddr.String())
		fmt.Println(intermediaryAccount.String())
		newRecord := distrtypes.DelegatorStartingInfoRecord{
			DelegatorAddress: intermediaryAccount.String(),
			ValidatorAddress: info.ValidatorAddress,
			StartingInfo:     startingInfo,
		}
		newDelegatorStartingInfos = append(newDelegatorStartingInfos, newRecord)
	}

	newDelegatorWithdrawInfos := make([]distrtypes.DelegatorWithdrawInfo, 0)
	for _, info := range oldState.DelegatorWithdrawInfos {
		delAddr := sdk.AccAddress(info.DelegatorAddress)
		intermediaryAccount := types.IntermediaryAccount(delAddr)
		newRecord := distrtypes.DelegatorWithdrawInfo{
			DelegatorAddress: intermediaryAccount.String(),
			WithdrawAddress:  info.WithdrawAddress,
		}
		newDelegatorWithdrawInfos = append(newDelegatorWithdrawInfos, newRecord)
	}

	newValHistoryRewards := make([]distrtypes.ValidatorHistoricalRewardsRecord, 0)
	for _, info := range oldState.ValidatorHistoricalRewards {
		oldPeriod, _ := strconv.ParseUint(info.Period, 10, 64)
		newRecord := distrtypes.ValidatorHistoricalRewardsRecord{
			ValidatorAddress: info.ValidatorAddress,
			Rewards:          info.Rewards,
			Period:           oldPeriod,
		}
		newValHistoryRewards = append(newValHistoryRewards, newRecord)
	}

	newValCurrentRewards := make([]distrtypes.ValidatorCurrentRewardsRecord, 0)
	for _, info := range oldState.ValidatorCurrentRewards {
		oldPeriod, _ := strconv.ParseUint(info.Rewards.Period, 10, 64)
		newRecord := distrtypes.ValidatorCurrentRewardsRecord{
			ValidatorAddress: info.ValidatorAddress,
			Rewards: distrtypes.ValidatorCurrentRewards{
				Period:  oldPeriod,
				Rewards: info.Rewards.Rewards,
			},
		}
		newValCurrentRewards = append(newValCurrentRewards, newRecord)
	}

	newValSlashEvents := make([]distrtypes.ValidatorSlashEventRecord, 0)
	for _, info := range oldState.ValidatorSlashEvents {
		oldPeriod, _ := strconv.ParseUint(info.Period, 10, 64)
		oldHeight, _ := strconv.ParseUint(info.Height, 10, 64)
		oldValPeriod, _ := strconv.ParseUint(info.ValidatorSlashEvent.ValidatorPeriod, 10, 64)

		newRecord := distrtypes.ValidatorSlashEventRecord{
			ValidatorAddress: info.ValidatorAddress,
			Height: oldHeight,
			Period: oldPeriod,
			ValidatorSlashEvent: distrtypes.NewValidatorSlashEvent(oldValPeriod, info.ValidatorSlashEvent.Fraction),
		}
		newValSlashEvents = append(newValSlashEvents, newRecord)
	}

	newState := distrtypes.GenesisState{
		Params:                          oldState.Params,
		FeePool:                         oldState.FeePool,
		DelegatorWithdrawInfos:          newDelegatorWithdrawInfos,
		PreviousProposer:                oldState.PreviousProposer,
		OutstandingRewards:              oldState.OutstandingRewards,
		ValidatorAccumulatedCommissions: oldState.ValidatorAccumulatedCommissions,
		ValidatorHistoricalRewards:      newValHistoryRewards,
		ValidatorCurrentRewards:         newValCurrentRewards,
		DelegatorStartingInfos:          newDelegatorStartingInfos,
		ValidatorSlashEvents:            newValSlashEvents,
	}

	newData, err := json.Marshal(&newState)
	if err != nil {
		return nil, err
	}

	genesisState[types.ModuleName] = newData

	return genesisState, nil
}
