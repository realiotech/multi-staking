package simapp

import (
	evmenc "github.com/evmos/ethermint/encoding"

	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
)

func MakeEncodingConfig() simappparams.EncodingConfig {
	return evmenc.MakeConfig(ModuleBasics)
}
