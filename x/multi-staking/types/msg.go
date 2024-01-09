package types

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// staking message types
const (
	TypeMsgUndelegate         = "begin_unbonding"
	TypeMsgCancelUnbonding    = "cancel_unbond"
	TypeMsgEditValidator      = "edit_validator"
	TypeMsgCreateValidator    = "create_validator"
	TypeMsgDelegate           = "delegate"
	TypeMsgBeginRedelegate    = "begin_redelegate"
	TypeMsgVote               = "vote"
	TypeMsgVoteWeighted       = "weighted_vote"
	TypeMsgSetWithdrawAddress = "set_withdraw_address"
	TypeMsgWithdrawReward     = "withdraw_delegator_reward"
)

var (
	_ sdk.Msg                            = &MsgCreateValidator{}
	_ codectypes.UnpackInterfacesMessage = (*MsgCreateValidator)(nil)
	_ sdk.Msg                            = &MsgCreateValidator{}
	_ sdk.Msg                            = &MsgEditValidator{}
	_ sdk.Msg                            = &MsgDelegate{}
	_ sdk.Msg                            = &MsgUndelegate{}
	_ sdk.Msg                            = &MsgBeginRedelegate{}
	_ sdk.Msg                            = &MsgCancelUnbonding{}
)

// NewMsgCreateValidator creates a new MsgCreateValidator instance.
// Delegator address and validator address are the same.
func NewMsgCreateValidator(
	valAddr sdk.ValAddress, pubKey cryptotypes.PubKey, //nolint:interfacer
	selfDelegation sdk.Coin, description stakingtypes.Description, commission stakingtypes.CommissionRates, minSelfDelegation math.Int,
) (*MsgCreateValidator, error) {
	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	}
	return &MsgCreateValidator{
		Description:        description,
		MultiStakerAddress: sdk.AccAddress(valAddr).String(),
		ValidatorAddress:   valAddr.String(),
		Pubkey:             pkAny,
		Value:              selfDelegation,
		Commission:         commission,
		MinSelfDelegation:  minSelfDelegation,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgCreateValidator) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCreateValidator) Type() string { return TypeMsgCreateValidator }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.
func (msg MsgCreateValidator) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	multiStakerAddr, valAddr, _ := AccAddrAndValAddrFromStrings(msg.MultiStakerAddress, msg.ValidatorAddress)
	addrs := []sdk.AccAddress{multiStakerAddr}

	valAccAddr := sdk.AccAddress(valAddr)
	if !multiStakerAddr.Equals(valAccAddr) {
		addrs = append(addrs, valAccAddr)
	}

	return addrs
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCreateValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreateValidator) ValidateBasic() error {
	// note that unmarshaling from bech32 ensures both non-empty and valid
	multiStakerAddr, valAddr, err := AccAddrAndValAddrFromStrings(msg.MultiStakerAddress, msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap(err.Error())
	}

	if !sdk.AccAddress(valAddr).Equals(multiStakerAddr) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "validator address is invalid")
	}

	if msg.Pubkey == nil {
		return stakingtypes.ErrEmptyValidatorPubKey
	}

	if !msg.Value.IsValid() || !msg.Value.Amount.IsPositive() {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid delegation amount")
	}

	if msg.Description == (stakingtypes.Description{}) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.Commission == (stakingtypes.CommissionRates{}) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "empty commission")
	}

	if err := msg.Commission.Validate(); err != nil {
		return err
	}

	if !msg.MinSelfDelegation.IsPositive() {
		return errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"minimum self delegation must be a positive integer",
		)
	}

	if msg.Value.Amount.LT(msg.MinSelfDelegation) {
		return stakingtypes.ErrSelfDelegationBelowMinimum
	}

	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateValidator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.Pubkey, &pubKey)
}

// NewMsgEditValidator creates a new MsgEditValidator instance
//
//nolint:interfacer
func NewMsgEditValidator(valAddr sdk.ValAddress, description stakingtypes.Description, newRate *sdk.Dec, newMinSelfDelegation *math.Int) *MsgEditValidator {
	return &MsgEditValidator{
		Description:       description,
		CommissionRate:    newRate,
		ValidatorAddress:  valAddr.String(),
		MinSelfDelegation: newMinSelfDelegation,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgEditValidator) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgEditValidator) Type() string { return TypeMsgEditValidator }

