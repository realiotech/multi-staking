package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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
	bondedPoolAddress   = authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String()
	unbondedPoolAddress = authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName).String()
	newBondedTokenDenom = "stake"
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
	return genesisState, nil
}

func migrateBank(genesisState AppMap, ctx client.Context) (AppMap, error) {
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

	// validator map
	validatorMap := make(map[string]v0.Validator, 0)
	for _, validator := range oldStakingState.Validators {
		validatorMap[validator.OperatorAddress] = validator
	}

	for _, delegation := range oldStakingState.Delegations {
		// change the delegation from staking bonded(unbonded) pool to intermedary account
		DelAddr := sdk.AccAddress(delegation.DelegatorAddress)
		intermediaryAccount := types.IntermediaryAccount(DelAddr)

		validator, ok := validatorMap[delegation.ValidatorAddress]
		if !ok {
			return nil, fmt.Errorf("Error validator not found delegation %v", delegation.ValidatorAddress)
		}
		denom := validator.BondDenom

		if validator.Status == v0.Bonded {
			// if bonded => transfer from bondedPoolAddress to intermediary address
			amount := (delegation.Shares.MulInt(validator.Tokens)).Quo(validator.DelegatorShares)

			// Calculate new bonded token amount in bondedPoolAddress
		} else {
			// if unbonded or unbonding => transfer from unbondedPoolAddress to intermediary address
			amount := (delegation.Shares.MulInt(validator.Tokens)).Quo(validator.DelegatorShares)
			// Caculate new bonded token amount in unbondedPoolAddress
		}

	}

	for _, ubdDelegation := range oldStakingState.UnbondingDelegations {
		// change the delegation from staking unbonded pool to intermedary account
		DelAddr := sdk.AccAddress(ubdDelegation.DelegatorAddress)
		intermediaryAccount := types.IntermediaryAccount(DelAddr)

		validator, ok := validatorMap[ubdDelegation.ValidatorAddress]
		if !ok {
			return nil, fmt.Errorf("Error validator not found ubdDelegation %v", ubdDelegation.ValidatorAddress)
		}
		denom := validator.BondDenom

		for _, entry := range ubdDelegation.Entries {
			// move balance in entry from unbondedPoolAddress to intermediary address

			// Caculate new bonded token amount in unbondedPoolAddress
		}

	}

	return genesisState, nil
}

func migrateStaking(genesisState AppMap) (AppMap, error) {

	rawData := genesisState[stakingtypes.ModuleName]
	var oldState v0.GenesisState
	err := json.Unmarshal(rawData, &oldState)
	if err != nil {
		return nil, err
	}

	newState := types.GenesisState{}
	// Migrate state.StakingGenesisState
	stakingGenesisState := stakingtypes.GenesisState{}

	stakingGenesisState.Params = stakingtypes.Params(oldState.Params)
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

	newData, err := json.Marshal(&newState)
	if err != nil {
		return nil, err
	}

	genesisState[stakingtypes.ModuleName] = newData

	return genesisState, nil
}

func convertValidators(validators []v0.Validator) []stakingtypes.Validator {
	newValidators := make([]stakingtypes.Validator, 0)
	for _, val := range validators {
		newVal := stakingtypes.Validator{
			OperatorAddress: val.OperatorAddress,
			ConsensusPubkey: val.ConsensusPubkey,
			Jailed:          val.Jailed,
			Status:          stakingtypes.BondStatus(val.Status),
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