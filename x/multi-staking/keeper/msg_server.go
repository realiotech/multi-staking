package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

	delAcc, valAcc, err := types.DelAccAndValAccFromStrings(msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	intermediaryAccount := types.IntermediaryAccount(delAcc)
	if !k.IsIntermediaryAccount(ctx, intermediaryAccount) {
		k.SetIntermediaryAccount(ctx, intermediaryAccount)
	}

	sdkBondToken, err := k.Keeper.LockMultiStakingTokenAndMintBondToken(ctx, delAcc, valAcc, msg.Value)
	if err != nil {
		return nil, err
	}

	sdkMsg := stakingtypes.MsgCreateValidator{
		Description:       msg.Description,
		Commission:        msg.Commission,
		MinSelfDelegation: msg.MinSelfDelegation,
		DelegatorAddress:  intermediaryAccount.String(),
		ValidatorAddress:  msg.ValidatorAddress,
		Pubkey:            msg.Pubkey,
		Value:             sdkBondToken,
	}

	k.SetValidatorAllowedToken(ctx, valAcc, msg.Value.Denom)

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

	delAcc, valAcc, err := types.DelAccAndValAccFromStrings(msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	if !k.IsAllowedToken(ctx, valAcc, msg.Amount) {
		return nil, fmt.Errorf("not allowed token")
	}

	intermediaryAccount := types.IntermediaryAccount(delAcc)
	if !k.IsIntermediaryAccount(ctx, intermediaryAccount) {
		k.SetIntermediaryAccount(ctx, intermediaryAccount)
	}

	mintedBondToken, err := k.Keeper.LockMultiStakingTokenAndMintBondToken(ctx, delAcc, valAcc, msg.Amount)
	if err != nil {
		return nil, err
	}

	sdkMsg := stakingtypes.MsgDelegate{
		DelegatorAddress: intermediaryAccount.String(),
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           mintedBondToken,
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

	delAcc := sdk.MustAccAddressFromBech32(msg.DelegatorAddress)

	srcValAcc, err := sdk.ValAddressFromBech32(msg.ValidatorSrcAddress)
	if err != nil {
		return nil, err
	}
	dstValAcc, err := sdk.ValAddressFromBech32(msg.ValidatorDstAddress)
	if err != nil {
		return nil, err
	}

	if !k.IsAllowedToken(ctx, srcValAcc, msg.Amount) || !k.IsAllowedToken(ctx, dstValAcc, msg.Amount) {
		return nil, fmt.Errorf("not allowed Token")
	}

	bondAmount, err := k.LockedAmountToBondAmount(ctx, delAcc, srcValAcc, msg.Amount.Amount)
	if err != nil {
		return nil, err
	}

	sdkBondCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), bondAmount)
	intermediaryAccount := types.IntermediaryAccount(delAcc)

	sdkMsg := &stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    intermediaryAccount.String(),
		ValidatorSrcAddress: msg.ValidatorSrcAddress,
		ValidatorDstAddress: msg.ValidatorDstAddress,
		Amount:              sdkBondCoin,
	}
	_, err = k.stakingMsgServer.BeginRedelegate(goCtx, sdkMsg)
	if err != nil {
		return nil, err
	}

	err = k.MoveLockedMultistakingToken(ctx, delAcc, srcValAcc, dstValAcc, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgBeginRedelegateResponse{}, err
}

func (k Keeper) GetDelegation(ctx sdk.Context, delAcc sdk.AccAddress, val sdk.ValAddress) (stakingtypes.Delegation, bool) {
	return k.stakingKeeper.GetDelegation(ctx, types.IntermediaryAccount(delAcc), val)
}

