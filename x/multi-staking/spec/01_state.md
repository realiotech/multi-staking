<!--
order: 1
-->

# State

## Store

### Bond Token Weight

* BondTokenWeight: `0x00 | BondDenom -> BondTokenWeight (sdk.Dec)`

### Validator Bond Denom

* ValidatorBondDenom: `0x01 | ValOperatorAddr -> BondDenom (string)`

### Intermediary Account Delegator

* IntermediaryAccountDelegator: `0x02 | IntermediaryAccount -> DelegatorAddr`

### DV Pair SDK Bond Tokens

* DVPairSDKBondToken: `0x03 | DVPair -> SDKBondTokens`

### DelegatorTotalSDKBondToken

* DelegatorTotalSDKBondToken; `0x04 | Delegator -> TotalSDKBondToken`

## MemStore

### CompletedDelegations

* CompletedDelegations :`0x04 -> store(delegations)`