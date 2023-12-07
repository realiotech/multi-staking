package types

import (
	"fmt"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeRemoveBondDenom = "RemoveBondDenom"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeRemoveBondDenom)
}

func NewRemoveBondDenomProposal(title, description, bondDenom string) govtypes.Content {
	return &RemoveBondTokenProposal{
		Title:           title,
		Description:     description,
		BondTokenRemove: bondDenom,
	}
}

var (
	_ govtypes.Content = &RemoveBondTokenProposal{}
)

func (p *RemoveBondTokenProposal) ProposalRoute() string { return RouterKey }

func (p *RemoveBondTokenProposal) ProposalType() string {
	return ProposalTypeRemoveBondDenom
}

func (p *RemoveBondTokenProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if p.BondTokenRemove == "" {
		return fmt.Errorf("denom %s does not exist", p.BondTokenRemove)
	}

	return nil
}
