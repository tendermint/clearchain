package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cosmos/cosmos-sdk/client/builder"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/clearchain/types"
	crypto "github.com/tendermint/go-crypto"
)

const (
	flagName       = "name"
	flagPubKeyFile = "pubkey-file"
	flagEntityName = "entity-name"
	flagEntityType = "entity-type"
	flagSequence   = "seq"
)

type Commander struct {
	Cdc *wire.Codec
}

// GetCreateAdminTxCmd returns a CreateAdminTxCmd.
func GetCreateAdminTxCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := Commander{Cdc: cdc}
	cmd := &cobra.Command{
		Use:   "create-admin",
		Short: "Create and sign a CreateAdminTx",
		RunE:  cmdr.createAdminTxCmd,
	}
	cmd.Flags().String(flagPubKeyFile, "", "Load new asset's pubkey from file")
	cmd.Flags().String(flagEntityName, "", "New admin's entity name")
	cmd.Flags().String(flagEntityType, "", "New admin's entity type (ch|gcm|icm|custodian)")
	cmd.MarkFlagRequired(flagPubKeyFile)
	cmd.MarkFlagRequired(flagEntityName)
	cmd.MarkFlagRequired(flagEntityType)
	cmd.MarkFlagFilename(flagPubKeyFile)
	return cmd
}

func (c Commander) createAdminTxCmd(cmd *cobra.Command, args []string) error {
	name := viper.GetString(flagName)
	info, err := getKey(name)
	if err != nil {
		return fmt.Errorf("getKey(): %v", err)
	}
	pub, err := pubKeyFromFile()
	if err != nil {
		return err
	}
	creator := info.PubKey.Address()
	msg := types.NewCreateAdminMsg(creator, pub, viper.GetString(flagEntityName), viper.GetString(flagEntityType))
	if err != nil {
		return err
	}

	res, err := builder.SignBuildBroadcast(name, msg, c.Cdc)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil

}

func pubKeyFromFile() (key crypto.PubKey, err error) {
	filename := viper.GetString(flagPubKeyFile)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return key, fmt.Errorf("Error reading from %s: %v", filename, err)
	}
	err = key.UnmarshalJSON(bytes)
	return
}
