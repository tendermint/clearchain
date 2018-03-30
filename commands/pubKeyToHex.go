package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	wire "github.com/tendermint/go-wire"
)

func GetPubToHexCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "pub2hex",
		RunE: pubToHexCmd,
	}
	return cmd
}

func pubToHexCmd(cmd *cobra.Command, args []string) error {
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return err
	}
	if len(args) < 1 {
		return fmt.Errorf("insufficient arguments")
	}
	name := args[0]
	info, err := keybase.Get(name)
	if err != nil {
		return errors.Errorf("No key for: %s", name)
	}
	fmt.Println(hex.EncodeToString(info.PubKey.Bytes()))
	return nil
}
