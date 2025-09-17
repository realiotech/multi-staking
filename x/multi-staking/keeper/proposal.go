package keeper

import (
	"bytes"
	"fmt"

	erc20types "github.com/cosmos/evm/x/erc20/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AddMultiStakingCoinProposal handles the proposals to add a new bond token
func (k Keeper) AddMultiStakingCoinProposal(
	ctx sdk.Context,
	p *types.AddMultiStakingCoinProposal,
) error {
	_, found := k.GetBondWeight(ctx, p.Denom)
	if found {
		return fmt.Errorf("Error MultiStakingCoin %s already exist", p.Denom) //nolint:stylecheck
	}

	bondWeight := *p.BondWeight
	if bondWeight.LTE(math.LegacyZeroDec()) {
		return fmt.Errorf("Error MultiStakingCoin BondWeight %s invalid", bondWeight) //nolint:stylecheck
	}

	k.SetBondWeight(ctx, p.Denom, bondWeight)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddMultiStakingCoin,
			sdk.NewAttribute(types.AttributeKeyDenom, p.Denom),
			sdk.NewAttribute(types.AttributeKeyBondWeight, p.BondWeight.String()),
		),
	)
	return nil
}

// AddMultiStakingEVMCoinProposal handles the proposals to add a new bond token
func (k Keeper) AddMultiStakingEVMCoinProposal(
	ctx sdk.Context,
	p *types.AddMultiStakingEVMCoinProposal,
) error {
	// Check if the contract address is already registered in erc20 module
	tokenId := k.erc20keeper.GetTokenPairID(ctx, p.ContractAddress)
	if !bytes.Equal(tokenId, []byte{}) {
		return fmt.Errorf("Error ERC20 token %s already registered", p.ContractAddress) //nolint:stylecheck
	}

	// Register the erc20 token
	_, err := k.erc20keeper.RegisterERC20(ctx, &erc20types.MsgRegisterERC20{
		Signer:         k.authority,
		Erc20Addresses: []string{p.ContractAddress},
	})
	if err != nil {
		return err
	}

	// Get the denom of the registered erc20 token
	tokenDenom, err := k.erc20keeper.GetTokenDenom(ctx, common.HexToAddress(p.ContractAddress))
	if err != nil {
		return err
	}

	// Register the token as a multistaking coin
	return k.AddMultiStakingCoinProposal(ctx, &types.AddMultiStakingCoinProposal{
		Title:       p.Title,
		Description: p.Description,
		Denom:       tokenDenom,
		BondWeight:  p.BondWeight,
	})
}

func (k Keeper) BondWeightProposal(
	ctx sdk.Context,
	p *types.UpdateBondWeightProposal,
) error {
	_, found := k.GetBondWeight(ctx, p.Denom)
	if !found {
		return fmt.Errorf("Error MultiStakingCoin %s not found", p.Denom) //nolint:stylecheck
	}

	bondWeight := *p.UpdatedBondWeight
	if bondWeight.LTE(math.LegacyZeroDec()) {
		return fmt.Errorf("Error MultiStakingCoin BondWeight %s invalid", bondWeight) //nolint:stylecheck
	}

	k.SetBondWeight(ctx, p.Denom, bondWeight)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddMultiStakingCoin,
			sdk.NewAttribute(types.AttributeKeyDenom, p.Denom),
			sdk.NewAttribute(types.AttributeKeyBondWeight, p.UpdatedBondWeight.String()),
		),
	)
	return nil
}

// RemoveMultiStakingCoinProposal handles the proposals to remove a bond token
// We will force undelegate all the delegation of the removed bond token
// Remove bond token from store
func (k Keeper) RemoveMultiStakingCoinProposal(
	ctx sdk.Context,
	p *types.RemoveMultiStakingCoinProposal,
) error {
	_, found := k.GetBondWeight(ctx, p.Denom)
	if !found {
		return fmt.Errorf("Error MultiStakingCoin %s not found", p.Denom) //nolint:stylecheck
	}

	var ubdErr error
	k.MultiStakingLockIterator(ctx, func(stakingLock types.MultiStakingLock) bool {
		if stakingLock.LockedCoin.Denom != p.Denom {
			return false
		}
		// Call the Keeper method directly instead of going through MsgServer
		_, err := k.Undelegate(ctx, &stakingtypes.MsgUndelegate{
			DelegatorAddress: stakingLock.LockID.MultiStakerAddr,
			ValidatorAddress: stakingLock.LockID.ValAddr,
			Amount:           stakingLock.LockedCoin.ToCoin(),
		})
		if err != nil {
			ubdErr = err
			return true
		}

		return false
	})
	if ubdErr != nil {
		return ubdErr
	}

	k.RemoveBondWeight(ctx, p.Denom)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveMultiStakingCoin,
			sdk.NewAttribute(types.AttributeKeyDenom, p.Denom),
		),
	)
	return nil
}
