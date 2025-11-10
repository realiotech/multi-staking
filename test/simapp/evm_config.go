package simapp

import (
	evmtypes "github.com/cosmos/evm/x/vm/types"
	// realionetworktypes "github.com/realiotech/realio-network/types"
)

// EVMOptionsFn defines a function type for setting app options specifically for
// the Cosmos EVM app. The function should receive the chainID and return an error if
// any.
type EVMOptionsFn func(string) error

// NoOpEVMOptions is a no-op function that can be used when the app does not
// need any specific configuration.
func NoOpEVMOptions(_ string) error {
	return nil
}

var sealed = false

// ChainsCoinInfo is a map of the chain id and its corresponding EvmCoinInfo
// that allows initializing the app with different coin info based on the
// chain id
var ChainsCoinInfo = evmtypes.EvmCoinInfo{
	Denom:         "ustake",
	ExtendedDenom: "ustake",
	DisplayDenom:  "ustake",
	Decimals:      6,
}
