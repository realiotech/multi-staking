package multistaking

import (
	"github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	sdkerrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// NewMultiStakingProposalHandler creates a governance handler to manage Mult-Staking proposals.
func NewMultiStakingProposalHandler(k *keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.AddMultiStakingCoinProposal:
			return handleAddMultiStakingCoinProposal(ctx, k, c)
		case *types.UpdateBondWeightProposal:
			return handleChangeTokenWeightProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(errortypes.ErrUnknownRequest, "unrecognized %s proposal content type: %T", types.ModuleName, c)
		}
	}
}

// handleAddMultiStakingCoinProposal handles the proposals to add a new bond token
func handleAddMultiStakingCoinProposal(
	ctx sdk.Context,
	k *keeper.Keeper,
	p *types.AddMultiStakingCoinProposal,
) error {
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

// handleChangeTokenWeightProposal handles the proposals to change a bond tokens weight
func handleChangeTokenWeightProposal(
	ctx sdk.Context,
	k *keeper.Keeper,
	p *types.UpdateBondWeightProposal,
) error {
	k.SetBondWeight(ctx, p.Denom, *p.UpdatedBondWeight)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateBondWeight,
			sdk.NewAttribute(types.AttributeKeyDenom, p.Denom),
			sdk.NewAttribute(types.AttributeKeyBondWeight, p.UpdatedBondWeight.String()),
		),
	)
	return nil
}
