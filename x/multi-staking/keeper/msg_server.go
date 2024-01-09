package keeper

import (
	"context"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type msgServer struct {
	keeper           Keeper
	stakingMsgServer stakingtypes.MsgServer
}

var _ stakingtypes.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) stakingtypes.MsgServer {
	return &msgServer{
		keeper:           keeper,
		stakingMsgServer: stakingkeeper.NewMsgServerImpl(keeper.stakingKeeper),
	}
}

// CreateValidator defines a method for creating a new validator
func (k msgServer) CreateValidator(goCtx context.Context, msg *stakingtypes.MsgCreateValidator) (*stakingtypes.MsgCreateValidatorResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)

	// multiStakerAddr, valAcc, err := types.AccAddrAndValAddrFromStrings(msg.MultiStakerAddress, msg.ValidatorAddress)
	// if err != nil {
	// 	return nil, err
	// }

	// lockID := types.MultiStakingLockID(msg.MultiStakerAddress, msg.ValidatorAddress)

	// mintedBondCoin, err := k.Keeper.LockCoinAndMintBondCoin(ctx, lockID, multiStakerAddr, multiStakerAddr, msg.Value)
	// if err != nil {
	// 	return nil, err
	// }

	// sdkMsg := stakingtypes.MsgCreateValidator{
	// 	Description:       msg.Description,
	// 	Commission:        msg.Commission,
	// 	MinSelfDelegation: msg.MinSelfDelegation,
	// 	DelegatorAddress:  msg.MultiStakerAddress,
	// 	ValidatorAddress:  msg.ValidatorAddress,
	// 	Pubkey:            msg.Pubkey,
	// 	Value:             mintedBondCoin,
	// }

	// k.SetValidatorMultiStakingCoin(ctx, valAcc, msg.Value.Denom)

	// _, err = k.stakingMsgServer.CreateValidator(ctx, &sdkMsg)

	// if err != nil {
	// 	return nil, err
	// }

	return &stakingtypes.MsgCreateValidatorResponse{}, nil
}

// EditValidator defines a method for editing an existing validator
func (k msgServer) EditValidator(goCtx context.Context, msg *stakingtypes.MsgEditValidator) (*stakingtypes.MsgEditValidatorResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)

	// sdkMsg := stakingtypes.MsgEditValidator{
	// 	Description:       msg.Description,
	// 	CommissionRate:    msg.CommissionRate,
	// 	MinSelfDelegation: msg.MinSelfDelegation,
	// 	ValidatorAddress:  msg.ValidatorAddress,
	// }

	// _, err := k.stakingMsgServer.EditValidator(ctx, &sdkMsg)
	// if err != nil {
	// 	return nil, err
	// }
	return &stakingtypes.MsgEditValidatorResponse{}, nil
}

// Delegate defines a method for performing a delegation of coins from a delegator to a validator
func (k msgServer) Delegate(goCtx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)

	// multiStakerAddr, valAcc, err := types.AccAddrAndValAddrFromStrings(msg.MultiStakerAddress, msg.ValidatorAddress)
	// if err != nil {
	// 	return nil, err
	// }

	// if !k.isValMultiStakingCoin(ctx, valAcc, msg.Amount) {
	// 	return nil, fmt.Errorf("not allowed coin")
	// }

	// lockID := types.MultiStakingLockID(msg.MultiStakerAddress, msg.ValidatorAddress)

	// mintedBondCoin, err := k.Keeper.LockCoinAndMintBondCoin(ctx, lockID, multiStakerAddr, multiStakerAddr, msg.Amount)
	// if err != nil {
	// 	return nil, err
	// }

	// sdkMsg := stakingtypes.MsgDelegate{
	// 	DelegatorAddress: msg.MultiStakerAddress,
	// 	ValidatorAddress: msg.ValidatorAddress,
	// 	Amount:           mintedBondCoin,
	// }

	// _, err = k.stakingMsgServer.Delegate(ctx, &sdkMsg)
	// if err != nil {
	// 	return nil, err
	// }

	return &stakingtypes.MsgDelegateResponse{}, nil
}

