package main

import (
	"fmt"
	"os"

	"github.com/realio-tech/multi-staking-module/test/simapp"
	"github.com/realio-tech/multi-staking-module/test/simapp/simd/cmd"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "simd", simapp.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err) //nolint
		os.Exit(1)
	}
}
