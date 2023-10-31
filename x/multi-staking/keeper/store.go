package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) GetBondTokenWeight(ctx sdk.Context, tokenDenom string) (math.LegacyDec, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBondTokenWeightKey(tokenDenom))
	if bz == nil {
		return math.LegacyZeroDec(), false
	}

	bondTokenWeight := &math.LegacyDec{}
	err := bondTokenWeight.Unmarshal(bz)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal bond token weight %v", err))

	}
	return *bondTokenWeight, true
}

func (k Keeper) SetBondTokenWeight(ctx sdk.Context, tokenDenom string, tokenWeight math.LegacyDec) {
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

func (k Keeper) GetIntermediaryAccountDelegator(ctx sdk.Context, intermediaryAccount sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetIntermediaryAccountDelegatorKey(intermediaryAccount))

	return bz
}

func (k Keeper) SetIntermediaryAccountDelegator(ctx sdk.Context, intermediaryAccount sdk.AccAddress, delegator sdk.AccAddress) {
	if k.GetIntermediaryAccountDelegator(ctx, intermediaryAccount) != nil {
		panic("intermediary account for delegator already set")
	}

	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetIntermediaryAccountDelegatorKey(intermediaryAccount), delegator)
}

func (k Keeper) GetDVPairSDKBondAmount(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) math.Int {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetDVPairSDKBondAmountKey(delAddr, valAddr))
	if bz == nil {
		return math.ZeroInt()
	}

	sdkBondAmount := &math.Int{}
	err := sdkBondAmount.Unmarshal(bz)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal sdk bond amount %v", err))
	}

	return *sdkBondAmount
}

func (k Keeper) SetDVPairSDKBondAmount(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, sdkBondAmount math.Int) {
	store := ctx.KVStore(k.storeKey)

	bz, err := sdkBondAmount.Marshal()

	if err != nil {
		panic(fmt.Errorf("unable to marshal sdk bond amount %v", err))
	}

	store.Set(types.GetDVPairSDKBondAmountKey(delAddr, valAddr), bz)
}

func (k Keeper) GetDVPairBondAmount(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) math.Int {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetDVPairBondAmountKey(delAddr, valAddr))
	if bz == nil {
		return sdk.ZeroInt()
	}

	bondAmount := &math.Int{}
	err := bondAmount.Unmarshal(bz)

	if err != nil {
		panic(fmt.Errorf("unable to unmarshal bond amount %v", err))
	}

	return *bondAmount
}

func (k Keeper) SetDVPairBondAmount(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, bondAmount math.Int) {
	store := ctx.KVStore(k.storeKey)

	bz, err := bondAmount.Marshal()

	if err != nil {
		panic(fmt.Errorf("unable to marshal bond amount %v", err))
	}

	store.Set(types.GetDVPairBondAmountKey(delAddr, valAddr), bz)
}
