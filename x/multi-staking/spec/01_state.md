<!--
order: 1
-->

# State

## Data Structure

* 






## Bond Token Weight

* BondTokenWeight: `0x00 | BondDenom -> BondTokenWeight (sdk.Dec)`

## Validator Bond Denom

* ValidatorBondDenom: `0x01 | ValOperatorAddr -> BondDenom (string)`

## Delegated Tokens

Uses for accounting of locked tokens

* Delegated Tokens: `0x02 | DelAddr | ValAddr -> Delegated Tokens (sdk.Coin)`

## Staking Records

Uses for claiming reward

* StakingRecords: `0x03 | DelAddr -> store(StakingRecords)`

##


---




















