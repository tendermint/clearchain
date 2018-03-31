package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	wire "github.com/tendermint/go-wire"
)

const (
	flagName = "name"
)

func GetPubToHexCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "pub2hex",
		RunE: pubToHexCmd,
	}
	cmd.Flags().String(flagName, "", "Account's pubkey")
	return cmd
}

func pubToHexCmd(cmd *cobra.Command, args []string) error {
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return err
	}
	name := viper.GetString(flagName)
	info, err := keybase.Get(name)
	if err != nil {
		return errors.Errorf("No key for: %s", name)
	}
	fmt.Println(hex.EncodeToString(info.PubKey.Bytes()))
	return nil
}
