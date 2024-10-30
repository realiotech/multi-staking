package keeper

import (
	"context"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type queryServer struct {
	Keeper
	stakingQuerier stakingkeeper.Querier
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &queryServer{
		Keeper: keeper,
		stakingQuerier: stakingkeeper.Querier{
			Keeper: keeper.stakingKeeper,
		},
	}
}

var _ types.QueryServer = queryServer{}

// BondWeights implements types.QueryServer.
func (k queryServer) MultiStakingCoinInfos(ctx context.Context, req *types.QueryMultiStakingCoinInfosRequest) (*types.QueryMultiStakingCoinInfosResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var infos []*types.MultiStakingCoinInfo

	store := sdkCtx.KVStore(k.storeKey)
	coinInfoStore := prefix.NewStore(store, types.BondWeightKey)

	pageRes, err := query.Paginate(coinInfoStore, req.Pagination, func(key []byte, value []byte) error {
		bondCoinWeight := &math.LegacyDec{}
		err := bondCoinWeight.Unmarshal(value)
		if err != nil {
			return err
		}
		coinInfo := types.MultiStakingCoinInfo{
			Denom:      string(key),
			BondWeight: *bondCoinWeight,
		}

		infos = append(infos, &coinInfo)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMultiStakingCoinInfosResponse{Infos: infos, Pagination: pageRes}, nil
}

// BondWeight implements types.QueryServer.
func (k queryServer) BondWeight(ctx context.Context, req *types.QueryBondWeightRequest) (*types.QueryBondWeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	weight, found := k.Keeper.GetBondWeight(sdkCtx, req.Denom)

	return &types.QueryBondWeightResponse{
		Weight: weight,
		Found:  found,
	}, nil
}

// MultiStakingLock implements types.QueryServer.
func (k queryServer) MultiStakingLock(ctx context.Context, req *types.QueryMultiStakingLockRequest) (*types.QueryMultiStakingLockResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	lockId := types.MultiStakingLockID(req.MultiStakerAddress, req.ValidatorAddress)
	lock, found := k.Keeper.GetMultiStakingLock(sdkCtx, lockId)

	return &types.QueryMultiStakingLockResponse{
		Lock:  &lock,
		Found: found,
	}, nil
}

// MultiStakingLocks implements types.QueryServer.
func (k queryServer) MultiStakingLocks(ctx context.Context, req *types.QueryMultiStakingLocksRequest) (*types.QueryMultiStakingLocksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var locks []*types.MultiStakingLock

	store := sdkCtx.KVStore(k.storeKey)
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
func (k queryServer) MultiStakingUnlock(ctx context.Context, req *types.QueryMultiStakingUnlockRequest) (*types.QueryMultiStakingUnlockResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	unlockId := types.MultiStakingUnlockID(req.MultiStakerAddress, req.ValidatorAddress)
	unlock, found := k.Keeper.GetMultiStakingUnlock(sdkCtx, unlockId)

	return &types.QueryMultiStakingUnlockResponse{
		Unlock: &unlock,
		Found:  found,
	}, nil
}

// MultiStakingUnlocks implements types.QueryServer.
func (k queryServer) MultiStakingUnlocks(ctx context.Context, req *types.QueryMultiStakingUnlocksRequest) (*types.QueryMultiStakingUnlocksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var unlocks []*types.MultiStakingUnlock

	store := sdkCtx.KVStore(k.storeKey)
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
func (k queryServer) ValidatorMultiStakingCoin(ctx context.Context, req *types.QueryValidatorMultiStakingCoinRequest) (*types.QueryValidatorMultiStakingCoinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	valAcc, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid validator address")
	}

	denom := k.Keeper.GetValidatorMultiStakingCoin(sdkCtx, valAcc)

	return &types.QueryValidatorMultiStakingCoinResponse{
		Denom: denom,
	}, nil
}

func (k queryServer) Validators(ctx context.Context, req *types.QueryValidatorsRequest) (*types.QueryValidatorsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkReq := stakingtypes.QueryValidatorsRequest{
		Status:     req.Status,
		Pagination: req.Pagination,
	}

	resp, err := k.stakingQuerier.Validators(ctx, &sdkReq)
	if err != nil {
		return nil, err
	}

	vals := []types.ValidatorInfo{}
	for _, val := range resp.Validators {
		valAcc, err := sdk.ValAddressFromBech32(val.OperatorAddress)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid validator address")
		}

		denom := k.Keeper.GetValidatorMultiStakingCoin(sdkCtx, valAcc)
		valInfo := types.ValidatorInfo{
			OperatorAddress:   val.OperatorAddress,
			ConsensusPubkey:   val.ConsensusPubkey,
			Jailed:            val.Jailed,
			Status:            val.Status,
			Tokens:            val.Tokens,
			DelegatorShares:   val.DelegatorShares,
			Description:       val.Description,
			UnbondingHeight:   val.UnbondingHeight,
			UnbondingTime:     val.UnbondingTime,
			Commission:        val.Commission,
			MinSelfDelegation: val.MinSelfDelegation,
			BondDenom:         denom,
		}
		vals = append(vals, valInfo)
	}

	return &types.QueryValidatorsResponse{Validators: vals, Pagination: resp.Pagination}, nil
}

func (k queryServer) Validator(ctx context.Context, req *types.QueryValidatorRequest) (*types.QueryValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ValidatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	validator, err := k.stakingKeeper.GetValidator(sdkCtx, valAddr)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get validator with address %s: %s", req.ValidatorAddr, err.Error())
	}

	denom := k.Keeper.GetValidatorMultiStakingCoin(sdkCtx, valAddr)
	valInfo := types.ValidatorInfo{
		OperatorAddress:   validator.OperatorAddress,
		ConsensusPubkey:   validator.ConsensusPubkey,
		Jailed:            validator.Jailed,
		Status:            validator.Status,
		Tokens:            validator.Tokens,
		DelegatorShares:   validator.DelegatorShares,
		Description:       validator.Description,
		UnbondingHeight:   validator.UnbondingHeight,
		UnbondingTime:     validator.UnbondingTime,
		Commission:        validator.Commission,
		MinSelfDelegation: validator.MinSelfDelegation,
		BondDenom:         denom,
	}
	return &types.QueryValidatorResponse{Validator: valInfo}, nil
}
