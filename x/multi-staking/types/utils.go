package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// BaseDenom defines the base denomination
	BaseDenom = "ario"

	// GovDenom defines the governance denomination
	GovDenom = "arst"
)

func AccAddrAndValAddrFromStrings(accAddrString string, valAddrStraing string) (sdk.AccAddress, sdk.ValAddress, error) {
	accAddr, err := sdk.AccAddressFromBech32(accAddrString)
	if err != nil {
		return sdk.AccAddress{}, sdk.ValAddress{}, err
	}
	valAcc, err := sdk.ValAddressFromBech32(valAddrStraing)
	if err != nil {
		return sdk.AccAddress{}, sdk.ValAddress{}, err
	}

	return accAddr, valAcc, nil
}
