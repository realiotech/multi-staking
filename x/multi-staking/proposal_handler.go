package multistaking

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

// NewMultiStakingProposalHandler creates a governance handler to manage Mult-Staking proposals.
func NewMultiStakingProposalHandler(k *keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {

		switch c := content.(type) {
		case *types.AddBondTokenProposal:
			return handleAddBondTokenProposal(ctx, k, c)
		case *types.ChangeBondTokenWeightProposal:
			return handleChangeTokenWeightProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(errortypes.ErrUnknownRequest, "unrecognized %s proposal content type: %T", types.ModuleName, c)
		}
	}
}

// handleAddBondTokenProposal handles the proposals to add a new bond token 
func handleAddBondTokenProposal(
	ctx sdk.Context,
	k *keeper.Keeper,
	p *types.AddBondTokenProposal,
) error {
	k.SetBondTokenWeight(ctx, p.BondToken, p.TokenWeight)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddBondToken,
			sdk.NewAttribute(types.AttributeKeyBondToken, p.BondToken),
			sdk.NewAttribute(types.AttributeKeyBondTokenWeight, p.TokenWeight.String()),
		),
	)
	return nil
}

// handleChangeTokenWeightProposal handles the proposals to change a bond tokens weight
func handleChangeTokenWeightProposal(
	ctx sdk.Context,
	k *keeper.Keeper,
	p *types.ChangeBondTokenWeightProposal,
) error {
	k.SetBondTokenWeight(ctx, p.BondToken, p.TokenWeight)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeBondTokenWeight,
			sdk.NewAttribute(types.AttributeKeyBondToken, p.BondToken),
			sdk.NewAttribute(types.AttributeKeyBondTokenWeight, p.TokenWeight.String()),
		),
	)
	return nil
}