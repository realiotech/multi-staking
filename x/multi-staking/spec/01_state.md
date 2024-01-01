<!--
order: 1
-->

# State

## Store

### Bond Token Weight

* BondWeight: `0x00 | BondDenom -> BondWeight (sdk.Dec)`

### Validator Bond Denom

* ValidatorMultiStakingCoin: `0x01 | ValOperatorAddr -> BondDenom (string)`

### Intermediary Account Delegator

* IntermediaryDelegator: `0x02 | IntermediaryAccount -> DelegatorAddr`

### DV Pair SDK Bond Tokens

* DVPairSDKBondToken: `0x03 | DVPair -> SDKBondTokens`

### DV Pair Bond Token

* DVPairBondToken: `0x04 | DVPair -> BondTokens`

## MemStore

### CompletedDelegations

* CompletedDelegations :`0x00 -> store(delegations)`