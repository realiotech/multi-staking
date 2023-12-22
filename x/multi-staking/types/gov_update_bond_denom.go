package types

import (
	"fmt"

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
	_ govtypes.Content = &UpdateBondCoinWeightProposals{}
)

func NewUpdateBondDenomProposal(title, description, bondDenom string, bondDenomWeight sdk.Dec) govtypes.Content {
	return &UpdateBondCoinWeightProposals{
		Title:                title,
		Description:          description,
		BondDenomChange:      bondDenom,
		BondCoinWeightChange: &bondDenomWeight,
	}
}

func (p *UpdateBondCoinWeightProposals) ProposalRoute() string { return RouterKey }

func (p *UpdateBondCoinWeightProposals) ProposalType() string {
	return ProposalTypeUpdateBondDenom
}

func (p *UpdateBondCoinWeightProposals) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	if p.BondDenomChange == "" {
		return fmt.Errorf("denom %s does not exist", p.BondDenomChange)
	}

	if p.BondCoinWeightChange.LT(sdk.ZeroDec()) {
		return fmt.Errorf("BondCoinWeight cannot be less than 0")
	}

	return nil
}
