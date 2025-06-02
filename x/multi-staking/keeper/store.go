package keeper

import (
	"context"
	"fmt"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetBondWeight(ctx context.Context, tokenDenom string) (math.LegacyDec, bool) {
	store := k.storeService.OpenKVStore(ctx)
	bz, _ := store.Get(types.GetBondWeightKey(tokenDenom))
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

func (k Keeper) SetBondWeight(ctx context.Context, tokenDenom string, tokenWeight math.LegacyDec) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := tokenWeight.Marshal()
	if err != nil {
		panic(fmt.Errorf("unable to marshal bond coin weight %v", err))
	}

	err = store.Set(types.GetBondWeightKey(tokenDenom), bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) RemoveBondWeight(ctx context.Context, tokenDenom string) {
	store := k.storeService.OpenKVStore(ctx)

	err := store.Delete(types.GetBondWeightKey(tokenDenom))
	if err != nil {
		panic(err)
	}
}

func (k Keeper) GetValidatorMultiStakingCoin(ctx context.Context, operatorAddr sdk.ValAddress) string {
	store := k.storeService.OpenKVStore(ctx)
	bz, _ := store.Get(types.GetValidatorMultiStakingCoinKey(operatorAddr))

	return string(bz)
}

func (k Keeper) SetValidatorMultiStakingCoin(ctx context.Context, operatorAddr sdk.ValAddress, bondDenom string) {
	if k.GetValidatorMultiStakingCoin(ctx, operatorAddr) != "" {
		panic("validator multi staking coin already set")
	}

	store := k.storeService.OpenKVStore(ctx)

	err := store.Set(types.GetValidatorMultiStakingCoinKey(operatorAddr), []byte(bondDenom))
	if err != nil {
		panic(err)
	}
}

func (k Keeper) ValidatorMultiStakingCoinIterator(ctx context.Context, cb func(valAddr string, denom string) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
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
	store := k.storeService.OpenKVStore(ctx)

	bz, _ := store.Get(multiStakingLockID.ToBytes())

	if bz == nil {
		return types.MultiStakingLock{}, false
	}

	multiStakingLock := types.MultiStakingLock{}
	k.cdc.MustUnmarshal(bz, &multiStakingLock)
	return multiStakingLock, true
}

func (k Keeper) SetMultiStakingLock(ctx context.Context, multiStakingLock types.MultiStakingLock) {
	if multiStakingLock.IsEmpty() {
		k.RemoveMultiStakingLock(ctx, multiStakingLock.LockID)
		return
	}

	store := k.storeService.OpenKVStore(ctx)

	bz := k.cdc.MustMarshal(&multiStakingLock)

	err := store.Set(multiStakingLock.LockID.ToBytes(), bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) RemoveMultiStakingLock(ctx context.Context, multiStakingLockID types.LockID) {
	store := k.storeService.OpenKVStore(ctx)

	err := store.Delete(multiStakingLockID.ToBytes())
	if err != nil {
		panic(err)
	}
}

func (k Keeper) MultiStakingLockIterator(ctx context.Context, cb func(stakingLock types.MultiStakingLock) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
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

func (k Keeper) MultiStakingUnlockIterator(ctx context.Context, cb func(multiStakingUnlock types.MultiStakingUnlock) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
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

func (k Keeper) BondWeightIterator(ctx context.Context, cb func(denom string, bondWeight math.LegacyDec) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
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

func (k Keeper) GetMultiStakingUnlock(ctx context.Context, multiStakingUnlockID types.UnlockID) (unlock types.MultiStakingUnlock, found bool) {
	store := k.storeService.OpenKVStore(ctx)
	value, _ := store.Get(multiStakingUnlockID.ToBytes())

	if value == nil {
		return unlock, false
	}

	unlock = types.MultiStakingUnlock{}
	k.cdc.MustUnmarshal(value, &unlock)

	return unlock, true
}

// SetMultiStakingUnlock sets the unbonding delegation and associated index.
func (k Keeper) SetMultiStakingUnlock(ctx context.Context, unlock types.MultiStakingUnlock) {
	store := k.storeService.OpenKVStore(ctx)

	bz := k.cdc.MustMarshal(&unlock)

	err := store.Set(unlock.UnlockID.ToBytes(), bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) DeleteMultiStakingUnlock(ctx context.Context, unlockID types.UnlockID) {
	store := k.storeService.OpenKVStore(ctx)

	err := store.Delete(unlockID.ToBytes())
	if err != nil {
		panic(err)
	}
}
