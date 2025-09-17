package keeper

import (
	"context"
	"fmt"

	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type msgServer struct {
	keeper           Keeper
	stakingMsgServer stakingtypes.MsgServer
}

var (
	_ stakingtypes.MsgServer = msgServer{}
	_ types.MsgServer        = msgServer{}
)

func NewMultiStakingMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

func (k msgServer) UpdateMultiStakingParams(goCtx context.Context, msg *types.MsgUpdateMultiStakingParams) (*types.MsgUpdateMultiStakingParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.keeper.authority != msg.Authority {
		return nil, fmt.Errorf("invalid authority; expected %s, got %s", k.keeper.authority, msg.Authority)
	}

	// store params
	if err := k.keeper.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateMultiStakingParamsResponse{}, nil
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) *msgServer {
	return &msgServer{
		keeper:           keeper,
		stakingMsgServer: stakingkeeper.NewMsgServerImpl(keeper.stakingKeeper),
	}
}

// UpdateParams updates the staking params
func (k msgServer) UpdateParams(goCtx context.Context, msg *stakingtypes.MsgUpdateParams) (*stakingtypes.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.keeper.authority != msg.Authority {
		return nil, fmt.Errorf("invalid authority; expected %s, got %s", k.keeper.authority, msg.Authority)
	}

	// store params
	if err := k.keeper.stakingKeeper.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &stakingtypes.MsgUpdateParamsResponse{}, nil
}

// CreateValidator defines a method for creating a new validator
func (k msgServer) CreateValidator(goCtx context.Context, msg *stakingtypes.MsgCreateValidator) (*stakingtypes.MsgCreateValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	multiStakerAddr, err := k.keeper.validatorAddressCodec.StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	lockID := types.MultiStakingLockID(sdk.AccAddress(multiStakerAddr).String(), msg.ValidatorAddress)

	mintedBondCoin, err := k.keeper.LockCoinAndMintBondCoin(ctx, lockID, multiStakerAddr, multiStakerAddr, msg.Value)
	if err != nil {
		return nil, err
	}

	sdkMsg := stakingtypes.MsgCreateValidator{
		Description:       msg.Description,
		Commission:        msg.Commission,
		MinSelfDelegation: msg.MinSelfDelegation,
		ValidatorAddress:  msg.ValidatorAddress,
		Pubkey:            msg.Pubkey,
		Value:             mintedBondCoin, // replace lock coin with bond coin
	}

	k.keeper.SetValidatorMultiStakingCoin(ctx, sdk.ValAddress(multiStakerAddr), msg.Value.Denom)

	return k.stakingMsgServer.CreateValidator(ctx, &sdkMsg)
}

// EditValidator defines a method for editing an existing validator
func (k msgServer) EditValidator(goCtx context.Context, msg *stakingtypes.MsgEditValidator) (*stakingtypes.MsgEditValidatorResponse, error) {
	return k.stakingMsgServer.EditValidator(goCtx, msg)
}

// Delegate defines a method for performing a delegation of coins from a delegator to a validator
func (k msgServer) Delegate(goCtx context.Context, msg *stakingtypes.MsgDelegate) (*stakingtypes.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	multiStakerAddr, valAcc, err := types.AccAddrAndValAddrFromStrings(msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	if !k.keeper.isValMultiStakingCoin(ctx, valAcc, msg.Amount) {
		return nil, fmt.Errorf("not allowed coin")
	}

	lockID := types.MultiStakingLockID(msg.DelegatorAddress, msg.ValidatorAddress)

	mintedBondCoin, err := k.keeper.LockCoinAndMintBondCoin(ctx, lockID, multiStakerAddr, multiStakerAddr, msg.Amount)
	if err != nil {
		return nil, err
	}

	sdkMsg := stakingtypes.MsgDelegate{
		DelegatorAddress: msg.DelegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           mintedBondCoin, // replace lock coin with bond coin
	}

	return k.stakingMsgServer.Delegate(ctx, &sdkMsg)
}

// BeginRedelegate defines a method for performing a redelegation of coins from a delegator and source validator to a destination validator
func (k msgServer) BeginRedelegate(goCtx context.Context, msg *stakingtypes.MsgBeginRedelegate) (*stakingtypes.MsgBeginRedelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	multiStakerAddr := sdk.MustAccAddressFromBech32(msg.DelegatorAddress)

	srcValAcc, err := sdk.ValAddressFromBech32(msg.ValidatorSrcAddress)
	if err != nil {
		return nil, err
	}
	dstValAcc, err := sdk.ValAddressFromBech32(msg.ValidatorDstAddress)
	if err != nil {
		return nil, err
	}

	if !k.keeper.isValMultiStakingCoin(ctx, srcValAcc, msg.Amount) || !k.keeper.isValMultiStakingCoin(ctx, dstValAcc, msg.Amount) {
		return nil, fmt.Errorf("not allowed Coin")
	}

	fromLockID := types.MultiStakingLockID(msg.DelegatorAddress, msg.ValidatorSrcAddress)
	fromLock, found := k.keeper.GetMultiStakingLock(ctx, fromLockID)
	if !found {
		return nil, fmt.Errorf("lock not found")
	}

	toLockID := types.MultiStakingLockID(msg.DelegatorAddress, msg.ValidatorDstAddress)
	toLock := k.keeper.GetOrCreateMultiStakingLock(ctx, toLockID)

	multiStakingCoin := fromLock.MultiStakingCoin(msg.Amount.Amount)

	err = fromLock.MoveCoinToLock(&toLock, multiStakingCoin)
	if err != nil {
		return nil, err
	}
	k.keeper.SetMultiStakingLock(ctx, fromLock)
	k.keeper.SetMultiStakingLock(ctx, toLock)

	redelegateAmount := multiStakingCoin.BondValue()
	redelegateAmount, err = k.keeper.AdjustUnbondAmount(ctx, multiStakerAddr, srcValAcc, redelegateAmount)
	if err != nil {
		return nil, err
	}

	bondDenom, err := k.keeper.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}
	bondCoin := sdk.NewCoin(bondDenom, redelegateAmount)

	sdkMsg := &stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    msg.DelegatorAddress,
		ValidatorSrcAddress: msg.ValidatorSrcAddress,
		ValidatorDstAddress: msg.ValidatorDstAddress,
		Amount:              bondCoin, // replace lockCoin with bondCoin
	}

	return k.stakingMsgServer.BeginRedelegate(ctx, sdkMsg)
}

