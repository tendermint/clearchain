package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/clearchain"
	"github.com/tendermint/clearchain/app"
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-crypto/keys"
	"github.com/tendermint/go-crypto/keys/words"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

const (
	defaultClearingHouseName = "ClearingHouse"
	defaultConfigBaseDir     = ".clearchaind"
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

func main() {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stderr)).With("module", "main")
	clearchaindCmd.AddCommand(
		server.InitCmd(defaultOptions, logger),
		server.StartCmd(generateApp, logger),
		server.UnsafeResetAllCmd(logger),
		server.ShowNodeIdCmd(logger),
		versionCmd,
	)
	// prepare and add flags
	rootDir := os.ExpandEnv(defaultConfigBaseDir)
	executor := cli.PrepareBaseCmd(clearchaindCmd, "CC", rootDir)
	executor.Execute()
}

// defaultOptions sets up the app_options for the
// default genesis file
func defaultOptions(args []string) (json.RawMessage, error) {
	var pubHex string
	if len(args) != 0 { // user has given a hexadecimal pubkey on the command line
		pubHex = args[0]
	} else {
		pub, secret, err := generateKey()
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "Secret phrase to access clearing house's admin account: %s\n", secret)
		pubHex = hex.EncodeToString(pub.Bytes())
	}
	opts := fmt.Sprintf(`{
      "ch_admin": {
		"public_key": "%s",
		"entity_name": "%s"
	  }
	}`, pubHex, defaultClearingHouseName)
	return json.RawMessage(opts), nil
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
