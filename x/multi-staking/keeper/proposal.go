package keeper

import (
	"bytes"
	"fmt"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	erc20types "github.com/cosmos/evm/x/erc20/types"
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
		Signer: k.authority,
		Erc20Addresses: []string{p.ContractAddress},
	})
	if err != nil {
		return err
	}

	// Get the denom of the registered erc20 token
	tokenId = k.erc20keeper.GetTokenPairID(ctx, p.ContractAddress) 
	if bytes.Equal(tokenId, []byte{}) {
		return fmt.Errorf("tokenId %s not found", p.ContractAddress)
	}
	tokenPair, found := k.erc20keeper.GetTokenPair(ctx, tokenId)
	if !found {
		return fmt.Errorf("token pair %s not found", p.ContractAddress)
	}
	tokenDenom := tokenPair.Denom

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
