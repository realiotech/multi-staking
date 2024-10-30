package keeper

import (
	"context"
	"fmt"
	"time"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Keeper struct {
	storeKey      storetypes.StoreKey
	cdc           codec.BinaryCodec
	accountKeeper types.AccountKeeper
	stakingKeeper *stakingkeeper.Keeper
	bankKeeper    types.BankKeeper
	authority     string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	accountKeeper types.AccountKeeper,
	stakingKeeper *stakingkeeper.Keeper,
	bankKeeper types.BankKeeper,
	key storetypes.StoreKey,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      key,
		accountKeeper: accountKeeper,
		stakingKeeper: stakingKeeper,
		bankKeeper:    bankKeeper,
		authority:     authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) DequeueAllMatureUBDQueue(ctx context.Context, currTime time.Time) (matureUnbonds []stakingtypes.DVPair, err error) {
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	unbondingTimesliceIterator, err := k.stakingKeeper.UBDQueueIterator(ctx, currTime)
	if err != nil {
		return nil, err
	}
	defer unbondingTimesliceIterator.Close()

	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := stakingtypes.DVPairs{}
		value := unbondingTimesliceIterator.Value()
		k.cdc.MustUnmarshal(value, &timeslice)

		matureUnbonds = append(matureUnbonds, timeslice.Pairs...)
	}

	return matureUnbonds, nil
}

func (k Keeper) GetMatureUnbondingDelegations(ctx context.Context) ([]stakingtypes.UnbondingDelegation, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var matureUnbondingDelegations []stakingtypes.UnbondingDelegation
	matureUnbonds, err := k.DequeueAllMatureUBDQueue(ctx, sdkCtx.BlockHeader().Time)
	if err != nil {
		return nil, err
	}
	for _, dvPair := range matureUnbonds {
		delAddr, valAddr, err := types.AccAddrAndValAddrFromStrings(dvPair.DelegatorAddress, dvPair.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		unbondingDelegation, err := k.stakingKeeper.GetUnbondingDelegation(ctx, delAddr, valAddr)
		if err != nil {
			if sdkerrors.IsOf(err, stakingtypes.ErrNoUnbondingDelegation) {
				continue
			}

			return nil, err
		}

		matureUnbondingDelegations = append(matureUnbondingDelegations, unbondingDelegation)
	}
	return matureUnbondingDelegations, nil
}

func (k Keeper) GetUnbondingEntryAtCreationHeight(ctx sdk.Context, delAcc sdk.AccAddress, valAcc sdk.ValAddress, creationHeight int64) (stakingtypes.UnbondingDelegationEntry, bool) {
	ubd, err := k.stakingKeeper.GetUnbondingDelegation(ctx, delAcc, valAcc)
	if err != nil {
		return stakingtypes.UnbondingDelegationEntry{}, false
	}

	var unbondingEntryAtHeight stakingtypes.UnbondingDelegationEntry
	found := false
	for _, entry := range ubd.Entries {
		if entry.CreationHeight == creationHeight {
			if !found {
				found = true
				unbondingEntryAtHeight = entry
			} else {
				unbondingEntryAtHeight.Balance = unbondingEntryAtHeight.Balance.Add(entry.Balance)
			}
		}
	}

	return unbondingEntryAtHeight, found
}

func (k Keeper) BurnCoin(ctx context.Context, accAddr sdk.AccAddress, coin sdk.Coin) error {
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, accAddr, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) isValMultiStakingCoin(ctx context.Context, valAcc sdk.ValAddress, lockedCoin sdk.Coin) bool {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return lockedCoin.Denom == k.GetValidatorMultiStakingCoin(sdkCtx, valAcc)
}

func (k Keeper) AdjustUnbondAmount(ctx context.Context, delAcc sdk.AccAddress, valAcc sdk.ValAddress, amount math.Int) (adjustedAmount math.Int, err error) {
	delegation, err := k.stakingKeeper.GetDelegation(ctx, delAcc, valAcc)
	if err != nil {
		return math.Int{}, fmt.Errorf("failed to get delegation: %s", err.Error())
	}
	validator, err := k.stakingKeeper.GetValidator(ctx, valAcc)
	if err != nil {
		return math.Int{}, fmt.Errorf("failed to get validator: %s", err.Error())
	}

	shares, err := validator.SharesFromTokens(amount)
	if err != nil {
		return math.Int{}, err
	}

	delShares := delegation.GetShares()
	// Cap the shares at the delegation's shares. Shares being greater could occur
	// due to rounding, however we don't want to truncate the shares or take the
	// minimum because we want to allow for the full withdraw of shares from a
	// delegation.
	if shares.GT(delShares) {
		shares = delShares
	}

	return validator.TokensFromShares(shares).TruncateInt(), nil
}

func (k Keeper) AdjustCancelUnbondingAmount(ctx context.Context, delAcc sdk.AccAddress, valAcc sdk.ValAddress, creationHeight int64, amount math.Int) (adjustedAmount math.Int, err error) {
	undelegation, err := k.stakingKeeper.GetUnbondingDelegation(ctx, delAcc, valAcc)
	if err != nil {
		return math.Int{}, fmt.Errorf("failed to get undelegation: %s", err.Error())
	}

	totalUnbondingAmount := math.ZeroInt()
	for _, entry := range undelegation.Entries {
		if entry.CreationHeight == creationHeight {
			totalUnbondingAmount = totalUnbondingAmount.Add(entry.Balance)
		}
	}

	return math.MinInt(totalUnbondingAmount, amount), nil
}

func (k Keeper) BondDenom(ctx context.Context) string {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	bondDenom := k.GetParams(sdkCtx).MainBondDenom
	return bondDenom
}

func (k Keeper) IterateDelegations(ctx context.Context, delegator sdk.AccAddress, fn func(index int64, delegation stakingtypes.DelegationI) (stop bool)) error {
	return k.stakingKeeper.IterateDelegations(ctx, delegator, fn)
}
