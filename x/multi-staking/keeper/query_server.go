package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

// NewQueryServerImpl returns an implementation of the module QueryServer.

var _ types.QueryServer = queryServer{}

func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{
		Keeper: k,
	}
}

type queryServer struct {
	Keeper
}

func (q queryServer) BondTokenWeight(c context.Context, params *types.QueryBondTokenWeightRequest) (*types.QueryBondTokenWeightResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	weight, isSet := q.Keeper.GetBondTokenWeight(ctx, params.TokenDenom)

	return &types.QueryBondTokenWeightResponse{
		Weight: weight,
		IsSet:  isSet,
	}, nil
}
