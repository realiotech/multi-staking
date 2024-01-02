package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
	}
}

var _ types.MsgServer = msgServer{}

// CreateValidator defines a method for creating a new validator
func (k msgServer) CreateValidator(goCtx context.Context, msg *types.MsgCreateValidator) (*types.MsgCreateValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	multiStakerAddr, valAcc, err := types.AccAddrAndValAddrFromStrings(msg.MultiStakerAddress, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	intermediaryDelegator := types.IntermediaryDelegator(multiStakerAddr)
	if !k.IsIntermediaryDelegator(ctx, intermediaryDelegator) {
		k.SetIntermediaryDelegator(ctx, intermediaryDelegator)
	}

	lockID := types.MultiStakingLockID(msg.MultiStakerAddress, msg.ValidatorAddress)

	mintedBondCoin, err := k.Keeper.LockCoinAndMintBondCoin(ctx, lockID, multiStakerAddr, intermediaryDelegator, msg.Value)
	if err != nil {
		return nil, err
	}

	sdkMsg := stakingtypes.MsgCreateValidator{
		Description:       msg.Description,
		Commission:        msg.Commission,
		MinSelfDelegation: msg.MinSelfDelegation,
		DelegatorAddress:  intermediaryDelegator.String(),
		ValidatorAddress:  msg.ValidatorAddress,
		Pubkey:            msg.Pubkey,
		Value:             mintedBondCoin,
	}

	k.SetValidatorMultiStakingCoin(ctx, valAcc, msg.Value.Denom)

	_, err = k.stakingMsgServer.CreateValidator(ctx, &sdkMsg)

	if err != nil {
		return nil, err
	}

	return &types.MsgCreateValidatorResponse{}, nil
}

// EditValidator defines a method for editing an existing validator
func (k msgServer) EditValidator(goCtx context.Context, msg *types.MsgEditValidator) (*types.MsgEditValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sdkMsg := stakingtypes.MsgEditValidator{
		Description:       msg.Description,
		CommissionRate:    msg.CommissionRate,
		MinSelfDelegation: msg.MinSelfDelegation,
		ValidatorAddress:  msg.ValidatorAddress,
	}

	_, err := k.stakingMsgServer.EditValidator(ctx, &sdkMsg)
	if err != nil {
		return nil, err
	}
	return &types.MsgEditValidatorResponse{}, nil
}

// Delegate defines a method for performing a delegation of coins from a delegator to a validator
func (k msgServer) Delegate(goCtx context.Context, msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	multiStakerAddr, valAcc, err := types.AccAddrAndValAddrFromStrings(msg.MultiStakerAddress, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	if !k.isValMultiStakingCoin(ctx, valAcc, msg.Amount) {
		return nil, fmt.Errorf("not allowed coin")
	}

	intermediaryDelegator := types.IntermediaryDelegator(multiStakerAddr)
	if !k.IsIntermediaryDelegator(ctx, intermediaryDelegator) {
		k.SetIntermediaryDelegator(ctx, intermediaryDelegator)
	}

	lockID := types.MultiStakingLockID(msg.MultiStakerAddress, msg.ValidatorAddress)

	mintedBondCoin, err := k.Keeper.LockCoinAndMintBondCoin(ctx, lockID, multiStakerAddr, intermediaryDelegator, msg.Amount)
	if err != nil {
		return nil, err
	}

	sdkMsg := stakingtypes.MsgDelegate{
		DelegatorAddress: intermediaryDelegator.String(),
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           mintedBondCoin,
	}

	_, err = k.stakingMsgServer.Delegate(ctx, &sdkMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgDelegateResponse{}, nil
}

// BeginRedelegate defines a method for performing a redelegation of coins from a delegator and source validator to a destination validator
func (k msgServer) BeginRedelegate(goCtx context.Context, msg *types.MsgBeginRedelegate) (*types.MsgBeginRedelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	multiStakerAddr := sdk.MustAccAddressFromBech32(msg.MultiStakerAddress)
	intermediaryDelegator := types.IntermediaryDelegator(multiStakerAddr)

	srcValAcc, err := sdk.ValAddressFromBech32(msg.ValidatorSrcAddress)
	if err != nil {
		return nil, err
	}
	dstValAcc, err := sdk.ValAddressFromBech32(msg.ValidatorDstAddress)
	if err != nil {
		return nil, err
	}

	if !k.isValMultiStakingCoin(ctx, srcValAcc, msg.Amount) || !k.isValMultiStakingCoin(ctx, dstValAcc, msg.Amount) {
		return nil, fmt.Errorf("not allowed Coin")
	}

	fromLockID := types.MultiStakingLockID(msg.MultiStakerAddress, msg.ValidatorSrcAddress)
	fromLock, found := k.GetMultiStakingLock(ctx, fromLockID)
	if !found {
		return nil, fmt.Errorf("lock not found")
	}

	toLockID := types.MultiStakingLockID(msg.MultiStakerAddress, msg.ValidatorDstAddress)
	toLock := k.GetOrCreateMultiStakingLock(ctx, toLockID)

	multiStakingCoin := fromLock.MultiStakingCoin(msg.Amount.Amount)

	err = fromLock.MoveCoinToLock(&toLock, multiStakingCoin)
	if err != nil {
		return nil, err
	}
	k.SetMultiStakingLock(ctx, fromLock)
	k.SetMultiStakingLock(ctx, toLock)

	bondAmount := multiStakingCoin.BondAmount()
	bondAmount, err = k.AdjustUnbondAmount(ctx, multiStakerAddr, srcValAcc, bondAmount)
	if err != nil {
		return nil, err
	}

	bondCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), bondAmount)

	sdkMsg := &stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    intermediaryDelegator.String(),
		ValidatorSrcAddress: msg.ValidatorSrcAddress,
		ValidatorDstAddress: msg.ValidatorDstAddress,
		Amount:              bondCoin,
	}
	_, err = k.stakingMsgServer.BeginRedelegate(goCtx, sdkMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgBeginRedelegateResponse{}, err
}

// Undelegate defines a method for performing an undelegation from a delegate and a validator
func (k msgServer) Undelegate(goCtx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	multiStakerAddr, valAcc, err := types.AccAddrAndValAddrFromStrings(msg.MultiStakerAddress, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	if !k.isValMultiStakingCoin(ctx, valAcc, msg.Amount) {
		return nil, fmt.Errorf("not allowed coin")
	}

	lockID := types.MultiStakingLockID(msg.MultiStakerAddress, msg.ValidatorAddress)
	lock, found := k.GetMultiStakingLock(ctx, lockID)
	if !found {
		return nil, fmt.Errorf("can't find multi staking lock")
	}

	multiStakingCoin := lock.MultiStakingCoin(msg.Amount.Amount)
	err = lock.RemoveCoinFromMultiStakingLock(multiStakingCoin)
	if err != nil {
		return nil, err
	}
	k.SetMultiStakingLock(ctx, lock)

	unbondAmount := multiStakingCoin.BondAmount()

	unbondCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), unbondAmount)
	intermediaryDelegator := types.IntermediaryDelegator(multiStakerAddr)

	sdkMsg := &stakingtypes.MsgUndelegate{
		DelegatorAddress: intermediaryDelegator.String(),
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           unbondCoin,
	}

	_, err = k.stakingMsgServer.Undelegate(goCtx, sdkMsg)
	if err != nil {
		return nil, err
	}

	k.SetMultiStakingUnlockEntry(ctx, types.MultiStakingUnlockID(msg.MultiStakerAddress, msg.ValidatorAddress), multiStakingCoin)

	return &types.MsgUndelegateResponse{}, err
}