// Undelegate defines a method for performing an undelegation from a delegate and a validator
func (k msgServer) Undelegate(goCtx context.Context, msg *stakingtypes.MsgUndelegate) (*stakingtypes.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	multiStakerAddr, valAcc, err := types.AccAddrAndValAddrFromStrings(msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	if !k.keeper.isValMultiStakingCoin(ctx, valAcc, msg.Amount) {
		return nil, fmt.Errorf("not allowed coin")
	}

	lockID := types.MultiStakingLockID(msg.DelegatorAddress, msg.ValidatorAddress)
	lock, found := k.keeper.GetMultiStakingLock(ctx, lockID)
	if !found {
		return nil, fmt.Errorf("can't find multi staking lock")
	}

	multiStakingCoin := lock.MultiStakingCoin(msg.Amount.Amount)
	err = lock.RemoveCoinFromMultiStakingLock(multiStakingCoin)
	if err != nil {
		return nil, err
	}
	k.keeper.SetMultiStakingLock(ctx, lock)

	unbondAmount := multiStakingCoin.BondValue()
	unbondAmount, err = k.keeper.AdjustUnbondAmount(ctx, multiStakerAddr, valAcc, unbondAmount)
	if err != nil {
		return nil, err
	}

	bondDenom, err := k.keeper.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}
	unbondCoin := sdk.NewCoin(bondDenom, unbondAmount)

	sdkMsg := &stakingtypes.MsgUndelegate{
		DelegatorAddress: msg.DelegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           unbondCoin, // replace with unbondCoin
	}

	k.keeper.SetMultiStakingUnlockEntry(ctx, types.MultiStakingUnlockID(msg.DelegatorAddress, msg.ValidatorAddress), multiStakingCoin)

	return k.stakingMsgServer.Undelegate(ctx, sdkMsg)
}

