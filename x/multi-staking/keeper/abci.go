package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {

}

func (k Keeper) EndBlocker(ctx sdk.Context, unbondedStakings types.UnbondedMultiStaking) {
	// for _, unbondedStaking := range unbondedStakings.Entries {
	// 	// k.BurnBondTokenAndUnlockMultiStakingToken(ctx, sdk.AccAddress(unbondedStaking.), sdk.ValAddress(unbondedStaking.ValAddr), unbondedStaking.Amount[0])
	// }
}
