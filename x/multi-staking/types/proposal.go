package types

import (
	sdkerrors "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// Proposal types
const (
	ProposalTypeAddBondToken          string = "AddBondToken"
	ProposalTypeChangeBondTokenWeight string = "ChangeBondTokenWeight"
)

// Assert module proposals implement govtypes.Content at compile-time
var (
	_ govv1beta1.Content = &AddBondTokenProposal{}
	_ govv1beta1.Content = &ChangeBondTokenWeightProposal{}
)

func init() {
	govv1beta1.RegisterProposalType(ProposalTypeAddBondToken)
	govv1beta1.RegisterProposalType(ProposalTypeChangeBondTokenWeight)
}

// NewAddBondTokenProposal returns new instance of AddBondTokenProposal
func NewAddBondTokenProposal(title, description, bondToken string, tokenWeight sdk.Dec) govv1beta1.Content {
	return &AddBondTokenProposal{
		Title:       title,
		Description: description,
		BondToken:   bondToken,
		TokenWeight: &tokenWeight,
	}
}

// GetTitle returns the title of a AddBondTokenProposal
func (abtp *AddBondTokenProposal) GetTitle() string { return abtp.Title }

// GetDescription returns the description of a AddBondTokenProposal
func (abtp *AddBondTokenProposal) GetDescription() string { return abtp.Description }

// ProposalRoute returns router key for a AddBondTokenProposal
func (*AddBondTokenProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for a AddBondTokenProposal
func (*AddBondTokenProposal) ProposalType() string {
	return ProposalTypeAddBondToken
}

// ValidateBasic runs basic stateless validity checks
func (abtp *AddBondTokenProposal) ValidateBasic() error {
	err := govv1beta1.ValidateAbstract(abtp)
	if err != nil {
		return err
	}

	if abtp.BondToken == "" {
		return sdkerrors.Wrap(ErrInvalidAddBondTokenProposal, "proposal bond token cannot be blank")
	}

	if !abtp.TokenWeight.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidAddBondTokenProposal, "proposal bond token weight must be positive")
	}

	return nil
}

// String implements the Stringer interface.
func (abtp AddBondTokenProposal) String() string {
	return fmt.Sprintf("AddBondTokenProposal: Title: %s Description: %s BondToken: %s TokenWeight: %s", abtp.Title, abtp.Description, abtp.BondToken, abtp.TokenWeight)
}

// NewChangeBondTokenWeightProposal returns new instance of ChangeBondTokenWeightProposal
func NewChangeBondTokenWeightProposal(title, description, bondToken string, tokenWeight sdk.Dec) govv1beta1.Content {
	return &ChangeBondTokenWeightProposal{
		Title:       title,
		Description: description,
		BondToken:   bondToken,
		TokenWeight: &tokenWeight,
	}
}

// GetTitle returns the title of a ChangeBondTokenWeightProposal
func (cbtp *ChangeBondTokenWeightProposal) GetTitle() string { return cbtp.Title }

// GetDescription returns the description of a ChangeBondTokenWeightProposal
func (cbtp *ChangeBondTokenWeightProposal) GetDescription() string { return cbtp.Description }

// ProposalRoute returns router key for a ChangeBondTokenWeightProposal
func (*ChangeBondTokenWeightProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for a ChangeBondTokenWeightProposal
func (*ChangeBondTokenWeightProposal) ProposalType() string {
	return ProposalTypeChangeBondTokenWeight
}

// String implements the Stringer interface.
func (cbtp ChangeBondTokenWeightProposal) String() string {
	return fmt.Sprintf("ChangeBondTokenWeightProposal: Title: %s Description: %s BondToken: %s TokenWeight: %s", cbtp.Title, cbtp.Description, cbtp.BondToken, cbtp.TokenWeight)
}

// ValidateBasic runs basic stateless validity checks
func (cbtp *ChangeBondTokenWeightProposal) ValidateBasic() error {
	err := govv1beta1.ValidateAbstract(cbtp)
	if err != nil {
		return err
	}

	if cbtp.BondToken == "" {
		return sdkerrors.Wrap(ErrInvalidChangeBondTokenWeightProposal, "proposal bond token cannot be blank")
	}

	if !cbtp.TokenWeight.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidChangeBondTokenWeightProposal, "proposal bond token weight must be positive")
	}

	return nil
}
