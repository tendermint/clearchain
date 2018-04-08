package commands

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/builder"
	"github.com/cosmos/cosmos-sdk/client/keys"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/clearchain/types"
	cryptokeys "github.com/tendermint/go-crypto/keys"
)

// GetCreateAssetAccountTxCmd returns a createAssetAccountTxCmd.
func GetCreateAssetAccountTxCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := Commander{Cdc: cdc}
	cmd := &cobra.Command{
		Use:   "create-asset",
		Short: "Create and sign a CreateAssetAccountTx",
		RunE:  cmdr.createAssetAccountTxCmd,
	}
	cmd.Flags().String(flagPubKeyFile, "", "Load new asset's pubkey from file")
	cmd.MarkFlagRequired(flagPubKeyFile)
	cmd.MarkFlagFilename(flagPubKeyFile)
	return cmd
}

func (c Commander) createAssetAccountTxCmd(cmd *cobra.Command, args []string) error {
	name := viper.GetString(flagName)
	creatorInfo, err := getKey(name)
	if err != nil {
		return fmt.Errorf("getKey(): %v", err)
	}
	pub, err := pubKeyFromFile()
	if err != nil {
		return err
	}
	msg := types.NewCreateAssetAccountMsg(creatorInfo.PubKey.Address(), pub)
	res, err := builder.SignBuildBroadcast(name, msg, c.Cdc)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil

}

func getKey(name string) (cryptokeys.Info, error) {
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return cryptokeys.Info{}, err
	}
	creatorInfo, err := keybase.Get(name)
	if err != nil {
		return cryptokeys.Info{}, fmt.Errorf("couldn't retrieve key name %q: %v", name, err)
	}
	return creatorInfo, nil
}
