package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/keys"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	flagFormat = "format"
)

func GetExportPubCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "export-pub",
		RunE: exportPubCmd,
	}
	cmd.Flags().String(flagFormat, "", "Alternative output format (armor|json)")
	return cmd
}

func exportPubCmd(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("insufficient arguments")
	}
	if len(args) > 1 {
		return fmt.Errorf("too many arguments: %v", args)
	}
	name := args[0]
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return err
	}
	if viper.GetString(flagFormat) == "armor" {
		out, err := keybase.Export(name)
		if err != nil {
			return err
		}
		fmt.Println(out)
		return nil
	}
	info, err := keybase.Get(name)
	if err != nil {
		return errors.Errorf("No key for: %s", name)
	}
	if viper.GetString(flagFormat) == "json" {
		out, err := info.PubKey.MarshalJSON()
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	}
	fmt.Println(hex.EncodeToString(info.PubKey.Bytes()))
	return nil
}
