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
	return &RemoveBondCoinProposal{
		Title:          title,
		Description:    description,
		BondCoinRemove: bondDenom,
	}
}

var (
	_ govtypes.Content = &RemoveBondCoinProposal{}
)

func (p *RemoveBondCoinProposal) ProposalRoute() string { return RouterKey }

func (p *RemoveBondCoinProposal) ProposalType() string {
	return ProposalTypeRemoveBondDenom
}

func (p *RemoveBondCoinProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if p.BondCoinRemove == "" {
		return fmt.Errorf("denom %s does not exist", p.BondCoinRemove)
	}

	return nil
}
