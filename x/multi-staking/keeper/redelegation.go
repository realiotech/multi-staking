package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) 

func (k Keeper) PreRedelegate(
	ctx sdk.Context, delAcc sdk.AccAddress, srcValAcc sdk.ValAddress, dstValAcc sdk.ValAddress,
	bondToken sdk.Coin,
) (sdk.Coin, error) {
	// check if bond denom match src val's bond denom
	srcValBondDenom := k.GetValidatorBondDenom(ctx, srcValAcc)
	if bondToken.Denom != srcValBondDenom {
		return sdk.Coin{}, fmt.Errorf("mismatch bond token; expect %s got %s", srcValBondDenom, bondToken.Denom)
	}

	dstValBondDenom := k.GetValidatorBondDenom(ctx, dstValAcc)
	if bondToken.Denom != dstValBondDenom {
		return sdk.Coin{}, fmt.Errorf("mismatch bond token; expect %s got %s", dstValBondDenom, bondToken.Denom)
	}





	// calculate converted sdk bond token
	sdkBondAmount, err := k.CalSDKBondAmount(ctx, bondToken)
	if err != nil {
		return sdk.Coin{}, err
	}
	sdkBondToken := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), sdkBondAmount)

	intermediaryAcc := types.GetIntermediaryAccount(delAcc.String(), valAcc.String())

	k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdkBondToken))

	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, intermediaryAcc, sdk.NewCoins(sdkBondToken))

	if k.GetIntermediaryAccountDelegator(ctx, intermediaryAcc) == nil {
		k.SetIntermediaryAccountDelegator(ctx, intermediaryAcc, delAcc)
	}

	k.UpdateDVPairBondAmount(ctx, delAcc, valAcc, bondToken.Amount)

	k.UpdateDVPairSDKBondAmount(ctx, delAcc, valAcc, sdkBondAmount)
	return sdkBondToken, nil
}