func (k Keeper) AdjustUnbondAmount(ctx sdk.Context, delAcc sdk.AccAddress, valAcc sdk.ValAddress, amount math.Int) (adjustedAmount math.Int, err error) {
	delegation, found := k.GetDelegation(ctx, delAcc, valAcc)
	if !found {
		return math.Int{}, fmt.Errorf("delegation not found")
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, valAcc)
	if !found {
		return math.Int{}, fmt.Errorf("validator not found")
	}

	shares, err := validator.SharesFromTokens(amount)
	if err != nil {
		return math.Int{}, err
	}

	delShares := delegation.GetShares()
	// Cap the shares at the delegation's shares. Shares being greater could occur
	// due to rounding, however we don't want to truncate the shares or take the
	// minimum because we want to allow for the full withdraw of shares from a
	// delegation.
	if shares.GT(delShares) {
		shares = delShares
	}

	return validator.TokensFromShares(shares).RoundInt(), nil
}

// Undelegate defines a method for performing an undelegation from a delegate and a validator
func (k msgServer) Undelegate(goCtx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	delAcc := sdk.MustAccAddressFromBech32(msg.DelegatorAddress)

	delAcc, valAcc, err := types.DelAccAndValAccFromStrings(msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	if !k.IsAllowedToken(ctx, valAcc, msg.Amount) {
		return nil, fmt.Errorf("not allowed token")
	}

	lockID := types.MultiStakingLockID(delAcc, valAcc)
	lock, found := k.GetMultiStakingLock(ctx, lockID)
	if !found {
		return nil, fmt.Errorf("can't find multi staking lock")
	}

	unbondAmount := lock.LockedAmountToBondAmount(msg.Amount.Amount).RoundInt()

	unbondCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), unbondAmount)
	intermediaryAccount := types.IntermediaryAccount(delAcc)

	sdkMsg := &stakingtypes.MsgUndelegate{
		DelegatorAddress: intermediaryAccount.String(),
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           unbondCoin,
	}

	resp, err := k.stakingMsgServer.Undelegate(goCtx, sdkMsg)

	k.RemoveTokenFromLock(ctx, delAcc, valAcc, msg.Amount.Amount)

	k.SetMultiStakingUnlockEntry(ctx, delAcc, valAcc, ctx.BlockHeight(), lock.ConversionRatio, resp.CompletionTime, msg.Amount.Amount)

	return &types.MsgUndelegateResponse{}, err
}

// CancelUnbondingDelegation defines a method for canceling the unbonding delegation
// and delegate back to the validator.
func (k msgServer) CancelUnbondingDelegation(goCtx context.Context, msg *types.MsgCancelUnbondingDelegation) (*types.MsgCancelUnbondingDelegationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	delAcc, valAcc, err := types.DelAccAndValAccFromStrings(msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	intermediaryAccount := types.IntermediaryAccount(delAcc)

	valDenom := k.GetValidatorAllowedToken(ctx, valAcc)

	if msg.Amount.Denom != valDenom {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, valDenom,
		)
	}

	ubd, found := k.GetMultiStakingUnlock(ctx, delAcc, valAcc)

	if !found {
		return nil, fmt.Errorf("not found unbonding recored")
	}

	var (
		unbondEntry      types.UnlockEntry
		unbondEntryIndex int64 = -1
	)

	for i, entry := range ubd.Entries {
		if entry.CreationHeight == msg.CreationHeight {
			unbondEntry = entry
			unbondEntryIndex = int64(i)
			break
		}
	}
	if unbondEntryIndex == -1 {
		return nil, sdkerrors.ErrNotFound.Wrapf("unbonding delegation entry is not found at block height %d", msg.CreationHeight)
	}

	if unbondEntry.Balance.LT(msg.Amount.Amount) {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("amount is greater than the unbonding delegation entry balance")
	}

	// delegate back the unbonding delegation amount to the validator
	// _, err = k.Keeper.Delegate(ctx, delegatorAddress, msg.Amount.Amount, types.Unbonding, validator, false)
	if err != nil {
		return nil, err
	}

	cancelAmt := unbondEntry.ConversionRatio.MulInt(msg.Amount.Amount).RoundInt()
	if err != nil {
		return nil, err
	}

	cancelCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), cancelAmt)

	sdkMsg := &stakingtypes.MsgCancelUnbondingDelegation{
		DelegatorAddress: intermediaryAccount.String(),
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           cancelCoin,
	}

	_, err = k.stakingMsgServer.CancelUnbondingDelegation(goCtx, sdkMsg)
	if err != nil {
		return nil, err
	}

	amount := unbondEntry.Balance.Sub(msg.Amount.Amount)
	if amount.IsZero() {
		ubd.RemoveEntry(unbondEntryIndex)
	} else {
		// update the unbondingDelegationEntryBalance and InitialBalance for ubd entry
		unbondEntry.Balance = amount
		ubd.Entries[unbondEntryIndex] = unbondEntry
	}

	// set the unbonding delegation or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		k.RemoveMultiStakingUnlock(ctx, ubd)
	} else {
		k.SetMultiStakingUnlock(ctx, ubd)
	}

	mintedBondToken, err := k.Keeper.LockMultiStakingTokenAndMintBondToken(ctx, delAcc, valAcc, msg.Amount)
	if err != nil {
		return nil, err
	}

	delegateMsg := stakingtypes.MsgDelegate{
		DelegatorAddress: intermediaryAccount.String(),
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           mintedBondToken,
	}

	_, err = k.stakingMsgServer.Delegate(ctx, &delegateMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgCancelUnbondingDelegationResponse{}, nil
}

