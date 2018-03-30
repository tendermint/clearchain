package commands

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/clearchain/types"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
)

// GetCreateAssetAccountTxCmd returns a CreateAssetAccountTxCmd.
func GetCreateAssetAccountTxCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{cdc}
	cmd := &cobra.Command{
		Use:   "createasset",
		Short: "Create and sign a CreateAssetAccountTxCmd",
		RunE:  cmdr.createAssetAccountTxCmd,
	}
	cmd.Flags().String(flagPubKey, "", "New asset account's pubkey")
	cmd.Flags().Int64(flagSequence, 0, "Sequence number")
	return cmd
}

func (c commander) createAssetAccountTxCmd(cmd *cobra.Command, args []string) error {
	txBytes, err := c.buildCreateAssetAccountTx()
	if err != nil {
		return err
	}

	res, err := client.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
	return nil
}

func (c commander) buildCreateAssetAccountTx() ([]byte, error) {
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	name := viper.GetString(client.FlagName)
	info, err := keybase.Get(name)
	if err != nil {
		return nil, errors.Errorf("No key for: %s", name)
	}
	creator := info.PubKey.Address()

	msg, err := buildCreateAssetAccountMsg(creator)
	if err != nil {
		return nil, err
	}

	// sign and build
	bz := msg.GetSignBytes()
	buf := client.BufferStdin()
	prompt := fmt.Sprintf("Password to sign with '%s':", name)
	passphrase, err := client.GetPassword(prompt, buf)
	if err != nil {
		return nil, err
	}
	sig, pubkey, err := keybase.Sign(name, passphrase, bz)
	if err != nil {
		return nil, err
	}
	sigs := []sdk.StdSignature{{
		PubKey:    pubkey,
		Signature: sig,
		Sequence:  viper.GetInt64(flagSequence),
	}}

	// marshal bytes
	tx := sdk.NewStdTx(msg, sigs)

	txBytes, err := c.cdc.MarshalBinary(tx)
	if err != nil {
		return nil, err
	}
	return txBytes, nil
}

func buildCreateAssetAccountMsg(creator crypto.Address) (sdk.Msg, error) {
	// parse new account pubkey
	pubKey, err := types.PubKeyFromHexString(viper.GetString(flagPubKey))
	if err != nil {
		return nil, err
	}
	msg := types.NewCreateAssetAccountMsg(creator, pubKey)
	return msg, nil
}
