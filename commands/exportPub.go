package commands

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/keys"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
)

const (
	flagFormat = "format"
)

func GetExportPubCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "export-pub",
		RunE: exportPubCmd,
		Args: cobra.ExactArgs(1),
	}
	cmd.Flags().String(flagFormat, "", "Alternative output format (armor|json)")
	return cmd
}

func exportPubCmd(cmd *cobra.Command, args []string) error {
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return err
	}
	key, err := keybase.Get(args[0])
	if err != nil {
		return err
	}
	bz, err := key.PubKey.MarshalJSON()
	if err != nil {
		return err
	}
	fmt.Println(string(bz))
	return nil
}
