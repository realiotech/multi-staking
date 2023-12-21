package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeAddBondDenom = "AddBondDenom"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeAddBondDenom)
}

var (
	_ govtypes.Content = &AddBondDenomProposal{}
)

func NewAddBondDenomProposal(title, description, bondDenom string, bondDenomWeight sdk.Dec) govtypes.Content {
	return &AddBondDenomProposal{
		Title:              title,
		Description:        description,
		BondTokenAdd:       bondDenom,
		BondTokenWeightAdd: bondDenomWeight,
	}
}

func (p *AddBondDenomProposal) ProposalRoute() string { return RouterKey }

func (p *AddBondDenomProposal) ProposalType() string {
	return ProposalTypeAddBondDenom
}

func (p *AddBondDenomProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if p.BondTokenAdd == "" {
		return fmt.Errorf("denom %s does not exist", p.BondTokenAdd)
	}

	if p.BondTokenWeightAdd.LT(sdk.ZeroDec()) {
		return fmt.Errorf("BondTokenWeight cannot be less than 0")
	}

	return nil
}
