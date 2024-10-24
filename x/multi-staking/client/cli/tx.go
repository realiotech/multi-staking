package cli

import (
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"github.com/spf13/cobra"

	"cosmossdk.io/core/address"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
)

// NewTxCmd returns a root CLI command handler for all x/exp transaction commands.
func NewTxCmd(valAddrCodec, ac address.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "multi-staking transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		cli.NewCreateValidatorCmd(valAddrCodec),
		cli.NewEditValidatorCmd(valAddrCodec),
		cli.NewDelegateCmd(valAddrCodec, ac),
		cli.NewRedelegateCmd(valAddrCodec, ac),
		cli.NewUnbondCmd(valAddrCodec, ac),
		cli.NewCancelUnbondingDelegation(valAddrCodec, ac),
	)

	return txCmd
}