// BeginRedelegate defines a method for performing a redelegation of coins from a delegator and source validator to a destination validator
func (k msgServer) BeginRedelegate(goCtx context.Context, msg *stakingtypes.MsgBeginRedelegate) (*stakingtypes.MsgBeginRedelegateResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)

	// multiStakerAddr := sdk.MustAccAddressFromBech32(msg.MultiStakerAddress)

	// srcValAcc, err := sdk.ValAddressFromBech32(msg.ValidatorSrcAddress)
	// if err != nil {
	// 	return nil, err
	// }
	// dstValAcc, err := sdk.ValAddressFromBech32(msg.ValidatorDstAddress)
	// if err != nil {
	// 	return nil, err
	// }

	// if !k.isValMultiStakingCoin(ctx, srcValAcc, msg.Amount) || !k.isValMultiStakingCoin(ctx, dstValAcc, msg.Amount) {
	// 	return nil, fmt.Errorf("not allowed Coin")
	// }

	// fromLockID := types.MultiStakingLockID(msg.MultiStakerAddress, msg.ValidatorSrcAddress)
	// fromLock, found := k.GetMultiStakingLock(ctx, fromLockID)
	// if !found {
	// 	return nil, fmt.Errorf("lock not found")
	// }

	// toLockID := types.MultiStakingLockID(msg.MultiStakerAddress, msg.ValidatorDstAddress)
	// toLock := k.GetOrCreateMultiStakingLock(ctx, toLockID)

	// multiStakingCoin := fromLock.MultiStakingCoin(msg.Amount.Amount)

	// err = fromLock.MoveCoinToLock(&toLock, multiStakingCoin)
	// if err != nil {
	// 	return nil, err
	// }
	// k.SetMultiStakingLock(ctx, fromLock)
	// k.SetMultiStakingLock(ctx, toLock)

	// redelegateAmount := multiStakingCoin.BondValue()
	// redelegateAmount, err = k.AdjustUnbondAmount(ctx, multiStakerAddr, srcValAcc, redelegateAmount)
	// if err != nil {
	// 	return nil, err
	// }

	// bondCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), redelegateAmount)

	// sdkMsg := &stakingtypes.MsgBeginRedelegate{
	// 	DelegatorAddress:    msg.MultiStakerAddress,
	// 	ValidatorSrcAddress: msg.ValidatorSrcAddress,
	// 	ValidatorDstAddress: msg.ValidatorDstAddress,
	// 	Amount:              bondCoin,
	// }
	// _, err = k.stakingMsgServer.BeginRedelegate(goCtx, sdkMsg)
	// if err != nil {
	// 	return nil, err
	// }

	return &stakingtypes.MsgBeginRedelegateResponse{}, nil
}

// Undelegate defines a method for performing an undelegation from a delegate and a validator
func (k msgServer) Undelegate(goCtx context.Context, msg *stakingtypes.MsgUndelegate) (*stakingtypes.MsgUndelegateResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)

	// multiStakerAddr, valAcc, err := types.AccAddrAndValAddrFromStrings(msg.MultiStakerAddress, msg.ValidatorAddress)
	// if err != nil {
	// 	return nil, err
	// }

	// if !k.isValMultiStakingCoin(ctx, valAcc, msg.Amount) {
	// 	return nil, fmt.Errorf("not allowed coin")
	// }

	// lockID := types.MultiStakingLockID(msg.MultiStakerAddress, msg.ValidatorAddress)
	// lock, found := k.GetMultiStakingLock(ctx, lockID)
	// if !found {
	// 	return nil, fmt.Errorf("can't find multi staking lock")
	// }

	// multiStakingCoin := lock.MultiStakingCoin(msg.Amount.Amount)
	// err = lock.RemoveCoinFromMultiStakingLock(multiStakingCoin)
	// if err != nil {
	// 	return nil, err
	// }
	// k.SetMultiStakingLock(ctx, lock)

	// unbondAmount := multiStakingCoin.BondValue()
	// unbondAmount, err = k.AdjustUnbondAmount(ctx, multiStakerAddr, valAcc, unbondAmount)
	// if err != nil {
	// 	return nil, err
	// }

	// unbondCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), unbondAmount)

	// sdkMsg := &stakingtypes.MsgUndelegate{
	// 	DelegatorAddress: msg.MultiStakerAddress,
	// 	ValidatorAddress: msg.ValidatorAddress,
	// 	Amount:           unbondCoin,
	// }

	// _, err = k.stakingMsgServer.Undelegate(goCtx, sdkMsg)
	// if err != nil {
	// 	return nil, err
	// }

	// k.SetMultiStakingUnlockEntry(ctx, types.MultiStakingUnlockID(msg.MultiStakerAddress, msg.ValidatorAddress), multiStakingCoin)

	return &stakingtypes.MsgUndelegateResponse{}, nil
}

// // CancelUnbondingDelegation defines a method for canceling the unbonding delegation
// // and delegate back to the validator.
func (k msgServer) CancelUnbondingDelegation(goCtx context.Context, msg *stakingtypes.MsgCancelUnbondingDelegation) (*stakingtypes.MsgCancelUnbondingDelegationResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)

	// multiStakerAddr, valAcc, err := types.AccAddrAndValAddrFromStrings(msg.MultiStakerAddress, msg.ValidatorAddress)
	// if err != nil {
	// 	return nil, err
	// }

	// if !k.isValMultiStakingCoin(ctx, valAcc, msg.Amount) {
	// 	return nil, fmt.Errorf("not allow coin")
	// }

	// unbondEntry, found := k.GetUnbondingEntryAtCreationHeight(ctx, multiStakerAddr, valAcc, msg.CreationHeight)
	// if !found {
	// 	return nil, fmt.Errorf("unbondEntry not found")
	// }
	// cancelUnbondingCoin := sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), unbondEntry.Balance)

	// sdkMsg := &stakingtypes.MsgCancelUnbondingDelegation{
	// 	DelegatorAddress: msg.MultiStakerAddress,
	// 	ValidatorAddress: msg.ValidatorAddress,
	// 	Amount:           cancelUnbondingCoin,
	// }

	// _, err = k.stakingMsgServer.CancelUnbondingDelegation(goCtx, sdkMsg)
	// if err != nil {
	// 	return nil, err
	// }

	// unlockID := types.MultiStakingUnlockID(msg.MultiStakerAddress, msg.ValidatorAddress)
	// err = k.DeleteUnlockEntryAtCreationHeight(ctx, unlockID, msg.CreationHeight)
	// if err != nil {
	// 	return nil, err
	// }

	return &stakingtypes.MsgCancelUnbondingDelegationResponse{}, nil
}
