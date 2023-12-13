package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func HandlerAddBondDenomProposal(ctx sdk.Context, k *Keeper, p *types.AddBondDenomProposal) error {
	_, found := k.GetBondTokenWeight(ctx, p.BondTokenAdd)
	if found {
		return fmt.Errorf("denom %s already exists", p.BondTokenAdd)
	}

	k.SetBondTokenWeight(ctx, p.BondTokenAdd, *p.BondTokenWeightAdd)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddBondToken,
			sdk.NewAttribute(types.AttributeKeyBondToken, p.BondTokenAdd),
			sdk.NewAttribute(types.AttributeKeyBondTokenWeight, p.BondTokenWeightAdd.String()),
		),
	)
	return nil
}

func HandlerUpdateBondTokenWeightProposals(ctx sdk.Context, k *Keeper, p *types.UpdateBondTokenWeightProposals) error {
	_, found := k.GetBondTokenWeight(ctx, p.BondDenomChange)
	if !found {
		return fmt.Errorf("denom %s does not exist", p.BondDenomChange)
	}
	k.RemoveBondTokenWeight(ctx, p.BondDenomChange)

	k.SetBondTokenWeight(ctx, p.BondDenomChange, *p.BondTokenWeightChange)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateBondTokenWeight,
			sdk.NewAttribute(types.AttributeKeyBondToken, p.BondDenomChange),
			sdk.NewAttribute(types.AttributeKeyBondTokenWeight, p.BondTokenWeightChange.String()),
		),
	)
	return nil
}

func HandlerRemoveBondTokenProposal(ctx sdk.Context, k *Keeper, p *types.RemoveBondTokenProposal) {
	_, found := k.GetBondTokenWeight(ctx, p.BondTokenRemove)
	if !found {
		return
	}

	k.RemoveBondTokenWeight(ctx, p.BondTokenRemove)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveBondTokenWeight,
			sdk.NewAttribute(types.AttributeKeyBondToken, p.BondTokenRemove),
		),
	)
}
