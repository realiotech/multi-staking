package types

// x/multistaking  module event types
const (
	EventTypeAddMultiStakingCoin = "add_bond_token"
	EventTypeUpdateBondWeight    = "change_bond_token_weight"

	AttributeKeyBondToken  = "bond_token"
	AttributeKeyBondWeight = "bond_token_weight"
)
