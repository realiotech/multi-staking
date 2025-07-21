package types

import (
	fmt "fmt"

	"cosmossdk.io/core/address"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
)

// staking message types
const (
	TypeMsgUpdateMultiStakingParams = "update_multistaking_params"
)

var (
	_ sdk.Msg                            = &MsgUpdateMultiStakingParams{}
	_ sdk.Msg                            = &MsgDelegateEVM{}
	_ sdk.Msg                            = &MsgBeginRedelegateEVM{}
	_ sdk.Msg                            = &MsgUndelegateEVM{}
	_ sdk.Msg                            = &MsgCancelUnbondingEVMDelegation{}
	_ codectypes.UnpackInterfacesMessage = (*MsgCreateEVMValidator)(nil)
)

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateEVMValidator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.Pubkey, &pubKey)
}

// NewMsgCreateEVMValidator creates a new MsgCreateEVMValidator instance.
// Delegator address and validator address are the same.
func NewMsgCreateEVMValidator(
	valAddr string, pubKey cryptotypes.PubKey, contractAddress string,
	selfDelegation math.Int, description stakingtypes.Description, commission stakingtypes.CommissionRates, minSelfDelegation math.Int,
) (*MsgCreateEVMValidator, error) {
	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	}
	return &MsgCreateEVMValidator{
		ContractAddress:   contractAddress,
		Description:       description,
		ValidatorAddress:  valAddr,
		Pubkey:            pkAny,
		Value:             selfDelegation,
		Commission:        commission,
		MinSelfDelegation: minSelfDelegation,
	}, nil
}

// Validate validates the MsgCreateEVMValidator sdk msg.
func (msg MsgCreateEVMValidator) Validate(ac address.Codec) error {
	// note that unmarshaling from bech32 ensures both non-empty and valid
	_, err := ac.StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if msg.Pubkey == nil {
		return stakingtypes.ErrEmptyValidatorPubKey
	}

	if !common.IsHexAddress(msg.ContractAddress) {
		return fmt.Errorf("invalid contract address")
	}
	
	if !!msg.Value.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid delegation amount")
	}

	if msg.Description == (stakingtypes.Description{}) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.Commission == (stakingtypes.CommissionRates{}) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "empty commission")
	}

	if err := msg.Commission.Validate(); err != nil {
		return err
	}

	if !msg.MinSelfDelegation.IsPositive() {
		return errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"minimum self delegation must be a positive integer",
		)
	}

	if msg.Value.LT(msg.MinSelfDelegation) {
		return stakingtypes.ErrSelfDelegationBelowMinimum
	}

	return nil
}
