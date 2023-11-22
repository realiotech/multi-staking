package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

)

func (k Keeper) GetBondTokenWeight(ctx sdk.Context, tokenDenom string) math.LegacyDec {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBondTokenWeightKey(tokenDenom))

	bondTokenWeight := &math.LegacyDec{}
	err := bondTokenWeight.Unmarshal(bz)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal bond token weight %v", err))

	}

	return *bondTokenWeight
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
