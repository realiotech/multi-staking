package keeper

import (
	"context"
	"fmt"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetBondWeight(ctx sdk.Context, tokenDenom string) (math.LegacyDec, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBondWeightKey(tokenDenom))
	if bz == nil {
		return math.LegacyDec{}, false
	}

	bondCoinWeight := &math.LegacyDec{}
	err := bondCoinWeight.Unmarshal(bz)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal bond coin weight %v", err))
	}
	return *bondCoinWeight, true
}

func (k Keeper) SetBondWeight(ctx sdk.Context, tokenDenom string, tokenWeight math.LegacyDec) {
	store := ctx.KVStore(k.storeKey)
	bz, err := tokenWeight.Marshal()
	if err != nil {
		panic(fmt.Errorf("unable to marshal bond coin weight %v", err))
	}

	store.Set(types.GetBondWeightKey(tokenDenom), bz)
}

func (k Keeper) RemoveBondWeight(ctx sdk.Context, tokenDenom string) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetBondWeightKey(tokenDenom))
}

func (k Keeper) GetValidatorMultiStakingCoin(ctx sdk.Context, operatorAddr sdk.ValAddress) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorMultiStakingCoinKey(operatorAddr))

	return string(bz)
}

func (k Keeper) SetValidatorMultiStakingCoin(ctx sdk.Context, operatorAddr sdk.ValAddress, bondDenom string) {
	if k.GetValidatorMultiStakingCoin(ctx, operatorAddr) != "" {
		panic("validator multi staking coin already set")
	}

	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetValidatorMultiStakingCoinKey(operatorAddr), []byte(bondDenom))
}

func (k Keeper) ValidatorMultiStakingCoinIterator(ctx sdk.Context, cb func(valAddr string, denom string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.ValidatorMultiStakingCoinKey)
	iterator := storetypes.KVStorePrefixIterator(prefixStore, nil)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		valAddr := sdk.ValAddress(iterator.Key()).String()
		denom := string(iterator.Value())
		if cb(valAddr, denom) {
			break
		}
	}
}

func (k Keeper) GetMultiStakingLock(ctx context.Context, multiStakingLockID types.LockID) (types.MultiStakingLock, bool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)

	bz := store.Get(multiStakingLockID.ToBytes())
	if bz == nil {
		return types.MultiStakingLock{}, false
	}

	multiStakingLock := types.MultiStakingLock{}
	k.cdc.MustUnmarshal(bz, &multiStakingLock)
	return multiStakingLock, true
}

func (k Keeper) SetMultiStakingLock(ctx context.Context, multiStakingLock types.MultiStakingLock) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if multiStakingLock.IsEmpty() {
		k.RemoveMultiStakingLock(sdkCtx, multiStakingLock.LockID)
		return
	}

	store := sdkCtx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&multiStakingLock)

	store.Set(multiStakingLock.LockID.ToBytes(), bz)
}

func (k Keeper) RemoveMultiStakingLock(ctx sdk.Context, multiStakingLockID types.LockID) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(multiStakingLockID.ToBytes())
}

func (k Keeper) MultiStakingLockIterator(ctx sdk.Context, cb func(stakingLock types.MultiStakingLock) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.MultiStakingLockPrefix)
	iterator := storetypes.KVStorePrefixIterator(prefixStore, nil)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var multiStakingLock types.MultiStakingLock
		k.cdc.MustUnmarshal(iterator.Value(), &multiStakingLock)
		if cb(multiStakingLock) {
			break
		}
	}
}

func (k Keeper) MultiStakingUnlockIterator(ctx sdk.Context, cb func(multiStakingUnlock types.MultiStakingUnlock) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.MultiStakingUnlockPrefix)
	iterator := storetypes.KVStorePrefixIterator(prefixStore, nil)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var multiStakingUnlock types.MultiStakingUnlock
		k.cdc.MustUnmarshal(iterator.Value(), &multiStakingUnlock)
		if cb(multiStakingUnlock) {
			break
		}
	}
}

func (k Keeper) BondWeightIterator(ctx sdk.Context, cb func(denom string, bondWeight math.LegacyDec) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.BondWeightKey)
	iterator := storetypes.KVStorePrefixIterator(prefixStore, nil)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		denom := string(iterator.Key())
		bondWeight := &math.LegacyDec{}
		err := bondWeight.Unmarshal(iterator.Value())
		if err != nil {
			panic(fmt.Errorf("unable to unmarshal bond coin weight %v", err))
		}
		if cb(denom, *bondWeight) {
			break
		}
	}
}

func (k Keeper) GetMultiStakingUnlock(ctx sdk.Context, multiStakingUnlockID types.UnlockID) (unlock types.MultiStakingUnlock, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(multiStakingUnlockID.ToBytes())

	if value == nil {
		return unlock, false
	}

	unlock = types.MultiStakingUnlock{}
	k.cdc.MustUnmarshal(value, &unlock)

	return unlock, true
}

// SetMultiStakingUnlock sets the unbonding delegation and associated index.
func (k Keeper) SetMultiStakingUnlock(ctx sdk.Context, unlock types.MultiStakingUnlock) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&unlock)

	store.Set(unlock.UnlockID.ToBytes(), bz)
}

func (k Keeper) DeleteMultiStakingUnlock(ctx sdk.Context, unlockID types.UnlockID) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(unlockID.ToBytes())
}
