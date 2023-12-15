package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

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

func (q queryServer) ValidatorAllowedToken(c context.Context, params *types.QueryValidatorAllowedTokenRequest) (*types.QueryValidatorAllowedTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	operatorAddress, err := sdk.ValAddressFromBech32(params.OperatorAddress)
	if err != nil {
		return &types.QueryValidatorAllowedTokenResponse{}, nil
	}

	allowedToken := q.Keeper.GetValidatorAllowedToken(ctx, operatorAddress)

	return &types.QueryValidatorAllowedTokenResponse{
		Denom: allowedToken,
	}, nil
}

func (q queryServer) MultiStakingLock(c context.Context, params *types.QueryMultiStakingLockRequest) (*types.QueryMultiStakingLockResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	multiStakingLock, found := q.Keeper.GetMultiStakingLock(ctx, params.MultiStakingLockId)

	return &types.QueryMultiStakingLockResponse{
		MultiStakingLock: &multiStakingLock,
		Found:            found,
	}, nil
}

func (q queryServer) MultiStakingUnlock(c context.Context, params *types.QueryMultiStakingUnlockRequest) (*types.QueryMultiStakingUnlockReponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	delAddr, valAddr, err := types.DelAccAndValAccFromStrings(params.DelegatorAddress, params.ValidatorAddress)
	if err != nil {
		return &types.QueryMultiStakingUnlockReponse{}, nil
	}

	multiStakingUnlock, found := q.Keeper.GetMultiStakingUnlock(ctx, delAddr, valAddr)

	return &types.QueryMultiStakingUnlockReponse{
		MultiStakingUnlock: &multiStakingUnlock,
		Found:              found,
	}, nil
}
