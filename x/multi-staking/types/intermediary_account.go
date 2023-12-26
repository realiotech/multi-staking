package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: make unit test for this
// this is against cosmos convention, doing this for more performance and less storage
func IntermediaryDelegator(delAddr sdk.AccAddress) sdk.AccAddress {
	return append(delAddr, 0x0)
}

func DelegatorAccount(intermediaryAcc sdk.AccAddress) sdk.AccAddress {
	return bytes.Clone(intermediaryAcc[:len(intermediaryAcc)-1])

}
