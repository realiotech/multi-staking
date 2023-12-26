package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeAddBondDenom    = "AddBondDenom"
	ProposalTypeUpdateBondDenom = "UpdateBondDenom"
	ProposalTypeRemoveBondDenom = "RemoveBondDenom"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeAddBondDenom)
	govtypes.RegisterProposalType(ProposalTypeUpdateBondDenom)
	govtypes.RegisterProposalType(ProposalTypeRemoveBondDenom)
}

var (
	_ govtypes.Content = &AddMultiStakingCoinProposal{}
	_ govtypes.Content = &UpdateBondWeightProposal{}
	_ govtypes.Content = &RemoveMultiStakingCoinProposal{}
)

func NewAddBondDenomProposal(title, description, denom string, bondWeight sdk.Dec) govtypes.Content {
	return &AddMultiStakingCoinProposal{
		Title:       title,
		Description: description,
		Denom:       denom,
		BondWeight:  &bondWeight,
	}
}

func (p *AddMultiStakingCoinProposal) ProposalRoute() string { return RouterKey }

func (p *AddMultiStakingCoinProposal) ProposalType() string {
	return ProposalTypeAddBondDenom
}

func (p *AddMultiStakingCoinProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if p.Denom == "" {
		return fmt.Errorf("denom %s does not exist", p.Denom)
	}

	if p.BondWeight.LT(sdk.ZeroDec()) {
		return fmt.Errorf("BondCoinWeight cannot be less than 0")
	}

	return nil
}

func NewUpdateBondDenomProposal(title, description, denom string, updatedBondWeight sdk.Dec) govtypes.Content {
	return &UpdateBondWeightProposal{
		Title:             title,
		Description:       description,
		Denom:             denom,
		UpdatedBondWeight: &updatedBondWeight,
	}
}

func (p *UpdateBondWeightProposal) ProposalRoute() string { return RouterKey }

func (p *UpdateBondWeightProposal) ProposalType() string {
	return ProposalTypeUpdateBondDenom
}

func (p *UpdateBondWeightProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if p.Denom == "" {
		return fmt.Errorf("denom %s does not exist", p.Denom)
	}

	if p.UpdatedBondWeight.LT(sdk.ZeroDec()) {
		return fmt.Errorf("BondCoinWeight cannot be less than 0")
	}

	return nil
}

func NewRemoveBondDenomProposal(title, description, denom string) govtypes.Content {
	return &RemoveMultiStakingCoinProposal{
		Title:       title,
		Description: description,
		Denom:       denom,
	}
}

func (p *RemoveMultiStakingCoinProposal) ProposalRoute() string { return RouterKey }

func (p *RemoveMultiStakingCoinProposal) ProposalType() string {
	return ProposalTypeRemoveBondDenom
}

func (p *RemoveMultiStakingCoinProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if p.Denom == "" {
		return fmt.Errorf("denom %s does not exist", p.Denom)
	}

	return nil
}
