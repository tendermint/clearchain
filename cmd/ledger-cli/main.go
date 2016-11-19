package main

import (
	"flag"
	"fmt"
	"github.com/tendermint/clearchain/client"
	crypto "github.com/tendermint/go-crypto"
)

const (
	GenerateKey string = "generateKey"
	DecodeKey string = "decodeKey"
)

func main() {
	flag.Bool(GenerateKey, false, "Generate secure private and public keys")
	flag.String(DecodeKey, "", "Prints binary representation of provided key")

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	flag.Visit(flagHandler)
}

func flagHandler(flag *flag.Flag) {
	switch flag.Name {
	case GenerateKey:
		generateKeys()
		return
	case DecodeKey:
		decodeKey(flag.Value.String())
		return
	default:
		panic(fmt.Sprintf(":( Unimplemented flag: %v", flag.Name))
		return
	}
}

func decodeKey(key string) {
	fmt.Println("Decoded Key: ", client.Decode(key))
}

func generateKeys() {
	privateKey := crypto.GenPrivKeyEd25519()
	publicKey := privateKey.PubKey()

	encodedPrivateKey := client.Encode(privateKey.Bytes())
	fmt.Println("EncodedPrivateKey: ", encodedPrivateKey)

	encodedPublicKey := client.Encode(publicKey.Bytes())
	fmt.Println("EncodedPublicKey: ", encodedPublicKey)
}
