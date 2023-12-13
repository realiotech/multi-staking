package legacy

import (
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	types1 "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BondStatus int32

type Params struct {
	UnbondingTime     time.Duration `json:"unbonding_time"`
	MaxValidators     uint32        `json:"max_validators"`
	MaxEntries        uint32        `json:"max_entries"`
	HistoricalEntries uint32        `json:"historical_entries"`
	BondDenom         string        `json:"bond_denom"`
	MinCommissionRate sdk.Dec       `json:"min_commission_rate"`
}

type LastValidatorPower struct {
	Address string `json:"address"`
	Power   int64  `json:"power"`
}

type Description struct {
	Moniker         string `json:"moniker"`
	Identity        string `json:"identity"`
	Website         string `json:"website"`
	SecurityContact string `json:"security_contact"`
	Details         string `json:"details"`
}

type CommissionRates struct {
	Rate          sdk.Dec `json:"rate"`
	MaxRate       sdk.Dec `json:"max_rate"`
	MaxChangeRate sdk.Dec `json:"max_change_rate"`
}

type Commission struct {
	CommissionRates `json:"commission_rates"`
	UpdateTime      time.Time `json:"update_time"`
}

type Validator struct {
	OperatorAddress   string      `json:"operator_address"`
	ConsensusPubkey   *types1.Any `json:"consensus_pubkey"`
	Jailed            bool        `json:"jailed"`
	Status            BondStatus  `json:"status"`
	Tokens            sdkmath.Int `json:"tokens"`
	DelegatorShares   sdk.Dec     `json:"delegator_shares"`
	Description       Description `json:"description"`
	UnbondingHeight   int64       `json:"unbonding_height"`
	UnbondingTime     time.Time   `json:"unbonding_time"`
	Commission        Commission  `json:"commission"`
	MinSelfDelegation sdkmath.Int `json:"min_self_delegation"`
	BondDenom         string      `json:"bond_denom"`
}

type Delegation struct {
	DelegatorAddress string  `json:"delegator_address"`
	ValidatorAddress string  `json:"validator_address"`
	Shares           sdk.Dec `json:"shares"`
}

type UnbondingDelegationEntry struct {
	CreationHeight int64     `json:"creation_height"`
	CompletionTime time.Time `json:"completion_time"`
	InitialBalance sdk.Coin  `json:"initial_balance"`
	Balance        sdk.Coin  `json:"balance"`
}

type UnbondingDelegation struct {
	DelegatorAddress string                     `json:"delegator_address"`
	ValidatorAddress string                     `json:"validator_address"`
	Entries          []UnbondingDelegationEntry `json:"entries"`
}

type RedelegationEntry struct {
	CreationHeight int64     `json:"creation_height"`
	CompletionTime time.Time `json:"completion_time"`
	InitialBalance sdk.Coin  `json:"initial_balance"`
	SharesDst      sdk.Dec   `json:"shares_dst"`
}

type Redelegation struct {
	DelegatorAddress    string              `json:"delegator_address"`
	ValidatorSrcAddress string              `json:"validator_src_address"`
	ValidatorDstAddress string              `json:"validator_dst_address"`
	Entries             []RedelegationEntry `json:"entries"`
}

type GenesisState struct {
	Params               Params                `json:"params"`
	LastTotalPower       sdkmath.Int           `json:"last_total_power"`
	LastValidatorPowers  []LastValidatorPower  `json:"last_validator_powers"`
	Validators           []Validator           `json:"validators"`
	Delegations          []Delegation          `json:"delegations"`
	UnbondingDelegations []UnbondingDelegation `json:"unbonding_delegations"`
	Redelegations        []Redelegation        `json:"redelegations"`
	Exported             bool                  `json:"exported"`
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateValidator{}, "cosmos-sdk/MsgCreateValidator")
	legacy.RegisterAminoMsg(cdc, &MsgEditValidator{}, "cosmos-sdk/MsgEditValidator")
	legacy.RegisterAminoMsg(cdc, &MsgDelegate{}, "cosmos-sdk/MsgDelegate")
	legacy.RegisterAminoMsg(cdc, &MsgUndelegate{}, "cosmos-sdk/MsgUndelegate")
	legacy.RegisterAminoMsg(cdc, &MsgBeginRedelegate{}, "cosmos-sdk/MsgBeginRedelegate")
	legacy.RegisterAminoMsg(cdc, &MsgCancelUnbondingDelegation{}, "cosmos-sdk/MsgCancelUnbondingDelegation")

	cdc.RegisterInterface((*isStakeAuthorization_Validators)(nil), nil)
	cdc.RegisterConcrete(&StakeAuthorization_AllowList{}, "cosmos-sdk/StakeAuthorization/AllowList", nil)
	cdc.RegisterConcrete(&StakeAuthorization_DenyList{}, "cosmos-sdk/StakeAuthorization/DenyList", nil)
	cdc.RegisterConcrete(&StakeAuthorization{}, "cosmos-sdk/StakeAuthorization", nil)
}
