package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/clearchain/app"
	"github.com/tendermint/clearchain/commands"
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-crypto/keys"
	"github.com/tendermint/go-crypto/keys/words"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

const (
	defaultClearingHouseName = "ClearingHouse"
	defaultConfigBaseDir     = "$HOME/.clearchaind"
)

var (
	clearchaindCmd = &cobra.Command{
		Use:   "clearchaind",
		Short: "Clearchain Daemon (server)",
	}
)

func main() {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stderr)).With("module", "main")
	clearchaindCmd.AddCommand(
		server.InitCmd(defaultOptions, logger),
		server.StartCmd(generateApp, logger),
		server.UnsafeResetAllCmd(logger),
		commands.VersionCmd,
	)
	// prepare and add flags
	rootDir := os.ExpandEnv(defaultConfigBaseDir)
	executor := cli.PrepareBaseCmd(clearchaindCmd, "CC", rootDir)
	executor.Execute()
}

// defaultOptions sets up the app_options for the
// default genesis file
func defaultOptions(args []string) (json.RawMessage, error) {
	pub, secret, err := generateKey()
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "Secret phrase to access clearing house's admin account: %s\n", secret)
	opts := fmt.Sprintf(`{
		"ch_admin": {
		  "public_key": "%s",
		  "entity_name": "%s"
		}
	  }`, hex.EncodeToString(pub.Bytes()), defaultClearingHouseName)
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
		return nil, "", err
	}
	keybase := keys.New(
		dbm.NewMemDB(),
		codec,
	)
	// generate a private key, with recovery phrase
	info, secret, err := keybase.Create("name", "pass", keys.AlgoEd25519)
	if err != nil {
		return nil, "", err
	}
	return info.PubKey, secret, nil
}
