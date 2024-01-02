package keeper

import (
	"fmt"

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