// SetWithdrawAddress defines a method for performing an undelegation from a delegate and a validator
func (k msgServer) SetWithdrawAddress(goCtx context.Context, msg *types.MsgSetWithdrawAddress) (*types.MsgSetWithdrawAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	delAcc, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	intermediaryAccount := types.IntermediaryAccount(delAcc)

	sdkMsg := distrtypes.MsgSetWithdrawAddress{
		DelegatorAddress: intermediaryAccount.String(),
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
	delAcc, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	intermediaryAccount := types.IntermediaryAccount(delAcc)

	if !k.IsIntermediaryAccount(ctx, intermediaryAccount) {
		k.SetIntermediaryAccount(ctx, intermediaryAccount)
	}

	sdkMsg := distrtypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: intermediaryAccount.String(),
		ValidatorAddress: msg.ValidatorAddress,
	}

	resp, err := k.distrMsgServer.WithdrawDelegatorReward(ctx, &sdkMsg)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.SendCoins(ctx, intermediaryAccount, delAcc, resp.Amount)
	if err != nil {
		return nil, err
	}
	return &types.MsgWithdrawDelegatorRewardResponse{Amount: resp.Amount}, nil
}

func (k msgServer) Vote(goCtx context.Context, msg *types.MsgVote) (*types.MsgVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	delAcc, err := sdk.AccAddressFromBech32(msg.Voter)
	if err != nil {
		return nil, err
	}
	intermediaryAcc := types.IntermediaryAccount(delAcc)
	if !k.IsIntermediaryAccount(ctx, intermediaryAcc) {
		k.SetIntermediaryAccount(ctx, intermediaryAcc)
	}

	sdkMsg := govv1.MsgVote{
		ProposalId: msg.ProposalId,
		Voter:      intermediaryAcc.String(),
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
	delAcc, err := sdk.AccAddressFromBech32(msg.Voter)
	if err != nil {
		return nil, err
	}

	intermediaryAcc := types.IntermediaryAccount(delAcc)
	if !k.IsIntermediaryAccount(ctx, intermediaryAcc) {
		k.SetIntermediaryAccount(ctx, intermediaryAcc)
	}

	sdkMsg := govv1.MsgVoteWeighted{
		ProposalId: msg.ProposalId,
		Voter:      intermediaryAcc.String(),
		Options:    msg.Options,
		Metadata:   msg.Metadata,
	}

	_, err = k.govMsgServer.VoteWeighted(ctx, &sdkMsg)
	if err != nil {
		return nil, err
	}
	return &types.MsgVoteWeightedResponse{}, nil
}
