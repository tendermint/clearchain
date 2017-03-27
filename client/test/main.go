package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/tendermint/clearchain/client"
	"github.com/tendermint/clearchain/types"
	"github.com/tendermint/go-crypto"
)

// Tendermint core  must be running
func main() {
	serverAddress := "127.0.0.1:46657"
	chainID := "test_chain_id"

	client.SetChainID(chainID)
	client.StartClient(serverAddress)

	privateKeyBytes := []byte{1, 52, 87, 91, 9, 73, 233, 187, 205, 69, 195, 81, 79, 241, 12, 155, 41, 163, 102, 240, 6, 176, 182, 105, 229, 175, 121, 183, 209, 203, 228, 212, 97, 132, 56, 120, 185, 50, 238, 73, 8, 164, 45, 36, 191, 252, 3, 158, 184, 223, 172, 212, 52, 12, 131, 52, 35, 16, 104, 32, 148, 4, 127, 175, 171}
	//ATRXWwlJ6bvNRcNRT/EMmymjZvAGsLZp5a95t9HL5NRhhDh4uTLuSQikLSS//AOeuN+s1DQMgzQjEGgglAR/r6s=

	privateKey, _ := crypto.PrivKeyFromBytes(privateKeyBytes) //User{b40cbf4e-5923-4ccd-beec-e22a9117b91b "Name" 31}

	userName := "userName"
	canCreate := true
	pubKey := crypto.GenPrivKeyEd25519().PubKey()
	client.CreateUser(privateKey, userName, pubKey, canCreate)

	accountID := uuid.NewV4().String()
	client.CreateAccount(privateKey, accountID)

	entityID := uuid.NewV4().String()
	parentID := uuid.NewV4().String()
	entityType := types.EntityTypeCHByte
	legalEntityName := "newLegalEntityName"
	client.CreateLegalEntity(privateKey, entityID, entityType, legalEntityName, parentID)

	fmt.Println("Account IDs:")
	var accountsRequested []string = client.GetAllAccounts().Accounts

	for _, account := range accountsRequested {
		fmt.Println("\t", account)
		
		var accounts []*types.Account = client.GetAccount(account).Account
		fmt.Println("accounts returned:")
		for _, accountRes := range accounts {
			fmt.Println("\t", accountRes)
		}
	}

	legalEntityIDs := client.GetAllLegalEntities()

	fmt.Println("legalEntity IDs:")
	for _, legalEntityID := range legalEntityIDs.Ids {
		fmt.Println("\t", legalEntityID)
		
		legalEntities := client.GetLegalEntity(legalEntityID).LegalEntities
		fmt.Println("legalEntities returned:")
		fmt.Println(legalEntities)
		for _, legalEntity := range legalEntities {
			fmt.Println("\t", legalEntity)
		}
	}

	//	serverAddress: tcp://127.0.0.1:46658
	//	chainID: test_chain_id
	//privateKey: ATRXWwlJ6bvNRcNRT/EMmymjZvAGsLZp5a95t9HL5NRhhDh4uTLuSQikLSS//AOeuN+s1DQMgzQjEGgglAR/r6s=
	senderID := "1d2df1ae-accb-11e6-bbbb-00ff5244ae7f"
	recipientID := "6b6d3a08-5527-4955-b4fd-f5ba7e083548"
	counterSignerAddresses := [][]byte{}
	amount := 10000
	currency := "EUR"
	client.TransferMoney(privateKey, senderID, recipientID, counterSignerAddresses, int64(amount), currency)
}

//
