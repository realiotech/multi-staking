package keeper

import (
	"context"
	"strings"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

// BondWeights implements types.QueryServer.
func (k queryServer) MultiStakingCoinInfos(c context.Context, req *types.QueryMultiStakingCoinInfosRequest) (*types.QueryMultiStakingCoinInfosResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	var infos []*types.MultiStakingCoinInfo

	store := ctx.KVStore(k.storeKey)
	coinInfoStore := prefix.NewStore(store, types.BondWeightKey)

	pageRes, err := query.Paginate(coinInfoStore, req.Pagination, func(key []byte, value []byte) error {
		bondCoinWeight := &sdk.Dec{}
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
	valAcc, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid validator address")
	}

	denom := k.Keeper.GetValidatorMultiStakingCoin(ctx, valAcc)

	return &types.QueryValidatorMultiStakingCoinResponse{
		Denom: denom,
	}, nil
}

func (k queryServer) Validators(c context.Context, req *types.QueryValidatorsRequest) (*types.QueryValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// validate the provided status, return all the validators if the status is empty
	if req.Status != "" && !(req.Status == stakingtypes.Bonded.String() || req.Status == stakingtypes.Unbonded.String() || req.Status == stakingtypes.Unbonding.String()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator status %s", req.Status)
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	valStore := prefix.NewStore(store, stakingtypes.ValidatorsKey)

	validators, pageRes, err := query.GenericFilteredPaginate(k.cdc, valStore, req.Pagination, func(key []byte, val *stakingtypes.Validator) (*stakingtypes.Validator, error) {
		if req.Status != "" && !strings.EqualFold(val.GetStatus().String(), req.Status) {
			return nil, nil
		}

		return val, nil
	}, func() *stakingtypes.Validator {
		return &stakingtypes.Validator{}
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	vals := []types.ValidatorInfo{}
	for _, val := range validators {
		valAcc, err := sdk.ValAddressFromBech32(val.OperatorAddress)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid validator address")
		}

		denom := k.Keeper.GetValidatorMultiStakingCoin(ctx, valAcc)
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

	return &types.QueryValidatorsResponse{Validators: vals, Pagination: pageRes}, nil
}

func (k queryServer) Validator(c context.Context, req *types.QueryValidatorRequest) (*types.QueryValidatorResponse, error) {
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

	ctx := sdk.UnwrapSDKContext(c)
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, status.Errorf(codes.NotFound, "validator %s not found", req.ValidatorAddr)
	}

	denom := k.Keeper.GetValidatorMultiStakingCoin(ctx, valAddr)
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
