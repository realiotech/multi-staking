package keeper

import (
	"fmt"

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

func (k Keeper) RemoveBondTokenWeight(ctx sdk.Context, tokenDenom string) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.GetBondTokenWeightKey(tokenDenom))
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

	bz := store.Get(multiStakingLockID)
	if bz == nil {
		return types.MultiStakingLock{}, false
	}

	multiStakingLock := types.MultiStakingLock{}
	k.cdc.MustUnmarshal(bz, &multiStakingLock)
	return multiStakingLock, true
}

func (k Keeper) SetMultiStakingLock(ctx sdk.Context, multiStakingLock types.MultiStakingLock) error {
	store := ctx.KVStore(k.storeKey)

	delAcc, valAcc, err := types.DelAccAndValAccFromStrings(multiStakingLock.DelAddr, multiStakingLock.ValAddr)
	if err != nil {
		return err
	}

	lockID := types.MultiStakingLockID(delAcc, valAcc)
	bz := k.cdc.MustMarshal(&multiStakingLock)

	store.Set(lockID, bz)

	return err
}

func (k Keeper) RemoveMultiStakingLock(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)

	store.Delete(types.MultiStakingLockID(delAddr, valAddr))
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
func (k Keeper) SetMultiStakingUnlock(ctx sdk.Context, unlock types.MultiStakingUnlock) {
	delAddr := sdk.MustAccAddressFromBech32(unlock.DelegatorAddress)
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&unlock)
	valAddr, err := sdk.ValAddressFromBech32(unlock.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	key := types.MultiStakingUnlockID(delAddr, valAddr)
	store.Set(key, bz)
}

// RemoveMultiStakingUnlock removes the unbonding delegation object and associated index.
func (k Keeper) RemoveMultiStakingUnlock(ctx sdk.Context, unlock types.MultiStakingUnlock) {
	delegatorAddress := sdk.MustAccAddressFromBech32(unlock.DelegatorAddress)

	store := ctx.KVStore(k.storeKey)
	addr, err := sdk.ValAddressFromBech32(unlock.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	key := types.MultiStakingUnlockID(delegatorAddress, addr)
	store.Delete(key)
}

// SetMultiStakingUnlockEntry adds an entry to the unbonding delegation at
// the given addresses. It creates the unbonding delegation if it does not exist.
func (k Keeper) SetMultiStakingUnlockEntry(
	ctx sdk.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
	creationHeight int64, rate sdk.Dec, balance math.Int,
) types.MultiStakingUnlock {
	unlock, found := k.GetMultiStakingUnlock(ctx, types.MultiStakingUnlockID(delegatorAddr, validatorAddr))
	if found {
		unlock.AddEntry(creationHeight, rate, balance)
	} else {
		unlock = types.NewMultiStakingUnlock(delegatorAddr, validatorAddr, creationHeight, rate, balance)
	}

	k.SetMultiStakingUnlock(ctx, unlock)
	return unlock
}
