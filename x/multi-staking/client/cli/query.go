package cli

import (
	"fmt"

	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	return cmd
}
