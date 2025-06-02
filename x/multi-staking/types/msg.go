package types

// import (
// 	sdkerrors "cosmossdk.io/errors"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// // staking message types
// const (
// 	TypeMsgUpdateMultiStakingParams = "update_multistaking_params"
// )

// var (
// 	_ sdk.Msg = &MsgUpdateMultiStakingParams{}
// )

// // GetSignBytes returns the raw bytes for a MsgUpdateParams message that
// // the expected signer needs to sign.
// func (m *MsgUpdateMultiStakingParams) GetSignBytes() []byte {
// 	bz := AminoCdc.MustMarshalJSON(m)
// 	return sdk.MustSortJSON(bz)
// }

// // ValidateBasic executes sanity validation on the provided data
// func (m *MsgUpdateMultiStakingParams) ValidateBasic() error {
// 	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
// 		return sdkerrors.Wrap(err, "invalid authority address")
// 	}
// 	return nil
// }

// // GetSigners returns the expected signers for a MsgUpdateParams message
// func (m *MsgUpdateMultiStakingParams) GetSigners() []sdk.AccAddress {
// 	addr, _ := sdk.AccAddressFromBech32(m.Authority)
// 	return []sdk.AccAddress{addr}
// }
