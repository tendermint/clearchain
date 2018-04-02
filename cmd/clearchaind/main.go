package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/clearchain"
	"github.com/tendermint/clearchain/app"
	"github.com/tendermint/clearchain/types"
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-crypto/keys"
	"github.com/tendermint/go-crypto/keys/words"
	"github.com/tendermint/tmlibs/cli"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

const (
	defaultClearingHouseName = "ClearingHouse"
	defaultConfigBaseDir     = ".clearchaind"
	flagClearingHouseName    = "clearing-house-name"
)

var (
	// clearchaindCmd is the entry point for this binary
	clearchaindCmd = &cobra.Command{
		Use:   "clearchaind",
		Short: "Clearchain Daemon (server)",
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the app version",
		Run:   doVersionCmd,
	}
)

func doVersionCmd(cmd *cobra.Command, args []string) {
	v := clearchain.Version
	if len(v) == 0 {
		fmt.Fprintln(os.Stderr, "unset")
		return
	}
	fmt.Fprintln(os.Stderr, v)
}

func initCommand(logger log.Logger) *cobra.Command {
	cmd := server.InitCmd(defaultOptions, logger)
	cmd.Flags().String(flagClearingHouseName, defaultClearingHouseName, "Clearing House name")
	cmd.Args = cobra.MaximumNArgs(1)
	return cmd
}

func main() {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stderr)).With("module", "main")
	initCmd := server.InitCmd(defaultOptions, logger)
	initCmd.Flags().String(flagClearingHouseName, defaultClearingHouseName, "Clearing House name")
	initCmd.Args = cobra.MaximumNArgs(1)
	clearchaindCmd.AddCommand(
		initCommand(logger),
		server.StartCmd(generateApp, logger),
		server.UnsafeResetAllCmd(logger),
		server.ShowNodeIdCmd(logger),
		server.ShowValidatorCmd(logger),
		versionCmd,
	)
	// prepare and add flags
	rootDir := os.ExpandEnv(defaultConfigBaseDir)
	executor := cli.PrepareBaseCmd(clearchaindCmd, "CC", rootDir)
	executor.Execute()
}

func readOrGenerateKey(args []string) (crypto.PubKey, string, error) {
	if len(args) != 0 { // user has given a hexadecimal pubkey on the command line
		pub, err := types.PubKeyFromHexString(args[0])
		if err != nil {
			return crypto.PubKey{}, "", err
		}
		return pub, "", nil
	}
	return generateKey()
}

// defaultOptions sets up the app_options for the
// default genesis file
func defaultOptions(args []string) (json.RawMessage, string, cmn.HexBytes, error) {
	pub, secret, err := readOrGenerateKey(args)
	if err != nil {
		return nil, "", nil, err
	}
	opts := fmt.Sprintf(`{
      "ch_admin": {
		"public_key": "%s",
		"entity_name": "%s"
	  }
	}`, hex.EncodeToString(pub.Bytes()), viper.GetString(flagClearingHouseName))
	return json.RawMessage(opts), secret, pub.Address(), nil
}

func generateApp(rootDir string, logger log.Logger) (abci.Application, error) {
	db, err := dbm.NewGoLevelDB("clearchain", rootDir)
	if err != nil {
		return nil, err
	}
	bapp := app.NewClearchainApp(logger, db)
	return bapp, nil
}

func generateKey() (crypto.PubKey, string, error) {
	// construct an in-memory key store
	codec, err := words.LoadCodec("english")
	if err != nil {
		return crypto.PubKey{}, "", err
	}
	keybase := keys.New(
		dbm.NewMemDB(),
		codec,
	)

	// generate a private key, with recovery phrase
	info, secret, err := keybase.Create("name", "pass", keys.AlgoEd25519)
	if err != nil {
		return crypto.PubKey{}, "", err
	}

	return info.PubKey, secret, nil
}
