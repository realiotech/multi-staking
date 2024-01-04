package types

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// Proposal types
const (
	ProposalTypeAddMultiStakingCoin    string = "AddMultiStakingCoin"
	ProposalTypeUpdateBondWeight       string = "UpdateBondWeight"
	ProposalTypeRemoveMultiStakingCoin string = "RemoveMultiStakingCoin"
)

// Assert module proposals implement govtypes.Content at compile-time
var (
	_ govv1beta1.Content = &AddMultiStakingCoinProposal{}
	_ govv1beta1.Content = &UpdateBondWeightProposal{}
	_ govv1beta1.Content = &RemoveMultiStakingCoinProposal{}
)

func init() {
	govv1beta1.RegisterProposalType(ProposalTypeAddMultiStakingCoin)
	govv1beta1.RegisterProposalType(ProposalTypeUpdateBondWeight)
	govv1beta1.RegisterProposalType(ProposalTypeRemoveMultiStakingCoin)

}

// NewAddMultiStakingCoinProposal returns new instance of AddMultiStakingCoinProposal
func NewAddMultiStakingCoinProposal(title, description, denom string, bondWeight sdk.Dec) govv1beta1.Content {
	return &AddMultiStakingCoinProposal{
		Title:       title,
		Description: description,
		Denom:       denom,
		BondWeight:  bondWeight,
	}
}

// GetTitle returns the title of a AddMultiStakingCoinProposal
func (abtp *AddMultiStakingCoinProposal) GetTitle() string { return abtp.Title }

// GetDescription returns the description of a AddMultiStakingCoinProposal
func (abtp *AddMultiStakingCoinProposal) GetDescription() string { return abtp.Description }

// ProposalRoute returns router key for a AddMultiStakingCoinProposal
func (*AddMultiStakingCoinProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for a AddMultiStakingCoinProposal
func (*AddMultiStakingCoinProposal) ProposalType() string {
	return ProposalTypeAddMultiStakingCoin
}

// ValidateBasic runs basic stateless validity checks
func (abtp *AddMultiStakingCoinProposal) ValidateBasic() error {
	err := govv1beta1.ValidateAbstract(abtp)
	if err != nil {
		return err
	}

	if abtp.Denom == "" {
		return sdkerrors.Wrap(ErrInvalidAddMultiStakingCoinProposal, "proposal bond token cannot be blank")
	}

	if !abtp.BondWeight.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidAddMultiStakingCoinProposal, "proposal bond token weight must be positive")
	}

	return nil
}

// String implements the Stringer interface.
func (abtp AddMultiStakingCoinProposal) String() string {
	return fmt.Sprintf("AddMultiStakingCoinProposal: Title: %s Description: %s Denom: %s TokenWeight: %s", abtp.Title, abtp.Description, abtp.Denom, abtp.BondWeight)
}

// NewUpdateBondWeightProposal returns new instance of UpdateBondWeightProposal
func NewUpdateBondWeightProposal(title, description, denom string, bondWeight sdk.Dec) govv1beta1.Content {
	return &UpdateBondWeightProposal{
		Title:             title,
		Description:       description,
		Denom:             denom,
		UpdatedBondWeight: bondWeight,
	}
}

// GetTitle returns the title of a UpdateBondWeightProposal
func (cbtp *UpdateBondWeightProposal) GetTitle() string { return cbtp.Title }

// GetDescription returns the description of a UpdateBondWeightProposal
func (cbtp *UpdateBondWeightProposal) GetDescription() string { return cbtp.Description }

// ProposalRoute returns router key for a UpdateBondWeightProposal
func (*UpdateBondWeightProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for a UpdateBondWeightProposal
func (*UpdateBondWeightProposal) ProposalType() string {
	return ProposalTypeUpdateBondWeight
}

// String implements the Stringer interface.
func (cbtp UpdateBondWeightProposal) String() string {
	return fmt.Sprintf("UpdateBondWeightProposal: Title: %s Description: %s Denom: %s TokenWeight: %s", cbtp.Title, cbtp.Description, cbtp.Denom, cbtp.UpdatedBondWeight)
}

// ValidateBasic runs basic stateless validity checks
func (cbtp *UpdateBondWeightProposal) ValidateBasic() error {
	err := govv1beta1.ValidateAbstract(cbtp)
	if err != nil {
		return err
	}

	if cbtp.Denom == "" {
		return sdkerrors.Wrap(ErrInvalidUpdateBondWeightProposal, "proposal bond token cannot be blank")
	}

	if !cbtp.UpdatedBondWeight.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidUpdateBondWeightProposal, "proposal bond token weight must be positive")
	}

	return nil
}

// NewRemoveMultiStakingCoinProposal returns new instance of RemoveMultiStakingCoinProposal
func NewRemoveMultiStakingCoinProposal(title, description, denom string) govv1beta1.Content {
	return &RemoveMultiStakingCoinProposal{
		Title:       title,
		Description: description,
		Denom:       denom,
	}
}

// GetTitle returns the title of a RemoveMultiStakingCoinProposal
func (abtp *RemoveMultiStakingCoinProposal) GetTitle() string { return abtp.Title }

// GetDescription returns the description of a RemoveMultiStakingCoinProposal
func (abtp *RemoveMultiStakingCoinProposal) GetDescription() string { return abtp.Description }

// ProposalRoute returns router key for a RemoveMultiStakingCoinProposal
func (*RemoveMultiStakingCoinProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for a RemoveMultiStakingCoinProposal
func (*RemoveMultiStakingCoinProposal) ProposalType() string {
	return ProposalTypeAddMultiStakingCoin
}

// ValidateBasic runs basic stateless validity checks
func (abtp *RemoveMultiStakingCoinProposal) ValidateBasic() error {
	err := govv1beta1.ValidateAbstract(abtp)
	if err != nil {
		return err
	}

	if abtp.Denom == "" {
		return sdkerrors.Wrap(ErrInvalidRemoveMultiStakingCoinProposal, "proposal bond token cannot be blank")
	}

	return nil
}

// String implements the Stringer interface.
func (abtp RemoveMultiStakingCoinProposal) String() string {
	return fmt.Sprintf("AddMultiStakingCoinProposal: Title: %s Description: %s Denom: %s", abtp.Title, abtp.Description, abtp.Denom)
}