// // CancelUnbondingDelegation defines a method for canceling the unbonding delegation
// // and delegate back to the validator.
func (k msgServer) CancelUnbondingDelegation(goCtx context.Context, msg *types.MsgCancelUnbondingDelegation) (*types.MsgCancelUnbondingDelegationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	multiStakerAddr, valAcc, err := types.AccAddrAndValAddrFromStrings(msg.MultiStakerAddress, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	intermediaryDelegator := types.IntermediaryDelegator(multiStakerAddr)
	if !k.isValMultiStakingCoin(ctx, valAcc, msg.Amount) {
		return nil, fmt.Errorf("not allow coin")
	}

	unbondEntry, found := k.GetUnbondingEntryAtCreationHeight(ctx, intermediaryDelegator, valAcc, msg.CreationHeight)
	if !found {
		return nil, fmt.Errorf("unbondEntry not found")
	}
	cancelUnbondingCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), unbondEntry.Balance)

	sdkMsg := &stakingtypes.MsgCancelUnbondingDelegation{
		DelegatorAddress: intermediaryDelegator.String(),
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           cancelUnbondingCoin,
	}

	_, err = k.stakingMsgServer.CancelUnbondingDelegation(goCtx, sdkMsg)
	if err != nil {
		return nil, err
	}

	unlockID := types.MultiStakingUnlockID(msg.MultiStakerAddress, msg.ValidatorAddress)
	err = k.DeleteUnlockEntryAtCreationHeight(ctx, unlockID, msg.CreationHeight)
	if err != nil {
		return nil, err
	}

	return &types.MsgCancelUnbondingDelegationResponse{}, nil
}

