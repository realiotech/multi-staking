package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) GetSDKDelegation(ctx sdk.Context, delAcc sdk.AccAddress, valAcc sdk.ValAddress) (stakingtypes.Delegation, bool) {
	intermediaryAccount := types.GetIntermediaryAccount(delAcc.String(), valAcc.String())

	return k.stakingKeeper.GetDelegation(ctx, intermediaryAccount, valAcc)
}

func (k Keeper) IsAllowedToken(ctx sdk.Context, valAcc sdk.ValAddress, lockedToken sdk.Coin) bool {
	return lockedToken.Denom == k.GetValidatorBondDenom(ctx, valAcc)
}

func (k Keeper) MoveLockedTokenAndRedelegate(ctx sdk.Context, delAcc sdk.AccAddress, srcValAcc sdk.ValAddress, dstValAcc sdk.ValAddress) {
	intermediaryAccount := k.GetIntermediaryAccountDelegator(ctx, delAcc)

}
