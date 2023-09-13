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

del a -> val b 100 ario -> 100 stake

a | b -> 100 stake

used to calculate the unlocked bond token when undbonding

sdk bond token -> bond token


500 stake

`sdk undelegate` 250 stake

50% off

`bond token` 50% off

a | b -> 500 stake

100 -> 5


unbonding amount bond token 100 ario

* DVPairSDKBondToken: `0x03 | DVPair -> SDKBondTokens`

### DelegatorTotalSDKBondToken

* DelegatorTotalSDKBondToken: `0x04 | Delegator -> TotalSDKBondToken`

## MemStore

### CompletedDelegations

* CompletedDelegations :`0x04 -> store(delegations)`