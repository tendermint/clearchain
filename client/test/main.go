package main

import (
	"fmt"
	"github.com/tendermint/clearchain/client"
	"github.com/tendermint/clearchain/types"
	"github.com/tendermint/go-crypto"
)

func main() {
	serverAddress := "tcp://127.0.0.1:46658"
	chainID := "test_chain_id"

	client.SetChainID(chainID)
	client.StartClient(serverAddress)

	privKeyBytes := []byte{1, 52, 87, 91, 9, 73, 233, 187, 205, 69, 195, 81, 79, 241, 12, 155, 41, 163, 102, 240, 6, 176, 182, 105, 229, 175, 121, 183, 209, 203, 228, 212, 97, 132, 56, 120, 185, 50, 238, 73, 8, 164, 45, 36, 191, 252, 3, 158, 184, 223, 172, 212, 52, 12, 131, 52, 35, 16, 104, 32, 148, 4, 127, 175, 171}
	//	fmt.Println("privKeyBytes: ", string(privKeyBytes))

	privKey, _ := crypto.PrivKeyFromBytes(privKeyBytes)

	fmt.Println("Account IDs:")
	var accountsRequested []string = client.GetAllAccounts(privKey).Accounts
	//	[]string{"1d2df1ae-accb-11e6-bbbb-00ff5244ae7f"}

	for _, account := range accountsRequested {
		fmt.Println("\t", account)
	}

	var accounts []*types.Account = client.GetAccounts(privKey, accountsRequested).Account

	fmt.Println("accounts returned:")
	for _, account := range accounts {
		fmt.Println("\t", account)
	}
}

//
