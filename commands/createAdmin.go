package commands

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/builder"
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/clearchain/types"
)

const (
	flagPubKey     = "pubkey"
	flagEntityName = "entityname"
	flagEntityType = "entitytype"
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
		Args:  cobra.ExactArgs(1),
	}
	cmd.Flags().String(flagPubKey, "", "New admin's pubkey")
	cmd.Flags().String(flagEntityName, "", "New admin's entity name")
	cmd.Flags().String(flagEntityType, "", "New admin's entity type (ch|gcm|icm|custodian)")
	cmd.Flags().Int64(flagSequence, 0, "Sequence number")
	return cmd
}

func (c Commander) createAdminTxCmd(cmd *cobra.Command, args []string) error {
	name := args[0]
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil
	}
	info, err := keybase.Get(name)
	if err != nil {
		return err
	}
	creator := info.PubKey.Address()
	msg, err := BuildCreateAdminMsg(creator, viper.GetString(flagEntityName), viper.GetString(flagEntityType), viper.GetString(flagPubKey))
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

// BuildCreateAdminMsg makes a new CreateAdminMsg.
func BuildCreateAdminMsg(creator sdk.Address, entityName, entityType, pubKey string) (sdk.Msg, error) {
	// parse new account pubkey
	pub, err := types.PubKeyFromHexString(pubKey)
	if err != nil {
		return nil, err
	}
	return types.NewCreateAdminMsg(creator, pub, entityName, entityType), nil
}
