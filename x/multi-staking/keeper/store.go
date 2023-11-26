package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) GetBondTokenWeight(ctx sdk.Context, tokenDenom string) (sdk.Dec, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBondTokenWeightKey(tokenDenom))
	if bz == nil {
		return sdk.Dec{}, false
	}

	bondTokenWeight := &sdk.Dec{}
	err := bondTokenWeight.Unmarshal(bz)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal bond token weight %v", err))

	}
	return *bondTokenWeight, true
}

func (k Keeper) SetBondTokenWeight(ctx sdk.Context, tokenDenom string, tokenWeight sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	bz, err := tokenWeight.Marshal()

	if err != nil {
		panic(fmt.Errorf("unable to marshal bond token weight %v", err))
	}

	store.Set(types.GetBondTokenWeightKey(tokenDenom), bz)
}

func (k Keeper) GetValidatorBondDenom(ctx sdk.Context, operatorAddr sdk.ValAddress) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorBondDenomKey(operatorAddr))

	return string(bz)
}

func (k Keeper) SetValidatorBondDenom(ctx sdk.Context, operatorAddr sdk.ValAddress, bondDenom string) {
	if k.GetValidatorBondDenom(ctx, operatorAddr) != "" {
		panic("validator denom already set")
	}

	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetValidatorBondDenomKey(operatorAddr), []byte(bondDenom))
}

func (k Keeper) GetIntermediaryAccountDelegator(ctx sdk.Context, delAcc sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetIntermediaryAccountDelegatorKey(delAcc))

	return bz
}

func (k Keeper) SetIntermediaryAccountDelegator(ctx sdk.Context, intermediaryAccount sdk.AccAddress, delegator sdk.AccAddress) {
	if k.GetIntermediaryAccountDelegator(ctx, intermediaryAccount) != nil {
		panic("intermediary account for delegator already set")
	}

	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetIntermediaryAccountDelegatorKey(intermediaryAccount), delegator)
}

func (k Keeper) GetMultiStakingLock(ctx sdk.Context, multiStakingLockID []byte) (types.MultiStakingLock, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(multiStakingLockID)
	if bz == nil {
		return types.MultiStakingLock{}, false
	}

	multiStakingLock := &types.MultiStakingLock{}
	k.cdc.MustUnmarshal(bz, multiStakingLock)
	return *multiStakingLock, true
}

func (k Keeper) SetMultiStakingLock(ctx sdk.Context, multiStakingLockID []byte, multiStakingLock types.MultiStakingLock) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&multiStakingLock)

	store.Set(multiStakingLockID, bz)
}

func (k Keeper) RemoveMultiStakingLock(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.MultiStakingLockID(delAddr, valAddr))
}
