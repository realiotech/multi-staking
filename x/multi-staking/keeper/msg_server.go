package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
// func NewMsgServerImpl(keeper Keeper) types.MsgServer {
// 	return &msgServer{Keeper: keeper}
// }

// var _ types.MsgServer = msgServer{}

// CreateValidator defines a method for creating a new validator
func (k msgServer) CreateValidator(goCtx context.Context, msg *types.MsgCreateValidator) (*types.MsgCreateValidatorResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	valAcc, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	delAcc := sdk.MustAccAddressFromBech32(msg.DelegatorAddress)

	if !k.IsAllowedToken(ctx, valAcc, msg.Value) {
		return nil, fmt.Errorf("not allowed token")
	}

	intermediaryAccount := k.GetIntermediaryAccountDelegator(ctx, delAcc)
	if intermediaryAccount == nil {
		intermediaryAccount = types.IntermediaryAccount(msg.DelegatorAddress)
		k.SetIntermediaryAccountDelegator(ctx, intermediaryAccount, delAcc)
		k.SetDelAddrByKeyIntermediaryAccount(ctx, intermediaryAccount, delAcc)
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

	valAcc, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	delAcc := sdk.MustAccAddressFromBech32(msg.DelegatorAddress)

	if !k.IsAllowedToken(ctx, valAcc, msg.Amount) {
		return nil, fmt.Errorf("not allowed token")
	}

	intermediaryAccount := k.GetIntermediaryAccountDelegator(ctx, delAcc)
	if intermediaryAccount == nil {
		intermediaryAccount = types.IntermediaryAccount(msg.DelegatorAddress)
		k.SetIntermediaryAccountDelegator(ctx, intermediaryAccount, delAcc)
		k.SetDelAddrByKeyIntermediaryAccount(ctx, intermediaryAccount, delAcc)
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

func (k Keeper) MoveLockedMultistakingToken(ctx sdk.Context, srcLockID []byte, dstLockID []byte, lockedToken sdk.Coin) (err error) {
	// get lock on source val
	srcLock, found := k.GetMultiStakingLock(ctx, srcLockID)
	if !found {
		return fmt.Errorf("can't find multi staking lock")
	}

	// remove token from lock on source val
	srcLock, err = srcLock.RemoveTokenFromMultiStakingLock(lockedToken.Amount)
	if err != nil {
		return err
	}

	// update lock on source val
	k.SetMultiStakingLock(ctx, srcLockID, srcLock)

	// get lock on destination val
	dstLock, found := k.GetMultiStakingLock(ctx, dstLockID)
	if !found {
		dstLock = types.NewMultiStakingLock(lockedToken.Amount, srcLock.ConversionRatio)
		k.SetMultiStakingLock(ctx, dstLockID, dstLock)
	} else {
		dstLock = dstLock.AddTokenToMultiStakingLock(lockedToken.Amount, srcLock.ConversionRatio)
	}

	// update lock on destination val
	k.SetMultiStakingLock(ctx, dstLockID, dstLock)

	return err
}

func (k Keeper) LockedAmountToBondAmount(ctx sdk.Context, multiStakingLockID []byte, lockedAmount math.Int) (math.Int, error) {
	// get lock on source val
	lock, found := k.GetMultiStakingLock(ctx, multiStakingLockID)
	if !found {
		return math.Int{}, fmt.Errorf("can't find multi staking lock")
	}

	return lock.LockedAmountToBondAmount(lockedAmount).RoundInt(), nil
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

	srcLockID := types.MultiStakingLockID(delAcc, srcValAcc)
	bondAmount, err := k.LockedAmountToBondAmount(ctx, srcLockID, msg.Amount.Amount)
	if err != nil {
		return nil, err
	}

	sdkBondCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), bondAmount)
	sdkMsg := &stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    msg.DelegatorAddress,
		ValidatorSrcAddress: msg.ValidatorSrcAddress,
		ValidatorDstAddress: msg.ValidatorDstAddress,
		Amount:              sdkBondCoin,
	}
	_, err = k.stakingMsgServer.BeginRedelegate(goCtx, sdkMsg)
	if err != nil {
		return nil, err
	}

	dstLockID := types.MultiStakingLockID(delAcc, dstValAcc)
	err = k.MoveLockedMultistakingToken(ctx, srcLockID, dstLockID, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgBeginRedelegateResponse{}, err
}

// Undelegate defines a method for performing an undelegation from a delegate and a validator
func (k msgServer) Undelegate(goCtx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	delAcc := sdk.MustAccAddressFromBech32(msg.DelegatorAddress)

	valAcc, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	lockID := types.MultiStakingLockID(delAcc, valAcc)
	unbondAmount, err := k.LockedAmountToBondAmount(ctx, lockID, msg.Amount.Amount)
	if err != nil {
		return nil, err
	}
	unbondCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), unbondAmount)

	sdkMsg := &stakingtypes.MsgUndelegate{
		DelegatorAddress: msg.DelegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           unbondCoin,
	}

	_, err = k.stakingMsgServer.Undelegate(goCtx, sdkMsg)

	return &types.MsgUndelegateResponse{}, err
}

// // CancelUnbondingDelegation defines a method for canceling the unbonding delegation
// // and delegate back to the validator.
// func (k msgServer) CancelUnbondingDelegation(goCtx context.Context, msg *types.MsgCancelUnbondingDelegation) (*types.MsgCancelUnbondingDelegationResponse, error) {
// 	ctx := sdk.UnwrapSDKContext(goCtx)

// 	intermediaryAccount := types.GetIntermediaryAccount(msg.DelegatorAddress, msg.ValidatorAddress)

// 	valAcc, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
// 	if err != nil {
// 		return nil, err
// 	}
// 	delAcc := sdk.MustAccAddressFromBech32(msg.DelegatorAddress)

// 	sdkMsg := stakingtypes.MsgCancelUnbondingDelegation{
// 		DelegatorAddress: intermediaryAccount.String(),
// 		ValidatorAddress: msg.ValidatorAddress,
// 		Amount:           exactDelegateValue,
// 	}

// 	k.Keeper.PreDelegate(ctx, delAcc, valAcc, msg.Amount)

// 	_, err = k.stakingMsgServer.CancelUnbondingDelegation(ctx, &sdkMsg)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &types.MsgCancelUnbondingDelegationResponse{}, nil
// }
