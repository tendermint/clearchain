package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/commands"
	"github.com/spf13/cobra"
	"github.com/tendermint/clearchain/commands"
	"github.com/tendermint/clearchain/types"
	"github.com/tendermint/tmlibs/cli"
)

const (
	defaultConfigBaseDir = ".clearchainctl"
)

var (
	clearchainctlCmd = &cobra.Command{
		Use:   "clearchainctl",
		Short: "Clearchain light-client",
	}
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := types.MakeCodec()
	rpc.AddCommands(clearchainctlCmd)
	clearchainctlCmd.AddCommand(client.LineBreak)
	tx.AddCommands(clearchainctlCmd, cdc)
	clearchainctlCmd.AddCommand(client.LineBreak)

	// add clearchain-specific commands
	clearchainctlCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("main", cdc, types.GetAccountDecoder(cdc)),
		)...)
	clearchainctlCmd.AddCommand(
		client.PostCommands(
			commands.GetCreateAdminTxCmd(cdc),
		)...)
	clearchainctlCmd.AddCommand(commands.GetExportPubCmd(cdc))
	//clearchainctlCmd.AddCommand(commands.GetImportPubCmd(cdc))

	// add proxy, version and key info
	clearchainctlCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(cdc),
		keys.Commands(),
		client.LineBreak,
		commands.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(clearchainctlCmd, "CC", os.ExpandEnv(defaultConfigBaseDir))
	executor.Execute()
}
