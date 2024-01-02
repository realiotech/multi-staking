package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) GetBondWeight(ctx sdk.Context, tokenDenom string) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBondWeightKey(tokenDenom))

	bondTokenWeight := &sdk.Dec{}
	err := bondTokenWeight.Unmarshal(bz)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal bond weight %v", err))

	}

	return *bondTokenWeight
}

func (k Keeper) SetBondWeight(ctx sdk.Context, tokenDenom string, bondWeight sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz, err := bondWeight.Marshal()

	if err != nil {
		panic(fmt.Errorf("unable to marshal bond weight %v", err))
	}

	store.Set(types.GetBondWeightKey(tokenDenom), bz)
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

func (k Keeper) GetIntermediaryDelegator(ctx sdk.Context, intermediaryAccount sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetIntermediaryDelegatorKey(intermediaryAccount))

	return bz
}

func (k Keeper) SetIntermediaryDelegator(ctx sdk.Context, intermediaryAccount sdk.AccAddress, delegator sdk.AccAddress) {
	if k.GetIntermediaryDelegator(ctx, intermediaryAccount) != nil {
		panic("intermediary delegator already set")
	}

	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetIntermediaryDelegatorKey(intermediaryAccount), delegator)
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
	if multiStakingLock.IsEmpty() {
		k.RemoveMultiStakingLock(ctx, *multiStakingLock.LockID)
		return
	}

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

func (k Keeper) RemoveMultiStakingUnlock(ctx sdk.Context, unlockID types.UnlockID) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(unlockID.ToBytes())
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

func (k Keeper) GetDVPairSDKBondTokens(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) sdk.Coin {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetDVPairSDKBondTokensKey(delAddr, valAddr))
	var sdkBondTokens sdk.Coin
	k.cdc.MustUnmarshal(bz, &sdkBondTokens)

	return sdkBondTokens
}

func (k Keeper) SetDVPairSDKBondTokens(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, sdkBondTokens sdk.Coin) {
	if sdkBondTokens.Denom != k.stakingKeeper.BondDenom(ctx) {
		panic("input token is not sdk bond token")
	}
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&sdkBondTokens)
	store.Set(types.GetDVPairSDKBondTokensKey(delAddr, valAddr), bz)
}

func (k Keeper) GetDVPairBondTokens(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) sdk.Coin {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetDVPairBondTokensKey(delAddr, valAddr))
	var bondTokens sdk.Coin
	k.cdc.MustUnmarshal(bz, &bondTokens)

	return bondTokens
}

func (k Keeper) SetDVPairBondTokens(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, bondTokens sdk.Coin) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&bondTokens)
	store.Set(types.GetDVPairBondTokensKey(delAddr, valAddr), bz)
}
