package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func HandlerAddBondDenomProposal(ctx sdk.Context, k *Keeper, p *types.AddBondDenomProposal) error {
	_, found := k.GetBondCoinWeight(ctx, p.BondCoinAdd)
	if found {
		return fmt.Errorf("denom %s already exists", p.BondCoinAdd)
	}

	k.SetBondCoinWeight(ctx, p.BondCoinAdd, *p.BondCoinWeightAdd)
	return nil
}

func HandlerUpdateBondCoinWeightProposals(ctx sdk.Context, k *Keeper, p *types.UpdateBondCoinWeightProposals) error {
	_, found := k.GetBondCoinWeight(ctx, p.BondDenomChange)
	if !found {
		return fmt.Errorf("denom %s does not exist", p.BondDenomChange)
	}
	k.RemoveBondCoinWeight(ctx, p.BondDenomChange)

	k.SetBondCoinWeight(ctx, p.BondDenomChange, *p.BondCoinWeightChange)
	return nil
}

func HandlerRemoveBondCoinProposal(ctx sdk.Context, k *Keeper, p *types.RemoveBondCoinProposal) {
	_, found := k.GetBondCoinWeight(ctx, p.BondCoinRemove)
	if !found {
		return
	}

	k.RemoveBondCoinWeight(ctx, p.BondCoinRemove)
}
