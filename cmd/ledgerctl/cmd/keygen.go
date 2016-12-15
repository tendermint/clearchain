package cmd

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/clearchain/client"
	crypto "github.com/tendermint/go-crypto"
)

var (
	flagWithSecret bool
	flagOutputFile string
)

func init() {
	keygenCmd.Flags().BoolVar(&flagWithSecret, "with-secret", false, "Generate keys from a secret")
	RootCmd.AddCommand(keygenCmd)
}

var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate secure private and public key pair",
	Run: func(cmd *cobra.Command, args []string) {
		var privateKey crypto.PrivKeyEd25519
		if flagWithSecret {
			fmt.Fprintf(os.Stderr, "Generating public/private key pair.\nEnter a secret: ")
			userInput, err := ReadLineBytes(os.Stdin)
			if err != nil {
				log.Fatal(err)
			}
			secret, err := bcrypt.GenerateFromPassword(userInput, 4)
			if err != nil {
				log.Fatal(err)
			}
			privateKey = crypto.GenPrivKeyEd25519FromSecret(secret)
		} else {
			privateKey = crypto.GenPrivKeyEd25519()
		}

		fmt.Println("PrivateKey     : ", client.Encode(privateKey.Bytes()))
		fmt.Println("PublicKeyAddr  : ", client.Encode(privateKey.PubKey().Address()))
		fmt.Println("PublicKey      : ", client.Encode(privateKey.PubKey().Bytes()))
	},
}
