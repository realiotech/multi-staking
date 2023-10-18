package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) Delegate(
	ctx sdk.Context, delAddress string, valAddress string,
	bondAmt sdk.Coin, sdkBondAmt sdk.Coin,
) {
	delAcc := sdk.AccAddress(delAddress)
	valAcc := sdk.ValAddress(valAddress)
	intermediaryAcc := types.GetIntermediaryAccount(delAddress)

	if k.GetIntermediaryAccountDelegator(ctx, intermediaryAcc) == nil {
		k.SetIntermediaryAccountDelegator(ctx, intermediaryAcc, delAcc)
	}

	oldBondTokens := k.GetDVPairBondTokens(ctx, delAcc, valAcc)
	if oldBondTokens.Denom == "" {
		k.SetDVPairBondTokens(ctx, delAcc, valAcc, bondAmt)

	} else {
		k.SetDVPairBondTokens(ctx, delAcc, valAcc, bondAmt.Add(oldBondTokens))
	}

	oldSDKBondTokens := k.GetDVPairSDKBondTokens(ctx, delAcc, valAcc)
	if oldSDKBondTokens.Denom == "" {
		k.SetDVPairSDKBondTokens(ctx, delAcc, valAcc, sdkBondAmt)

	} else {
		k.SetDVPairSDKBondTokens(ctx, delAcc, valAcc, sdkBondAmt.Add(oldSDKBondTokens))
	}
}
