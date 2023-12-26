<!--
order: 1
-->

# State

## Store

### Bond Coin Weight

* BondCoinWeight: `0x00 | BondDenom -> BondCoinWeight (sdk.Dec)`

### Validator Bond Denom

* ValidatorAllowedCoin: `0x01 | ValOperatorAddr -> BondDenom (string)`

### Intermediary Account Delegator

* IntermediaryDelegatorDelegator: `0x02 | IntermediaryDelegator -> DelegatorAddr`

### DV Pair SDK Bond Coins

* DVPairSDKBondCoin: `0x03 | DVPair -> SDKBondCoins`

### DV Pair Bond Coin

* DVPairBondCoin: `0x04 | DVPair -> BondCoins`

## MemStore

### CompletedDelegations

* CompletedDelegations :`0x00 -> store(delegations)`