package types

const (
	// ModuleName defines the module name
	ModuleName = "multi-staking"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "memory:capability"
)

// KVStore keys
var (
	BondTokenWeightKey    = []byte{0x00}
	ValidatorBondDenomKey = []byte{0x01}

	IntermediaryAccountDelegator = []byte{0x02}

	DVPairSDKBondTokens = []byte{0x03}

	DVPairBondTokens = []byte{0x04}

	// mem store key
	CompletedDelegationsPrefix = []byte{0x00}
)
