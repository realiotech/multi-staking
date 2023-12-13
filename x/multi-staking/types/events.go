package types

// x/multistaking  module event types
const (
	EventTypeAddBondToken          = "add_bond_token"
	EventTypeUpdateBondTokenWeight = "update_bond_token_weight"
	EventTypeRemoveBondTokenWeight = "remove_bond_token_weight"

	AttributeKeyBondToken       = "bond_token"
	AttributeKeyBondTokenWeight = "bond_token_weight"
)