// GetSigners implements the sdk.Msg interface.
func (msg MsgEditValidator) GetSigners() []sdk.AccAddress {
	valAddr, _ := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	return []sdk.AccAddress{sdk.AccAddress(valAddr)}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgEditValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgEditValidator) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if msg.Description == (stakingtypes.Description{}) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.MinSelfDelegation != nil && !msg.MinSelfDelegation.IsPositive() {
		return errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"minimum self delegation must be a positive integer",
		)
	}

	if msg.CommissionRate != nil {
		if msg.CommissionRate.GT(sdk.OneDec()) || msg.CommissionRate.IsNegative() {
			return errors.Wrap(sdkerrors.ErrInvalidRequest, "commission rate must be between 0 and 1 (inclusive)")
		}
	}

	return nil
}

// NewMsgDelegate creates a new MsgDelegate instance.
//
//nolint:interfacer
func NewMsgDelegate(multiStakerAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) *MsgDelegate {
	return &MsgDelegate{
		MultiStakerAddress: multiStakerAddr.String(),
		ValidatorAddress:   valAddr.String(),
		Amount:             amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgDelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgDelegate) Type() string { return TypeMsgDelegate }

// GetSigners implements the sdk.Msg interface.
func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
	delegator, _ := sdk.AccAddressFromBech32(msg.MultiStakerAddress)
	return []sdk.AccAddress{delegator}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgDelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgDelegate) ValidateBasic() error {
	if _, _, err := AccAddrAndValAddrFromStrings(msg.MultiStakerAddress, msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap(err.Error())
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid delegation amount",
		)
	}

	return nil
}

// NewMsgBeginRedelegate creates a new MsgBeginRedelegate instance.
//
//nolint:interfacer
func NewMsgBeginRedelegate(
	multiStakerAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, amount sdk.Coin,
) *MsgBeginRedelegate {
	return &MsgBeginRedelegate{
		MultiStakerAddress:  multiStakerAddr.String(),
		ValidatorSrcAddress: valSrcAddr.String(),
		ValidatorDstAddress: valDstAddr.String(),
		Amount:              amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgBeginRedelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg MsgBeginRedelegate) Type() string { return TypeMsgBeginRedelegate }

// GetSigners implements the sdk.Msg interface
func (msg MsgBeginRedelegate) GetSigners() []sdk.AccAddress {
	delegator, _ := sdk.AccAddressFromBech32(msg.MultiStakerAddress)
	return []sdk.AccAddress{delegator}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgBeginRedelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgBeginRedelegate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.MultiStakerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorSrcAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid source validator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorDstAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid destination validator address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid shares amount",
		)
	}

	return nil
}

// NewMsgUndelegate creates a new MsgUndelegate instance.
//
//nolint:interfacer
func NewMsgUndelegate(multiStakerAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) *MsgUndelegate {
	return &MsgUndelegate{
		MultiStakerAddress: multiStakerAddr.String(),
		ValidatorAddress:   valAddr.String(),
		Amount:             amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUndelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUndelegate) Type() string { return TypeMsgUndelegate }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUndelegate) GetSigners() []sdk.AccAddress {
	delegator, _ := sdk.AccAddressFromBech32(msg.MultiStakerAddress)
	return []sdk.AccAddress{delegator}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUndelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUndelegate) ValidateBasic() error {
	if _, _, err := AccAddrAndValAddrFromStrings(msg.MultiStakerAddress, msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap(err.Error())
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid shares amount",
		)
	}

	return nil
}

// NewMsgCancelUnbonding creates a new MsgCancelUnbonding instance.
//
//nolint:interfacer
func NewMsgCancelUnbonding(multiStakerAddr sdk.AccAddress, valAddr sdk.ValAddress, creationHeight int64, amount sdk.Coin) *MsgCancelUnbonding {
	return &MsgCancelUnbonding{
		MultiStakerAddress: multiStakerAddr.String(),
		ValidatorAddress:   valAddr.String(),
		Amount:             amount,
		CreationHeight:     creationHeight,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgCancelUnbonding) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCancelUnbonding) Type() string { return TypeMsgCancelUnbonding }

// GetSigners implements the sdk.Msg interface.
func (msg MsgCancelUnbonding) GetSigners() []sdk.AccAddress {
	delegator, _ := sdk.AccAddressFromBech32(msg.MultiStakerAddress)
	return []sdk.AccAddress{delegator}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgCancelUnbonding) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCancelUnbonding) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.MultiStakerAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid amount",
		)
	}

	if msg.CreationHeight <= 0 {
		return errors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid height",
		)
	}

	return nil
}
