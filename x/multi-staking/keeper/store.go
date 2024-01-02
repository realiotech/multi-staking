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

func (k Keeper) GetIntermediaryDelegator(ctx sdk.Context, intermediaryDelegator sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetIntermediaryDelegatorKey(intermediaryDelegator))

	return bz
}

func (k Keeper) SetIntermediaryDelegator(ctx sdk.Context, intermediaryDelegator sdk.AccAddress, delegator sdk.AccAddress) {
	if k.GetIntermediaryDelegator(ctx, intermediaryDelegator) != nil {
		panic("intermediary delegator already set")
	}

	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetIntermediaryDelegatorKey(intermediaryDelegator), []byte{1})
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
