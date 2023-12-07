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
	_ govtypes.Content = &UpdateBondTokenWeightProposals{}
)

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
		return fmt.Errorf("denom %s does not exist", p.BondDenomChange)
	}

	if p.BondTokenWeightChange.LT(sdk.ZeroDec()) {
		return fmt.Errorf("BondTokenWeight cannot be less than 0")
	}

	return nil
}