// SetWithdrawAddress defines a method for performing an undelegation from a delegate and a validator
func (k msgServer) SetWithdrawAddress(goCtx context.Context, msg *types.MsgSetWithdrawAddress) (*types.MsgSetWithdrawAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	multiStakerAddr, err := sdk.AccAddressFromBech32(msg.MultiStakerAddress)
	if err != nil {
		return nil, err
	}
	intermediaryDelegator := types.IntermediaryDelegator(multiStakerAddr)

	sdkMsg := distrtypes.MsgSetWithdrawAddress{
		DelegatorAddress: intermediaryDelegator.String(),
		WithdrawAddress:  msg.WithdrawAddress,
	}

	_, err = k.distrMsgServer.SetWithdrawAddress(ctx, &sdkMsg)
	if err != nil {
		return nil, err
	}
	return &types.MsgSetWithdrawAddressResponse{}, nil
}

func (k msgServer) WithdrawDelegatorReward(goCtx context.Context, msg *types.MsgWithdrawDelegatorReward) (*types.MsgWithdrawDelegatorRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	multiStakerAddr, err := sdk.AccAddressFromBech32(msg.MultiStakerAddress)
	if err != nil {
		return nil, err
	}

	intermediaryDelegator := types.IntermediaryDelegator(multiStakerAddr)

	if !k.IsIntermediaryDelegator(ctx, intermediaryDelegator) {
		k.SetIntermediaryDelegator(ctx, intermediaryDelegator)
	}

	sdkMsg := distrtypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: intermediaryDelegator.String(),
		ValidatorAddress: msg.ValidatorAddress,
	}

	resp, err := k.distrMsgServer.WithdrawDelegatorReward(ctx, &sdkMsg)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.SendCoins(ctx, intermediaryDelegator, multiStakerAddr, resp.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgWithdrawDelegatorRewardResponse{Amount: resp.Amount}, nil
}

func (k msgServer) Vote(goCtx context.Context, msg *types.MsgVote) (*types.MsgVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	multiStakerAddr, err := sdk.AccAddressFromBech32(msg.Voter)
	if err != nil {
		return nil, err
	}
	intermediaryDelegator := types.IntermediaryDelegator(multiStakerAddr)
	if !k.IsIntermediaryDelegator(ctx, intermediaryDelegator) {
		k.SetIntermediaryDelegator(ctx, intermediaryDelegator)
	}

	sdkMsg := govv1.MsgVote{
		ProposalId: msg.ProposalId,
		Voter:      intermediaryDelegator.String(),
		Option:     msg.Option,
		Metadata:   msg.Metadata,
	}

	_, err = k.govMsgServer.Vote(ctx, &sdkMsg)
	if err != nil {
		return nil, err
	}
	return &types.MsgVoteResponse{}, nil
}

func (k msgServer) VoteWeighted(goCtx context.Context, msg *types.MsgVoteWeighted) (*types.MsgVoteWeightedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	multiStakerAddr, err := sdk.AccAddressFromBech32(msg.Voter)
	if err != nil {
		return nil, err
	}

	intermediaryDelegator := types.IntermediaryDelegator(multiStakerAddr)
	if !k.IsIntermediaryDelegator(ctx, intermediaryDelegator) {
		k.SetIntermediaryDelegator(ctx, intermediaryDelegator)
	}

	sdkMsg := govv1.MsgVoteWeighted{
		ProposalId: msg.ProposalId,
		Voter:      intermediaryDelegator.String(),
		Options:    msg.Options,
		Metadata:   msg.Metadata,
	}

	_, err = k.govMsgServer.VoteWeighted(ctx, &sdkMsg)
	if err != nil {
		return nil, err
	}
	return &types.MsgVoteWeightedResponse{}, nil
}
