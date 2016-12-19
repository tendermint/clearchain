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
		`Write the private key binary format to the given file;
			     the public key and the address files will be saved
			     with .pub and .addr extensions respectively`)
	RootCmd.AddCommand(keygenCmd)
}

var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate secure private and public key pair",
	Run: func(cmd *cobra.Command, args []string) {
		var privKey crypto.PrivKey
		log.Println("Generating public/private key pair...")
		if flagWithSecret {
			fmt.Fprintf(os.Stderr, "Enter a secret: ")
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

		privKeyBytes := privKey.Bytes()
		pubKeyBytes := privKey.PubKey().Bytes()
		addrBytes := privKey.PubKey().Address()

		fmt.Println("\nPrivateKey:\n", client.Encode(privKeyBytes))
		fmt.Println("\nPublicKey:\n", client.Encode(pubKeyBytes))
		fmt.Println("\nAddress:\n", client.Encode(addrBytes))
		if len(flagOutputFile) != 0 {
			mustWriteToFile(mustCreateFile(flagOutputFile), privKeyBytes)
			mustWriteToFile(mustCreateFile(strings.Join([]string{flagOutputFile, "pub"}, ".")), pubKeyBytes)
			mustWriteToFile(mustCreateFile(strings.Join([]string{flagOutputFile, "addr"}, ".")), addrBytes)
		}
	},
}

func mustCreateFile(filename string) *os.File {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func mustWriteToFile(f *os.File, b []byte) {
	if _, err := f.Write(b); err != nil {
		log.Fatal(err)
	}
}
