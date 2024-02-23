package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

// handleAddMultiStakingCoinProposal handles the proposals to add a new bond token
func (k Keeper) AddMultiStakingCoinProposal(
	ctx sdk.Context,
	p *types.AddMultiStakingCoinProposal,
) error {

	_, found := k.GetBondWeight(ctx, p.Denom)
	if found {
		return fmt.Errorf("Error MultiStakingCoin %s already exist", p.Denom)
	}

	k.SetBondWeight(ctx, p.Denom, *p.BondWeight)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddMultiStakingCoin,
			sdk.NewAttribute(types.AttributeKeyDenom, p.Denom),
			sdk.NewAttribute(types.AttributeKeyBondWeight, p.BondWeight.String()),
		),
	)
	return nil
}
