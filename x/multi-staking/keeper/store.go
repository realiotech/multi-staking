package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) GetBondCoinWeight(ctx sdk.Context, tokenDenom string) (sdk.Dec, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBondCoinWeightKey(tokenDenom))
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

func (k Keeper) SetBondCoinWeight(ctx sdk.Context, tokenDenom string, tokenWeight sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz, err := tokenWeight.Marshal()

	if err != nil {
		panic(fmt.Errorf("unable to marshal bond coin weight %v", err))
	}

	store.Set(types.GetBondCoinWeightKey(tokenDenom), bz)
}

func (k Keeper) RemoveBondCoinWeight(ctx sdk.Context, tokenDenom string) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetBondCoinWeightKey(tokenDenom))
}

func (k Keeper) GetValidatorAllowedCoin(ctx sdk.Context, operatorAddr sdk.ValAddress) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorAllowedCoinKey(operatorAddr.String()))

	return string(bz)
}

func (k Keeper) SetValidatorAllowedCoin(ctx sdk.Context, operatorAddr sdk.ValAddress, bondDenom string) {
	if k.GetValidatorAllowedCoin(ctx, operatorAddr) != "" {
		panic("validator denom already set")
	}

	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetValidatorAllowedCoinKey(operatorAddr.String()), []byte(bondDenom))
}

func (k Keeper) ValidatorAllowedCoinIterator(ctx sdk.Context, cb func(valAddr string, denom string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.ValidatorAllowedCoinKey)
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

func (k Keeper) GetIntermediaryAccountKey(ctx sdk.Context, delAcc sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetIntermediaryAccountKey(delAcc))

	return bz
}

func (k Keeper) IsIntermediaryAccount(ctx sdk.Context, intermediaryAccount sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetIntermediaryAccountKey(intermediaryAccount))

	return bz != nil
}

func (k Keeper) SetIntermediaryAccount(ctx sdk.Context, intermediaryAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetIntermediaryAccountKey(intermediaryAccount), []byte{0x1})
}

func (k Keeper) GetMultiStakingLock(ctx sdk.Context, multiStakingLockID []byte) (types.MultiStakingLock, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(multiStakingLockID)
	if bz == nil {
		return types.MultiStakingLock{}, false
	}

	multiStakingLock := types.MultiStakingLock{}
	k.cdc.MustUnmarshal(bz, &multiStakingLock)
	return multiStakingLock, true
}

func (k Keeper) SetMultiStakingLock(ctx sdk.Context, multiStakingLockID []byte, multiStakingLock types.MultiStakingLock) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&multiStakingLock)

	store.Set(multiStakingLockID, bz)
}

func (k Keeper) RemoveMultiStakingLock(ctx sdk.Context, multiStakingLockID []byte) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(multiStakingLockID)
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

func (k Keeper) GetMultiStakingUnlock(ctx sdk.Context, multiStakingUnlockID []byte) (unlock types.MultiStakingUnlock, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(multiStakingUnlockID)

	if value == nil {
		return unlock, false
	}

	unlock = types.MultiStakingUnlock{}
	k.cdc.MustUnmarshal(value, &unlock)

	return unlock, true
}

// SetMultiStakingUnlock sets the unbonding delegation and associated index.
func (k Keeper) SetMultiStakingUnlock(ctx sdk.Context, unlockID []byte, unlock types.MultiStakingUnlock) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&unlock)

	store.Set(unlockID, bz)
}

// RemoveMultiStakingUnlock removes the unbonding delegation object and associated index.
func (k Keeper) RemoveMultiStakingUnlock(ctx sdk.Context, unlockID []byte) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(unlockID)
}

// SetMultiStakingUnlockEntry adds an entry to the unbonding delegation at
// the given addresses. It creates the unbonding delegation if it does not exist.
func (k Keeper) SetMultiStakingUnlockEntry(
	ctx sdk.Context, unlockID []byte,
	rate sdk.Dec, balance math.Int,
) types.MultiStakingUnlock {
	unlock, found := k.GetMultiStakingUnlock(ctx, unlockID)
	if found {
		unlock.AddEntry(ctx.BlockHeight(), rate, balance)
	} else {
		unlock = types.NewMultiStakingUnlock(ctx.BlockHeight(), rate, balance)
	}

	k.SetMultiStakingUnlock(ctx, unlockID, unlock)
	return unlock
}
