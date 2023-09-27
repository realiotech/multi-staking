package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) GetBondTokenWeight(ctx sdk.Context, tokenDenom string) sdk.Dec {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBondTokenWeightKey(tokenDenom))

	bondTokenWeight := &sdk.Dec{}
	err := bondTokenWeight.Unmarshal(bz)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal bond token weight %v", err))

	}

	return *bondTokenWeight
}
