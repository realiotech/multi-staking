package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {

}

func (k Keeper) EndBlocker(ctx sdk.Context, unbondedStakings types.UnbondedMultiStakings) {

}
