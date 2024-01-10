package keeper

import (
	"context"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

type queryServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &queryServer{
		Keeper: keeper,
	}
}

var _ types.QueryServer = queryServer{}

// BondWeight implements types.QueryServer.
func (k queryServer) BondWeight(c context.Context, req *types.QueryBondWeightRequest) (*types.QueryBondWeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	weight, found := k.Keeper.GetBondWeight(ctx, req.Denom)

	return &types.QueryBondWeightResponse{
		Weight: weight,
		Found:  found,
	}, nil
}

// MultiStakingLock implements types.QueryServer.
func (k queryServer) MultiStakingLock(c context.Context, req *types.QueryMultiStakingLockRequest) (*types.QueryMultiStakingLockResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	lockId := types.MultiStakingLockID(req.MultiStakerAddress, req.ValidatorAddress)
	lock, found := k.Keeper.GetMultiStakingLock(ctx, lockId)

	return &types.QueryMultiStakingLockResponse{
		Lock:  &lock,
		Found: found,
	}, nil
}

// MultiStakingLocks implements types.QueryServer.
func (k queryServer) MultiStakingLocks(c context.Context, req *types.QueryMultiStakingLocksRequest) (*types.QueryMultiStakingLocksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	var locks []*types.MultiStakingLock

	store := ctx.KVStore(k.storeKey)
	lockStore := prefix.NewStore(store, types.MultiStakingLockPrefix)

	pageRes, err := query.Paginate(lockStore, req.Pagination, func(key []byte, value []byte) error {
		var lock types.MultiStakingLock
		err := k.cdc.Unmarshal(value, &lock)
		if err != nil {
			return err
		}
		locks = append(locks, &lock)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMultiStakingLocksResponse{Locks: locks, Pagination: pageRes}, nil
}

// MultiStakingUnlock implements types.QueryServer.
func (k queryServer) MultiStakingUnlock(c context.Context, req *types.QueryMultiStakingUnlockRequest) (*types.QueryMultiStakingUnlockResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	unlockId := types.MultiStakingUnlockID(req.MultiStakerAddress, req.ValidatorAddress)
	unlock, found := k.Keeper.GetMultiStakingUnlock(ctx, unlockId)

	return &types.QueryMultiStakingUnlockResponse{
		Unlock: &unlock,
		Found:  found,
	}, nil
}

// MultiStakingUnlocks implements types.QueryServer.
func (k queryServer) MultiStakingUnlocks(c context.Context, req *types.QueryMultiStakingUnlocksRequest) (*types.QueryMultiStakingUnlocksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	var unlocks []*types.MultiStakingUnlock

	store := ctx.KVStore(k.storeKey)
	unlockStore := prefix.NewStore(store, types.MultiStakingUnlockPrefix)

	pageRes, err := query.Paginate(unlockStore, req.Pagination, func(key []byte, value []byte) error {
		var unlock types.MultiStakingUnlock
		err := k.cdc.Unmarshal(value, &unlock)
		if err != nil {
			return err
		}
		unlocks = append(unlocks, &unlock)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMultiStakingUnlocksResponse{Unlocks: unlocks, Pagination: pageRes}, nil
}

// ValidatorMultiStakingCoin implements types.QueryServer.
func (k queryServer) ValidatorMultiStakingCoin(c context.Context, req *types.QueryValidatorMultiStakingCoinRequest) (*types.QueryValidatorMultiStakingCoinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	denom := k.Keeper.GetValidatorMultiStakingCoin(ctx, sdk.ValAddress(req.ValidatorAddr))

	return &types.QueryValidatorMultiStakingCoinResponse{
		Denom: denom,
	}, nil
}
