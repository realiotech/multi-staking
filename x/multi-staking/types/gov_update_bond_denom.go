package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeUpdateBondDenom = "UpdateBondDenom"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeUpdateBondDenom)
}

var (
	_ govtypes.Content = &UpdateBondTokenWeightProposals{}
)

func NewUpdateBondDenomProposal(title, description, bondDenom string, bondDenomWeight sdk.Dec) govtypes.Content {
	return &UpdateBondTokenWeightProposals{
		Title:                 title,
		Description:           description,
		BondDenomChange:       bondDenom,
		BondTokenWeightChange: &bondDenomWeight,
	}
}

func (p *UpdateBondTokenWeightProposals) ProposalRoute() string { return RouterKey }

func (p *UpdateBondTokenWeightProposals) ProposalType() string {
	return ProposalTypeUpdateBondDenom
}

func (p *UpdateBondTokenWeightProposals) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if p.BondDenomChange == "" {
		return ErrBrondDenomDoesNotExist
	}

	if p.BondTokenWeightChange.LT(sdk.ZeroDec()) {
		return ErrLessThanZero
	}

	return nil
}
