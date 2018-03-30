package main

import (
	"errors"
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

// gaiacliCmd is the entry point for this binary
var (
	clearchainctlCmd = &cobra.Command{
		Use:   "clearchainctl",
		Short: "Clearchain light-client",
	}
)

func todoNotImplemented(_ *cobra.Command, _ []string) error {
	return errors.New("TODO: Command not yet implemented")
}

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := types.MakeCodec()

	// add standard rpc, and tx commands
	rpc.AddCommands(clearchainctlCmd)
	clearchainctlCmd.AddCommand(client.LineBreak)
	tx.AddCommands(clearchainctlCmd, cdc)
	clearchainctlCmd.AddCommand(client.LineBreak)

	// add query/post commands (custom to binary)
	clearchainctlCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("main", cdc, types.GetParseAccount(cdc)),
		)...)
	clearchainctlCmd.AddCommand(commands.GetPubToHexCmd(cdc))
	clearchainctlCmd.AddCommand(
		client.PostCommands(
			commands.GetCreateAdminTxCmd(cdc),
			commands.GetCreateOperatorTxCmd(cdc),
			commands.GetCreateAssetAccountTxCmd(cdc),
			//			bankcmd.SendTxCmd(cdc),
		)...)

	// add proxy, version and key info
	clearchainctlCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(),
		keys.Commands(),
		client.LineBreak,
		commands.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(clearchainctlCmd, "CC", os.ExpandEnv(defaultConfigBaseDir))
	executor.Execute()
}
