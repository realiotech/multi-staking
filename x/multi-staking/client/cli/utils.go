package cli

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func ParseAddMultiStakingCoinProposal(cdc *codec.LegacyAmino, proposalFile string) (types.AddMultiStakingCoinProposal, error) {
	proposal := types.AddMultiStakingCoinProposal{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