// CancelUnbondingDelegation defines a method for canceling the unbonding delegation
// and delegate back to the validator.
func (k msgServer) CancelUnbondingDelegation(goCtx context.Context, msg *stakingtypes.MsgCancelUnbondingDelegation) (*stakingtypes.MsgCancelUnbondingDelegationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	delAcc, valAcc, err := types.AccAddrAndValAddrFromStrings(msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	// prevent cancel unbonding for non multi-staking coin
	// incase denom was removed
	_, found := k.keeper.GetBondWeight(ctx, msg.Amount.Denom)
	if !found {
		return nil, fmt.Errorf("error MultiStakingCoin %s not found", msg.Amount.Denom)
	}

	if !k.keeper.isValMultiStakingCoin(ctx, valAcc, msg.Amount) {
		return nil, fmt.Errorf("not allow coin")
	}

	unlockID := types.MultiStakingUnlockID(msg.DelegatorAddress, msg.ValidatorAddress)
	cancelUnlockingCoin, err := k.keeper.DecreaseUnlockEntryAmount(ctx, unlockID, msg.Amount.Amount, msg.CreationHeight)
	if err != nil {
		return nil, err
	}

	cancelUnbondingAmount := cancelUnlockingCoin.BondValue()
	cancelUnbondingAmount, err = k.keeper.AdjustCancelUnbondingAmount(ctx, delAcc, valAcc, msg.CreationHeight, cancelUnbondingAmount)
	if err != nil {
		return nil, err
	}
	bondDenom, err := k.keeper.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}
	cancelUnbondingCoin := sdk.NewCoin(bondDenom, cancelUnbondingAmount)

	lockID := types.MultiStakingLockID(msg.DelegatorAddress, msg.ValidatorAddress)
	lock := k.keeper.GetOrCreateMultiStakingLock(ctx, lockID)
	err = lock.AddCoinToMultiStakingLock(cancelUnlockingCoin)
	if err != nil {
		return nil, err
	}
	k.keeper.SetMultiStakingLock(ctx, lock)

	sdkMsg := &stakingtypes.MsgCancelUnbondingDelegation{
		DelegatorAddress: msg.DelegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           cancelUnbondingCoin, // replace with cancelUnbondingCoin
		CreationHeight:   msg.CreationHeight,
	}

	return k.stakingMsgServer.CancelUnbondingDelegation(ctx, sdkMsg)
}

