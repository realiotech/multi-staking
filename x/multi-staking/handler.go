package multistaking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

// NewBondDenomProposalHandler returns bond denom module's proposals
func NewBondDenomProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.AddBondDenomProposal:
			return keeper.HandlerAddBondDenomProposal(ctx, &k, c)
		case *types.UpdateBondTokenWeightProposals:
			return keeper.HandlerUpdateBondTokenWeightProposals(ctx, &k, c)
		case *types.RemoveBondTokenProposal:
			keeper.HandlerRemoveBondTokenProposal(ctx, &k, c)
			return nil
		default:
			return types.ErrUnrecognized
		}
	}
}
