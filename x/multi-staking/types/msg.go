package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
