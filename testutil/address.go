package testutil

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

)

func GenAddress() sdk.AccAddress {
	priv := secp256k1.GenPrivKey()

	return sdk.AccAddress(priv.PubKey().Address())
}

func GenValAddress() sdk.ValAddress {
	priv := secp256k1.GenPrivKey()

	return sdk.ValAddress(priv.PubKey().Address())
}

func GenPubKey() cryptotypes.PubKey {
	priv := secp256k1.GenPrivKey()
	return priv.PubKey()
}