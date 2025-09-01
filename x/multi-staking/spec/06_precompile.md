<!--
order: 6
-->

# Precompiles

The multi-staking module provides Ethereum precompiles that enable MetaMask and other EVM wallets to interact with the multi-staking functionality directly through smart contract calls. These precompiles bridge the gap between EVM-based applications and the Cosmos SDK multi-staking module.

## Overview

The multi-staking precompiles allow users to:
- Create validators using ERC20 tokens
- Delegate ERC20 tokens to validators
- Undelegate ERC20 tokens from validators
- Redelegate ERC20 tokens between validators
- Cancel unbonding delegations
- Query delegation, unbonding delegation, and validator information

## Precompile Address

The multi-staking precompile is deployed at the following address:

```
0x0000000000000000000000000000000000000900
```

This address can be used to interact with the multi-staking functionality from MetaMask, web3 applications, or any EVM-compatible wallet.

## Precompile Interface

The precompiles implement the `IMultiStaking` interface defined in Solidity:

```solidity
interface IMultiStaking {
    // Transaction methods
    function delegate(string calldata erc20Token, string calldata validatorAddress, string calldata amount) external returns (bool success);
    function undelegate(string calldata erc20Token, string calldata validatorAddress, string calldata amount) external returns (int64 completionTime);
    function redelegate(string calldata erc20Token, string calldata srcValidatorAddress, string calldata dstValidatorAddress, string calldata amount) external returns (int64 completionTime);
    function cancelUnbondingDelegation(string calldata erc20Token, string calldata validatorAddress, string calldata amount, string calldata creationHeight) external returns (bool success);
    function createValidator(...) external returns (bool success);
    
    // Query methods
    function delegation(address delegatorAddress, string memory validatorAddress) external view returns (Coin calldata balance);
    function unbondingDelegation(address delegatorAddress, string memory validatorAddress) external view returns (UnbondingDelegationOutput calldata unbondingDelegation);
    function validator(address validatorAddress) external view returns (Validator calldata validator);
}
```

## Transaction Methods

### delegate

Delegates ERC20 tokens to a validator through the multi-staking module.

**Parameters:**
- `erc20Token` (string): The ERC20 contract address
- `validatorAddress` (string): The validator's operator address
- `amount` (string): The amount of ERC20 tokens to delegate

**Returns:**
- `success` (bool): True if delegation was successful

**Logic Flow:**
1. Get delegator address from sender
2. Creates a `MsgDelegateEVM` message
3. Executes the delegation through the multi-staking message server
4. Returns success status

### undelegate

Undelegates ERC20 tokens from a validator.

**Parameters:**
- `erc20Token` (string): The ERC20 contract address
- `validatorAddress` (string): The validator's operator address
- `amount` (string): The amount of ERC20 tokens to undelegate

**Returns:**
- `completionTime` (int64): Unix timestamp when unbonding will complete

**Logic Flow:**
1. Get delegator address from sender
2. Creates a `MsgUndelegateEVM` message
3. Executes the undelegation through the multi-staking message server
4. Returns the completion time of the unbonding period

### redelegate

Redelegates ERC20 tokens from one validator to another.

**Parameters:**
- `erc20Token` (string): The ERC20 contract address
- `srcValidatorAddress` (string): The source validator's operator address
- `dstValidatorAddress` (string): The destination validator's operator address
- `amount` (string): The amount of ERC20 tokens to redelegate

**Returns:**
- `completionTime` (int64): Unix timestamp when redelegation will complete

**Logic Flow:**
1. Get delegator address from sender
2. Creates a `MsgBeginRedelegateEVM` message
3. Executes the redelegation through the multi-staking message server
4. Returns the completion time

### cancelUnbondingDelegation

Cancels an unbonding delegation and re-delegates the tokens back to the validator.

**Parameters:**
- `erc20Token` (string): The ERC20 contract address
- `validatorAddress` (string): The validator's operator address
- `amount` (string): The amount of ERC20 tokens to cancel unbonding for
- `creationHeight` (string): The height at which the unbonding delegation was created

**Returns:**
- `success` (bool): True if cancellation was successful

**Logic Flow:**
1. Get delegator address from sender
2. Creates a `MsgCancelUnbondingEVMDelegation` message
3. Executes the cancellation through the multi-staking message server
4. Returns success status

### createValidator

Creates a new validator using ERC20 tokens for self-delegation.

**Parameters:**
- `pubkey` (string): The validator's consensus public key (base64 encoded)
- `contractAddress` (string): The ERC20 contract address for self-delegation
- `amount` (string): The amount of ERC20 tokens for self-delegation
- `moniker` (string): The validator's display name
- `identity` (string): The validator's identity signature
- `website` (string): The validator's website URL
- `security` (string): The validator's security contact
- `details` (string): Additional validator details
- `commissionRate` (string): The validator's commission rate
- `commissionMaxRate` (string): The maximum commission rate
- `commissionMaxChangeRate` (string): The maximum commission change rate per day
- `minSelfDelegation` (string): The minimum self-delegation amount

**Returns:**
- `success` (bool): True if validator creation was successful

**Logic Flow:**
1. Parses and validates all input parameters
2. Converts the public key from base64 to the appropriate format
3. Creates a `MsgCreateEVMValidator` message
4. Executes validator creation through the multi-staking message server
5. Returns success status

## Query Methods

### delegation

Queries the delegation amount for a specific delegator-validator pair.

**Parameters:**
- `delegatorAddress` (address): The delegator's Ethereum address
- `validatorAddress` (string): The validator's operator address

**Returns:**
- `balance` (Coin): The delegated amount with denomination and amount

### unbondingDelegation

Queries the unbonding delegation entries for a specific delegator-validator pair.

**Parameters:**
- `delegatorAddress` (address): The delegator's Ethereum address
- `validatorAddress` (string): The validator's operator address

**Returns:**
- `unbondingDelegation` (UnbondingDelegationOutput): Contains delegator address, validator address, and unbonding entries

### validator

Queries validator information by validator hex address.

**Parameters:**
- `validatorAddress` (address): The validator's Ethereum address

**Returns:**
- `validator` (Validator): Complete validator information including status, tokens, commission, etc.

## Data Types

### Validator
```solidity
struct Validator {
    string operatorAddress;
    string consensusPubkey;
    bool jailed;
    BondStatus status;
    uint256 tokens;
    uint256 delegatorShares;
    string description;
    int64 unbondingHeight;
    int64 unbondingTime;
    uint256 commission;
    uint256 minSelfDelegation;
    string bondDenom;
}
```

### UnbondingDelegationEntry
```solidity
struct UnbondingDelegationEntry {
    int64 creationHeight;
    uint256 balance;
}
```

### UnbondingDelegationOutput
```solidity
struct UnbondingDelegationOutput {
    string delegatorAddress;
    string validatorAddress;
    UnbondingDelegationEntry[] entries;
}
```

## Integration with MetaMask

The precompiles enable seamless integration with MetaMask and other EVM wallets:

1. **Contract Interaction**: Users can interact with the precompiles as if they were regular smart contracts
2. **Transaction Signing**: All transactions are signed using standard Ethereum transaction signing
3. **Gas Estimation**: Gas costs are calculated and displayed in MetaMask
4. **Event Monitoring**: Transaction results can be monitored through standard EVM event logs
