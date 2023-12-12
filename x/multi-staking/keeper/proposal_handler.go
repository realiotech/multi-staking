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

	if p.BondTokenWeightAdd.LTE(sdk.NewDec(0)) {
		return fmt.Errorf("bondTokenWeight must greater than zero")
	}

	k.SetBondTokenWeight(ctx, p.BondTokenAdd, *p.BondTokenWeightAdd)
	return nil
}

func HandlerUpdateBondTokenWeightProposals(ctx sdk.Context, k *Keeper, p *types.UpdateBondTokenWeightProposals) error {
	_, found := k.GetBondTokenWeight(ctx, p.BondDenomChange)
	if !found {
		return fmt.Errorf("denom %s does not exist", p.BondDenomChange)
	}
	if p.BondTokenWeightChange.LTE(sdk.NewDec(0)) {
		return fmt.Errorf("bondTokenWeight must greater than zero")
	}
	k.RemoveBondTokenWeight(ctx, p.BondDenomChange)

	k.SetBondTokenWeight(ctx, p.BondDenomChange, *p.BondTokenWeightChange)
	return nil
}

func HandlerRemoveBondTokenProposal(ctx sdk.Context, k *Keeper, p *types.RemoveBondTokenProposal) {
	_, found := k.GetBondTokenWeight(ctx, p.BondTokenRemove)
	if !found {
		return
	}

	k.RemoveBondTokenWeight(ctx, p.BondTokenRemove)
}
