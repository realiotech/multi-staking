package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) GetBondWeight(ctx sdk.Context, tokenDenom string) (sdk.Dec, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBondWeightKey(tokenDenom))
	if bz == nil {
		return sdk.Dec{}, false
	}

	bondCoinWeight := &sdk.Dec{}
	err := bondCoinWeight.Unmarshal(bz)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal bond coin weight %v", err))

	}
	return *bondCoinWeight, true
}

func (k Keeper) SetBondWeight(ctx sdk.Context, tokenDenom string, tokenWeight sdk.Dec) {
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
	bz := store.Get(types.GetValidatorMultiStakingCoinKey(operatorAddr.String()))

	return string(bz)
}

func (k Keeper) SetValidatorMultiStakingCoin(ctx sdk.Context, operatorAddr sdk.ValAddress, bondDenom string) {
	if k.GetValidatorMultiStakingCoin(ctx, operatorAddr) != "" {
		panic("validator denom already set")
	}

	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetValidatorMultiStakingCoinKey(operatorAddr.String()), []byte(bondDenom))
}

func (k Keeper) ValidatorMultiStakingCoinIterator(ctx sdk.Context, cb func(valAddr string, denom string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.ValidatorMultiStakingCoinKey)
	iterator := sdk.KVStorePrefixIterator(prefixStore, nil)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		valAddr := string(iterator.Key())
		denom := string(iterator.Value())
		if cb(valAddr, denom) {
			break
		}
	}
}

func (k Keeper) GetIntermediaryDelegatorKey(ctx sdk.Context, multiStakerAddr sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetIntermediaryDelegatorKey(multiStakerAddr))

	return bz
}

func (k Keeper) IsIntermediaryDelegator(ctx sdk.Context, intermediaryDelegator sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetIntermediaryDelegatorKey(intermediaryDelegator))

	return bz != nil
}

func (k Keeper) SetIntermediaryDelegator(ctx sdk.Context, intermediaryDelegator sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetIntermediaryDelegatorKey(intermediaryDelegator), []byte{0x1})
}

func (k Keeper) GetMultiStakingLock(ctx sdk.Context, multiStakingLockID types.LockID) (types.MultiStakingLock, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(multiStakingLockID.ToByte())
	if bz == nil {
		return types.MultiStakingLock{}, false
	}

	multiStakingLock := types.MultiStakingLock{}
	k.cdc.MustUnmarshal(bz, &multiStakingLock)
	return multiStakingLock, true
}

func (k Keeper) SetMultiStakingLock(ctx sdk.Context, multiStakingLock types.MultiStakingLock) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&multiStakingLock)

	store.Set(multiStakingLock.LockID.ToByte(), bz)
}

func (k Keeper) RemoveMultiStakingLock(ctx sdk.Context, multiStakingLockID types.LockID) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(multiStakingLockID.ToByte())
}

func (k Keeper) MultiStakingLockIterator(ctx sdk.Context, cb func(stakingLock types.MultiStakingLock) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.MultiStakingLockPrefix)
	iterator := sdk.KVStorePrefixIterator(prefixStore, nil)

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
	iterator := sdk.KVStorePrefixIterator(prefixStore, nil)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var multiStakingUnlock types.MultiStakingUnlock
		k.cdc.MustUnmarshal(iterator.Value(), &multiStakingUnlock)
		if cb(multiStakingUnlock) {
			break
		}
	}
}

func (k Keeper) BondWeightIterator(ctx sdk.Context, cb func(denom string, bondWeight sdk.Dec) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.BondWeightKey)
	iterator := sdk.KVStorePrefixIterator(prefixStore, nil)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		denom := string(iterator.Key())
		bondWeight := &sdk.Dec{}
		err := bondWeight.Unmarshal(iterator.Value())
		if err != nil {
			panic(fmt.Errorf("unable to unmarshal bond coin weight %v", err))

		}
		if cb(denom, *bondWeight) {
			break
		}
	}
}

func (k Keeper) IntermediaryDelegatorIterator(ctx sdk.Context, cb func(intermediaryDelegator sdk.AccAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.IntermediaryDelegatorKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		intermediaryDelegator := sdk.AccAddress(iterator.Key())

		if cb(intermediaryDelegator) {
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
