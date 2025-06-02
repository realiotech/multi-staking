package simapp

import (
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	simappparams "cosmossdk.io/simapp/params"
	"cosmossdk.io/x/tx/signing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
)

// MakeEncodingConfig creates the EncodingConfig for realio network
func MakeEncodingConfig() simappparams.EncodingConfig {
	legacyAmino := codec.NewLegacyAmino()
	signingOptions := signing.Options{
		AddressCodec: address.Bech32Codec{
			Bech32Prefix: sdk.GetConfig().GetBech32AccountAddrPrefix(),
		},
		ValidatorAddressCodec: address.Bech32Codec{
			Bech32Prefix: sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
		},
		CustomGetSigners: map[protoreflect.FullName]signing.GetSignersFunc{},
	}

	interfaceRegistry, _ := types.NewInterfaceRegistryWithOptions(types.InterfaceRegistryOptions{
		ProtoFiles:     proto.HybridResolver,
		SigningOptions: signingOptions,
	})
	codec := codec.NewProtoCodec(interfaceRegistry)

	txConfig := tx.NewTxConfig(codec, tx.DefaultSignModes)

	std.RegisterInterfaces(interfaceRegistry)

	ModuleBasics.RegisterLegacyAminoCodec(legacyAmino)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)

	return simappparams.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          txConfig,
		Amino:             legacyAmino,
	}
}
