package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) IsAllowedCoin(ctx sdk.Context, valAcc sdk.ValAddress, lockedCoin sdk.Coin) bool {
	return lockedCoin.Denom == k.GetValidatorAllowedCoin(ctx, valAcc)
}
