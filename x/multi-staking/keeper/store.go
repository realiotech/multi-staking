package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
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

func (k Keeper) GetValidatorAllowedToken(ctx sdk.Context, operatorAddr sdk.ValAddress) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorAllowedTokenKey(operatorAddr.String()))

	return string(bz)
}

func (k Keeper) SetValidatorAllowedToken(ctx sdk.Context, operatorAddr sdk.ValAddress, bondDenom string) {
	if k.GetValidatorAllowedToken(ctx, operatorAddr) != "" {
		panic("validator denom already set")
	}

	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetValidatorAllowedTokenKey(operatorAddr.String()), []byte(bondDenom))
}

func (k Keeper) ValidatorAllowedTokenIterator(ctx sdk.Context, cb func(valAddr string, denom string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.ValidatorAllowedTokenKey)
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
	prefixStore := prefix.NewStore(store, types.MultiStakingLockPrefix)

	bz := prefixStore.Get(multiStakingLockID)
	if bz == nil {
		return types.MultiStakingLock{}, false
	}

	multiStakingLock := types.MultiStakingLock{}
	k.cdc.MustUnmarshal(bz, &multiStakingLock)
	return multiStakingLock, true
}

func (k Keeper) SetMultiStakingLock(ctx sdk.Context, multiStakingLockID []byte, multiStakingLock types.MultiStakingLock) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.MultiStakingLockPrefix)

	bz := k.cdc.MustMarshal(&multiStakingLock)

	prefixStore.Set(multiStakingLockID, bz)
}

func (k Keeper) RemoveMultiStakingLock(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.MultiStakingLockPrefix)

	prefixStore.Delete(types.MultiStakingLockID(delAddr, valAddr))
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

func (k Keeper) GetUnbondedMultiStaking(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (ubd types.UnbondedMultiStaking, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDKey(delAddr, valAddr)
	value := store.Get(key)

	if value == nil {
		return ubd, false
	}

	ubd = types.MustUnmarshalUBD(k.cdc, value)

	return ubd, true
}

// SetUnbondedMultiStaking sets the unbonding delegation and associated index.
func (k Keeper) SetUnbondedMultiStaking(ctx sdk.Context, ubd types.UnbondedMultiStaking) {
	delAddr := sdk.MustAccAddressFromBech32(ubd.DelegatorAddress)

	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalUBD(k.cdc, ubd)
	valAddr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	key := types.GetUBDKey(delAddr, valAddr)
	store.Set(key, bz)
}

// RemoveUnbondedMultiStaking removes the unbonding delegation object and associated index.
func (k Keeper) RemoveUnbondedMultiStaking(ctx sdk.Context, ubd types.UnbondedMultiStaking) {
	delegatorAddress := sdk.MustAccAddressFromBech32(ubd.DelegatorAddress)

	store := ctx.KVStore(k.storeKey)
	addr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	key := types.GetUBDKey(delegatorAddress, addr)
	store.Delete(key)
}

// SetUnbondedMultiStakingEntry adds an entry to the unbonding delegation at
// the given addresses. It creates the unbonding delegation if it does not exist.
func (k Keeper) SetUnbondedMultiStakingEntry(
	ctx sdk.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
	creationHeight int64, rate sdk.Dec, minTime time.Time, balance math.Int,
) types.UnbondedMultiStaking {
	ubd, found := k.GetUnbondedMultiStaking(ctx, delegatorAddr, validatorAddr)
	if found {
		ubd.AddEntry(creationHeight, minTime,rate, balance)
	} else {
		ubd = types.NewUnbondedMultiStaking(delegatorAddr, validatorAddr, creationHeight, rate, minTime, balance)
	}

	k.SetUnbondedMultiStaking(ctx, ubd)
	return ubd
}
