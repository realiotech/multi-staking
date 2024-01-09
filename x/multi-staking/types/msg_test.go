package types_test

import (
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	// govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	// govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/realio-tech/multi-staking-module/testutil"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

// MsgBeginRedelegate
func TestMsgBeginRedelegate_ValidateBasic(t *testing.T) {
	mulStakeAddr := testutil.GenAddress()
	valSrcAddr := testutil.GenValAddress()
	valDstAddr := testutil.GenValAddress()
	denom := "ario"
	coin := sdk.NewCoin(denom, sdk.NewInt(10000))

	tests := []struct {
		name string
		msg  types.MsgBeginRedelegate
		err  error
	}{
		{
			name: "happy path",
			msg: types.MsgBeginRedelegate{
				MultiStakerAddress:  mulStakeAddr.String(),
				ValidatorSrcAddress: valSrcAddr.String(),
				ValidatorDstAddress: valDstAddr.String(),
				Amount:              coin,
			},
		},
		{
			name: "invalid address",
			msg: types.MsgBeginRedelegate{
				MultiStakerAddress:  "",
				ValidatorSrcAddress: "",
				ValidatorDstAddress: "",
				Amount:              coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid shares amount",
			msg: types.MsgBeginRedelegate{
				MultiStakerAddress:  mulStakeAddr.String(),
				ValidatorSrcAddress: valSrcAddr.String(),
				ValidatorDstAddress: valDstAddr.String(),
				Amount:              sdk.NewCoin(denom, sdk.ZeroInt()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

// MsgCancelUnbonding
func TestMsgCancelUnbondingDelegation_ValidateBasic(t *testing.T) {
	mulStakeAddr := testutil.GenAddress()
	valAddr := testutil.GenValAddress()
	denom := "ario"
	coin := sdk.NewCoin(denom, sdk.NewInt(10000))

	tests := []struct {
		name string
		msg  types.MsgCancelUnbonding
		err  error
	}{
		{
			name: "happy path",
			msg: types.MsgCancelUnbonding{
				MultiStakerAddress: mulStakeAddr.String(),
				ValidatorAddress:   valAddr.String(),
				Amount:             coin,
				CreationHeight:     1,
			},
		},
		{
			name: "invalid address",
			msg: types.MsgCancelUnbonding{
				MultiStakerAddress: "",
				ValidatorAddress:   "",
				Amount:             coin,
				CreationHeight:     1,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid height",
			msg: types.MsgCancelUnbonding{
				MultiStakerAddress: mulStakeAddr.String(),
				ValidatorAddress:   valAddr.String(),
				Amount:             coin,
				CreationHeight:     0,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid amount",
			msg: types.MsgCancelUnbonding{
				MultiStakerAddress: mulStakeAddr.String(),
				ValidatorAddress:   valAddr.String(),
				Amount:             sdk.NewCoin(denom, sdk.ZeroInt()),
				CreationHeight:     1,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

// MsgCreateValidator
func TestMsgCreateValidator_ValidateBasic(t *testing.T) {
	valPubKey := testutil.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())
	mulStakeAddr := sdk.AccAddress(valAddr)
	pubKey := codectypes.UnsafePackAny(valPubKey)
	denom := "ario"
	coin := sdk.NewCoin(denom, sdk.NewInt(10000))

	tests := []struct {
		name string
		msg  types.MsgCreateValidator
		err  error
	}{
		{
			name: "happy path",
			msg: types.MsgCreateValidator{
				MultiStakerAddress: mulStakeAddr.String(),
				ValidatorAddress:   valAddr.String(),
				Value:              coin,
				Pubkey:             pubKey,
				Description: stakingtypes.Description{
					Website: "https://validator.cosmos",
					Details: "Test validator",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          sdk.OneDec(),
					MaxRate:       sdk.OneDec(),
					MaxChangeRate: sdk.OneDec(),
				},
				MinSelfDelegation: sdk.NewInt(1),
			},
		},
		{
			name: "invalid address",
			msg: types.MsgCreateValidator{
				MultiStakerAddress: "",
				ValidatorAddress:   "",
				Value:              coin,
				Pubkey:             pubKey,
				Description: stakingtypes.Description{
					Website: "https://validator.cosmos",
					Details: "Test validator",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          sdk.OneDec(),
					MaxRate:       sdk.OneDec(),
					MaxChangeRate: sdk.OneDec(),
				},
				MinSelfDelegation: sdk.NewInt(1),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "validator address is invalid",
			msg: types.MsgCreateValidator{
				MultiStakerAddress: testutil.GenAddress().String(),
				ValidatorAddress:   valAddr.String(),
				Value:              coin,
				Pubkey:             pubKey,
				Description: stakingtypes.Description{
					Website: "https://validator.cosmos",
					Details: "Test validator",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          sdk.OneDec(),
					MaxRate:       sdk.OneDec(),
					MaxChangeRate: sdk.OneDec(),
				},
				MinSelfDelegation: sdk.NewInt(1),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "minimum self delegation must be a positive integer",
			msg: types.MsgCreateValidator{
				MultiStakerAddress: testutil.GenAddress().String(),
				ValidatorAddress:   valAddr.String(),
				Value:              coin,
				Pubkey:             pubKey,
				Description: stakingtypes.Description{
					Website: "https://validator.cosmos",
					Details: "Test validator",
				},
				Commission: stakingtypes.CommissionRates{
					Rate:          sdk.OneDec(),
					MaxRate:       sdk.OneDec(),
					MaxChangeRate: sdk.OneDec(),
				},
				MinSelfDelegation: sdk.NewInt(0),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "empty commission, description",
			msg: types.MsgCreateValidator{
				MultiStakerAddress: mulStakeAddr.String(),
				ValidatorAddress:   valAddr.String(),
				Value:              coin,
				Pubkey:             pubKey,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

// MsgDelegate
func TestMsgMsgDelegate_ValidateBasic(t *testing.T) {
	mulStakeAddr := testutil.GenAddress()
	valAddr := testutil.GenValAddress()
	denom := "ario"
	coin := sdk.NewCoin(denom, sdk.NewInt(10000))

	tests := []struct {
		name string
		msg  types.MsgDelegate
		err  error
	}{
		{
			name: "happy path",
			msg: types.MsgDelegate{
				MultiStakerAddress: mulStakeAddr.String(),
				ValidatorAddress:   valAddr.String(),
				Amount:             coin,
			},
		},
		{
			name: "empty address string is not allowed",
			msg: types.MsgDelegate{
				MultiStakerAddress: "",
				ValidatorAddress:   "",
				Amount:             coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid delegation amount",
			msg: types.MsgDelegate{
				MultiStakerAddress: mulStakeAddr.String(),
				ValidatorAddress:   valAddr.String(),
				Amount:             sdk.NewCoin(denom, sdk.ZeroInt()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

// MsgEditValidator
func TestMsgEditValidator_ValidateBasic(t *testing.T) {
	valAddr := testutil.GenValAddress()
	commissionRate := sdk.NewDec(1)
	minselfDelegation := sdk.NewInt(1)

	tests := []struct {
		name string
		msg  types.MsgEditValidator
		err  error
	}{
		{
			name: "happy path",
			msg: types.MsgEditValidator{
				Description: stakingtypes.Description{
					Website: "https://validator.cosmos",
					Details: "Test validator",
				},
				ValidatorAddress:  valAddr.String(),
				CommissionRate:    &commissionRate,
				MinSelfDelegation: &minselfDelegation,
			},
		},
		{
			name: "invalid validator address",
			msg: types.MsgEditValidator{
				Description: stakingtypes.Description{
					Website: "https://validator.cosmos",
					Details: "Test validator",
				},
				ValidatorAddress:  "",
				CommissionRate:    &commissionRate,
				MinSelfDelegation: &minselfDelegation,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid commissionrate and minselfdelegation",
			msg: types.MsgEditValidator{
				Description:       stakingtypes.Description{},
				ValidatorAddress:  valAddr.String(),
				CommissionRate:    nil,
				MinSelfDelegation: nil,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

// // MsgSetWithdrawAddress
// func TestMsgSetWithdrawAddress_ValidateBasic(t *testing.T) {
// 	mulStakeAddr := testutil.GenAddress()
// 	withdrawAddr := testutil.GenAddress()

// 	tests := []struct {
// 		name string
// 		msg  types.MsgSetWithdrawAddress
// 		err  error
// 	}{
// 		{
// 			name: "happy path",
// 			msg: types.MsgSetWithdrawAddress{
// 				MultiStakerAddress: mulStakeAddr.String(),
// 				WithdrawAddress:    withdrawAddr.String(),
// 			},
// 		},
// 		{
// 			name: "invalid address",
// 			msg: types.MsgSetWithdrawAddress{
// 				MultiStakerAddress: "",
// 				WithdrawAddress:    "",
// 			},
// 			err: sdkerrors.ErrInvalidAddress,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := tt.msg.ValidateBasic()
// 			if tt.err != nil {
// 				require.ErrorIs(t, err, tt.err)
// 				return
// 			}
// 			require.NoError(t, err)
// 		})
// 	}
// }

// MsgUndelegate
func TestMsgUndelegate_ValidateBasic(t *testing.T) {
	mulStakeAddr := testutil.GenAddress()
	valAddr := testutil.GenValAddress()
	denom := "ario"
	coin := sdk.NewCoin(denom, sdk.NewInt(10000))

	tests := []struct {
		name string
		msg  types.MsgUndelegate
		err  error
	}{
		{
			name: "happy path",
			msg: types.MsgUndelegate{
				MultiStakerAddress: mulStakeAddr.String(),
				ValidatorAddress:   valAddr.String(),
				Amount:             coin,
			},
		},
		{
			name: "invalid address",
			msg: types.MsgUndelegate{
				MultiStakerAddress: "",
				ValidatorAddress:   "",
				Amount:             coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid shares amount",
			msg: types.MsgUndelegate{
				MultiStakerAddress: mulStakeAddr.String(),
				ValidatorAddress:   valAddr.String(),
				Amount:             sdk.NewCoin(denom, sdk.ZeroInt()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

// // MsgVote
// func TestMsgVote_ValidateBasic(t *testing.T) {
// 	mulStakeAddr := testutil.GenAddress()

// 	tests := []struct {
// 		name string
// 		msg  types.MsgVote
// 		err  error
// 	}{
// 		{
// 			name: "happy path",
// 			msg: types.MsgVote{
// 				ProposalId: 1,
// 				Voter:      mulStakeAddr.String(),
// 				Option:     govv1.VoteOption_VOTE_OPTION_NO_WITH_VETO,
// 				Metadata:   "",
// 			},
// 		},
// 		{
// 			name: "invalid voter address",
// 			msg: types.MsgVote{
// 				ProposalId: 1,
// 				Voter:      "",
// 				Option:     govv1.OptionYes,
// 				Metadata:   "",
// 			},
// 			err: sdkerrors.ErrInvalidAddress,
// 		},
// 		{
// 			name: "invalid vote option",
// 			msg: types.MsgVote{
// 				ProposalId: 1,
// 				Voter:      mulStakeAddr.String(),
// 				Metadata:   "",
// 			},
// 			err: govtypes.ErrInvalidVote,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := tt.msg.ValidateBasic()
// 			if tt.err != nil {
// 				require.ErrorIs(t, err, tt.err)
// 				return
// 			}
// 			require.NoError(t, err)
// 		})
// 	}
// }

// // MsgVoteWeighted
// func TestMsgVoteWeighted_ValidateBasic(t *testing.T) {
// 	mulStakeAddr := testutil.GenAddress()

// 	tests := []struct {
// 		name string
// 		msg  types.MsgVoteWeighted
// 		err  error
// 	}{
// 		{
// 			name: "happy path",
// 			msg: types.MsgVoteWeighted{
// 				ProposalId: 1,
// 				Voter:      mulStakeAddr.String(),
// 				Options:    govv1.NewNonSplitVoteOption(govv1.OptionYes),
// 				Metadata:   "",
// 			},
// 		},
// 		{
// 			name: "invalid voter address",
// 			msg: types.MsgVoteWeighted{
// 				ProposalId: 1,
// 				Voter:      "",
// 				Options:    govv1.NewNonSplitVoteOption(govv1.OptionYes),
// 				Metadata:   "",
// 			},
// 			err: sdkerrors.ErrInvalidAddress,
// 		},
// 		{
// 			name: "invalid vote options",
// 			msg: types.MsgVoteWeighted{
// 				ProposalId: 1,
// 				Voter:      mulStakeAddr.String(),
// 				Metadata:   "",
// 			},
// 			err: sdkerrors.ErrInvalidRequest,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := tt.msg.ValidateBasic()
// 			if tt.err != nil {
// 				require.ErrorIs(t, err, tt.err)
// 				return
// 			}
// 			require.NoError(t, err)
// 		})
// 	}
// }

// // MsgWithdrawDelegatorReward
// func TestMsgWithdrawDelegatorReward_ValidateBasic(t *testing.T) {
// 	mulStakeAddr := testutil.GenAddress()
// 	valAddr := testutil.GenValAddress()

// 	tests := []struct {
// 		name string
// 		msg  types.MsgWithdrawDelegatorReward
// 		err  error
// 	}{
// 		{
// 			name: "happy path",
// 			msg: types.MsgWithdrawDelegatorReward{
// 				MultiStakerAddress: mulStakeAddr.String(),
// 				ValidatorAddress:   valAddr.String(),
// 			},
// 		},
// 		{
// 			name: "invalid address",
// 			msg: types.MsgWithdrawDelegatorReward{
// 				MultiStakerAddress: "",
// 				ValidatorAddress:   "",
// 			},
// 			err: sdkerrors.ErrInvalidAddress,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := tt.msg.ValidateBasic()
// 			if tt.err != nil {
// 				require.ErrorIs(t, err, tt.err)
// 				return
// 			}
// 			require.NoError(t, err)
// 		})
// 	}
// }
