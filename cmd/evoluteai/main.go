package main

import (
	"fmt"
	"os"

	"github.com/evoluteai-network/evoluteai-chain/app"
	"github.com/evoluteai-network/evoluteai-chain/cmd/evoluteaid/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {

	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
