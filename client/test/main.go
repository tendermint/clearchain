package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/tendermint/clearchain/client"
	"github.com/tendermint/clearchain/types"
	"github.com/tendermint/go-crypto"
)

func main() {
	serverAddress := "tcp://127.0.0.1:46658"
	chainID := "test_chain_id"

	client.SetChainID(chainID)
	client.StartClient(serverAddress)

	privateKeyBytes := []byte{1, 52, 87, 91, 9, 73, 233, 187, 205, 69, 195, 81, 79, 241, 12, 155, 41, 163, 102, 240, 6, 176, 182, 105, 229, 175, 121, 183, 209, 203, 228, 212, 97, 132, 56, 120, 185, 50, 238, 73, 8, 164, 45, 36, 191, 252, 3, 158, 184, 223, 172, 212, 52, 12, 131, 52, 35, 16, 104, 32, 148, 4, 127, 175, 171}

	privateKey, _ := crypto.PrivKeyFromBytes(privateKeyBytes)

	userName := "userName"
	canCreate := true
	pubKey := crypto.GenPrivKeyEd25519().PubKey()
	client.CreateUser(privateKey, userName, pubKey, canCreate)

	accountID := uuid.NewV4().String()
	client.CreateAccount(privateKey, accountID)

	entityID := uuid.NewV4().String()
	entityType := types.EntityTypeCHByte
	legalEntityName := "newLegalEntityName"
	client.CreateLegalEntity(privateKey, entityID, entityType, legalEntityName)

	fmt.Println("Account IDs:")
	var accountsRequested []string = client.GetAllAccounts(privateKey).Accounts

	for _, account := range accountsRequested {
		fmt.Println("\t", account)
	}

	var accounts []*types.Account = client.GetAccounts(privateKey, accountsRequested).Account

	fmt.Println("accounts returned:")
	for _, account := range accounts {
		fmt.Println("\t", account)
	}

	legalEntityIDs := client.GetAllLegalEntities(privateKey)

	fmt.Println("legalEntity IDs:")
	for _, legalEntityID := range legalEntityIDs.Ids {
		fmt.Println("\t", legalEntityID)

	}

	legalEntities := client.GetLegalEntities(privateKey, legalEntityIDs.Ids).LegalEntities
	fmt.Println("legalEntities returned:")
	for _, legalEntity := range legalEntities {
		fmt.Println("\t", legalEntity)
	}
}

//
