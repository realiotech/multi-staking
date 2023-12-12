package cli

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

const flagGenesisTime = "genesis-time"

// MigrateGenesisCmd returns a command to execute genesis state migration.
func MigrateStakingGenesisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-staking [target-version] [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "override genesis_time with this flag")
	cmd.Flags().String(flags.FlagChainID, "", "override chain_id with this flag")

	return cmd
}
