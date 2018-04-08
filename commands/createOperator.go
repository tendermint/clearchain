package commands

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/builder"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/clearchain/types"
)

// GetCreateOperatorTxCmd returns a createOperatorTxCmd.
func GetCreateOperatorTxCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := Commander{Cdc: cdc}
	cmd := &cobra.Command{
		Use:   "create-operator",
		Short: "Create and sign a CreateOperatorTx",
		RunE:  cmdr.createOperatorTxCmd,
		Args:  cobra.ExactArgs(1),
	}
	cmd.Flags().String(flagPubKeyFile, "", "Load new asset's pubkey from file")
	cmd.MarkFlagRequired(flagPubKeyFile)
	return cmd
}

func (c Commander) createOperatorTxCmd(cmd *cobra.Command, args []string) error {
	name := viper.GetString(flagName)
	creatorInfo, err := getKey(name)
	if err != nil {
		return fmt.Errorf("getKey(): %v", err)
	}
	pub, err := pubKeyFromFile()
	if err != nil {
		return err
	}
	msg := types.NewCreateOperatorMsg(creatorInfo.PubKey.Address(), pub)
	res, err := builder.SignBuildBroadcast(name, msg, c.Cdc)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil

}
