package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

// SetMultiStakingUnlockEntry adds an entry to the unbonding delegation at
// the given addresses. It creates the unbonding delegation if it does not exist.
func (k Keeper) SetMultiStakingUnlockEntry(
	ctx sdk.Context, unlockID types.UnlockID,
	multistakingCoin types.MultiStakingCoin,
) types.MultiStakingUnlock {
	unlock, found := k.GetMultiStakingUnlock(ctx, unlockID)
	if found {
		unlock.AddEntry(ctx.BlockHeight(), multistakingCoin)
	} else {
		unlock = types.NewMultiStakingUnlock(ctx.BlockHeight(), multistakingCoin)
	}

	k.SetMultiStakingUnlock(ctx, unlock)
	return unlock
}

func (k Keeper) DeleteUnlockEntryAtCreationHeight(
	ctx sdk.Context, unlockID types.UnlockID,
	creationHeight int64,
) error {
	unlock, found := k.GetMultiStakingUnlock(ctx, unlockID)
	if found {
		unlock.RemoveEntryAtCreationHeight(creationHeight)
	} else {
		return fmt.Errorf("can't found unlock entry")
	}

	if len(unlock.Entries) == 0 {
		k.DeleteMultiStakingUnlock(ctx, unlockID)
		return nil
	}

	k.SetMultiStakingUnlock(ctx, unlock)
	return nil
}
