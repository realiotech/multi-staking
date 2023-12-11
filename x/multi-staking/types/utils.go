package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DelAccAndValAccFromStrings(delAddrString string, valAddrStraing string) (sdk.AccAddress, sdk.ValAddress, error) {
	delAcc, err := sdk.AccAddressFromBech32(delAddrString)
	if err != nil {
		return sdk.AccAddress{}, sdk.ValAddress{}, err
	}
	valAcc, err := sdk.ValAddressFromBech32(valAddrStraing)
	if err != nil {
		return sdk.AccAddress{}, sdk.ValAddress{}, err
	}

	return delAcc, valAcc, nil
}
