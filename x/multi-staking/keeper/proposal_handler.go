package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func HandlerAddMultiStakingCoinProposal(ctx sdk.Context, k *Keeper, p *types.AddMultiStakingCoinProposal) error {
	_, found := k.GetBondCoinWeight(ctx, p.Denom)
	if found {
		return fmt.Errorf("denom %s already exists", p.Denom)
	}

	k.SetBondCoinWeight(ctx, p.Denom, *p.BondWeight)
	return nil
}

func HandlerUpdateBondWeightProposals(ctx sdk.Context, k *Keeper, p *types.UpdateBondWeightProposal) error {
	_, found := k.GetBondCoinWeight(ctx, p.Denom)
	if !found {
		return fmt.Errorf("denom %s does not exist", p.Denom)
	}

	k.SetBondCoinWeight(ctx, p.Denom, *p.UpdatedBondWeight)
	return nil
}

func HandlerRemoveMultiStakingCoinProposal(ctx sdk.Context, k *Keeper, p *types.RemoveMultiStakingCoinProposal) {
	_, found := k.GetBondCoinWeight(ctx, p.Denom)
	if !found {
		return
	}

	k.RemoveBondCoinWeight(ctx, p.Denom)
}
