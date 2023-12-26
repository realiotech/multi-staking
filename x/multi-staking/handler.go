package multistaking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

// NewBondDenomProposalHandler returns bond denom module's proposals
func NewBondDenomProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.AddMultiStakingCoinProposal:
			return keeper.HandlerAddMultiStakingCoinProposal(ctx, &k, c)
		case *types.UpdateBondWeightProposal:
			return keeper.HandlerUpdateBondWeightProposals(ctx, &k, c)
		case *types.RemoveMultiStakingCoinProposal:
			keeper.HandlerRemoveMultiStakingCoinProposal(ctx, &k, c)
			return nil
		default:
			return fmt.Errorf("unrecognized brond denom proposal content type")
		}
	}
}
