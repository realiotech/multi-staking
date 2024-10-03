package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	v1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

var (
	amino = codec.NewLegacyAmino()

	// AminoCdc references the global x/relationships module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/relationships and
	// defined at the application level.
	AminoCdc = codec.NewAminoCodec(amino)
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateMultiStakingParams{}, "multistaking/MsgUpdateMSParams")

	cdc.RegisterConcrete(&AddMultiStakingCoinProposal{}, "multistaking/AddMultiStakingCoinProposal", nil)
	cdc.RegisterConcrete(&UpdateBondWeightProposal{}, "multistaking/UpdateBondWeightProposal", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateMultiStakingParams{},
	)
	registry.RegisterImplementations(
		(*v1beta1types.Content)(nil),
		&AddMultiStakingCoinProposal{},
		&UpdateBondWeightProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)
}