// CreateEVMValidator defines a method for creating a new validator using ERC20 token
func (k msgServer) CreateEVMValidator(goCtx context.Context, msg *types.MsgCreateEVMValidator) (*types.MsgCreateEVMValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Converting ERC20 to cosmos coin
	multiStakerAddr, err := k.keeper.validatorAddressCodec.StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	multiStakerHex := common.BytesToAddress(multiStakerAddr).Hex()

	_, err = k.keeper.erc20keeper.ConvertERC20(goCtx, &erc20types.MsgConvertERC20{
		Sender:          multiStakerHex,
		ContractAddress: msg.ContractAddress,
		Amount:          msg.Value,
		Receiver:        sdk.AccAddress(multiStakerAddr).String(),
	})
	if err != nil {
		return nil, err
	}

	tokenDenom, err := k.keeper.erc20keeper.GetTokenDenom(ctx, common.HexToAddress(msg.ContractAddress))
	if err != nil {
		return nil, err
	}

	createValMsg := &stakingtypes.MsgCreateValidator{
		Description:       msg.Description,
		Commission:        msg.Commission,
		MinSelfDelegation: msg.MinSelfDelegation,
		ValidatorAddress:  msg.ValidatorAddress,
		Pubkey:            msg.Pubkey,
		Value:             sdk.NewCoin(tokenDenom, msg.Value),
	}

	_, err = k.CreateValidator(goCtx, createValMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateEVMValidatorResponse{}, nil
}

// Delegate defines a method for performing a delegation of coins from a delegator to a validator
func (k msgServer) DelegateEVM(goCtx context.Context, msg *types.MsgDelegateEVM) (*types.MsgDelegateEVMResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	delAddr := sdk.MustAccAddressFromBech32(msg.DelegatorAddress)
	delHex := common.BytesToAddress(delAddr).Hex()

	_, err := k.keeper.erc20keeper.ConvertERC20(goCtx, &erc20types.MsgConvertERC20{
		Sender:          delHex,
		ContractAddress: msg.ContractAddress,
		Amount:          msg.Amount,
		Receiver:        msg.DelegatorAddress,
	})
	if err != nil {
		return nil, err
	}

	tokenDenom, err := k.keeper.erc20keeper.GetTokenDenom(ctx, common.HexToAddress(msg.ContractAddress))
	if err != nil {
		return nil, err
	}

	delMsg := &stakingtypes.MsgDelegate{
		DelegatorAddress: msg.DelegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           sdk.NewCoin(tokenDenom, msg.Amount),
	}

	_, err = k.Delegate(goCtx, delMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgDelegateEVMResponse{}, nil
}

// BeginRedelegate defines a method for performing a redelegation of coins from a delegator and source validator to a destination validator
func (k msgServer) BeginRedelegateEVM(goCtx context.Context, msg *types.MsgBeginRedelegateEVM) (*types.MsgBeginRedelegateEVMResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	tokenDenom, err := k.keeper.erc20keeper.GetTokenDenom(ctx, common.HexToAddress(msg.ContractAddress))
	if err != nil {
		return nil, err
	}

	redelegateMsg := &stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    msg.DelegatorAddress,
		ValidatorSrcAddress: msg.ValidatorSrcAddress,
		ValidatorDstAddress: msg.ValidatorDstAddress,
		Amount:              sdk.NewCoin(tokenDenom, msg.Amount),
	}

	res, err := k.BeginRedelegate(goCtx, redelegateMsg)
	if err != nil {
		return nil, err
	}
	return &types.MsgBeginRedelegateEVMResponse{CompletionTime: res.CompletionTime}, nil
}

// Undelegate defines a method for performing an undelegation from a delegate and a validator
func (k msgServer) UndelegateEVM(goCtx context.Context, msg *types.MsgUndelegateEVM) (*types.MsgUndelegateEVMResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	tokenDenom, err := k.keeper.erc20keeper.GetTokenDenom(ctx, common.HexToAddress(msg.ContractAddress))
	if err != nil {
		return nil, err
	}

	undelegateMsg := &stakingtypes.MsgUndelegate{
		DelegatorAddress: msg.DelegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           sdk.NewCoin(tokenDenom, msg.Amount),
	}

	res, err := k.Undelegate(goCtx, undelegateMsg)
	if err != nil {
		return nil, err
	}
	return &types.MsgUndelegateEVMResponse{CompletionTime: res.CompletionTime, ContractAddress: msg.ContractAddress, Amount: msg.Amount}, nil
}

// CancelUnbondingDelegation defines a method for canceling the unbonding delegation
// and delegate back to the validator.
func (k msgServer) CancelUnbondingEVMDelegation(goCtx context.Context, msg *types.MsgCancelUnbondingEVMDelegation) (*types.MsgCancelUnbondingEVMDelegationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	tokenDenom, err := k.keeper.erc20keeper.GetTokenDenom(ctx, common.HexToAddress(msg.ContractAddress))
	if err != nil {
		return nil, err
	}

	cancelUnbondingMsg := &stakingtypes.MsgCancelUnbondingDelegation{
		DelegatorAddress: msg.DelegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           sdk.NewCoin(tokenDenom, msg.Amount),
		CreationHeight:   msg.CreationHeight,
	}

	_, err = k.CancelUnbondingDelegation(goCtx, cancelUnbondingMsg)
	if err != nil {
		return nil, err
	}
	return &types.MsgCancelUnbondingEVMDelegationResponse{}, nil
}
