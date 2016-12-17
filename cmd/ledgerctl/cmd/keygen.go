package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tendermint/clearchain/client"
	crypto "github.com/tendermint/go-crypto"
	"golang.org/x/crypto/bcrypt"
)

var (
	flagWithSecret bool
	flagOutputFile string
)

func init() {
	keygenCmd.Flags().BoolVar(&flagWithSecret, "with-secret", false, "Generate keys from a secret")
	keygenCmd.Flags().StringVarP(&flagOutputFile, "output-file", "O", "",
		`Write the private key to the given file;
			     the public key will be saved with the .pub extension`)
	RootCmd.AddCommand(keygenCmd)
}

var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate secure private and public key pair",
	Run: func(cmd *cobra.Command, args []string) {
		var privKey crypto.PrivKey
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
			privKey = crypto.GenPrivKeyEd25519FromSecret(secret)
		} else {
			privKey = crypto.GenPrivKeyEd25519()
		}
		pubKey := privKey.PubKey()

		fmt.Println("Fingerprint:\n", client.Encode(pubKey.Address()))
		if len(flagOutputFile) == 0 {
			fmt.Println("\nPrivateKey:\n", client.Encode(privKey.Bytes()))
			fmt.Println("\nPublicKey:\n", client.Encode(pubKey.Bytes()))
		} else {
			privKeyFile, err := os.Create(flagOutputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
			privKeyFile.WriteString(client.Encode(privKey.Bytes()))

			pubKeyFile, err := os.Create(strings.Join([]string{flagOutputFile, "pub"}, "."))
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
			pubKeyFile.WriteString(client.Encode(pubKey.Bytes()))
		}
	},
}
