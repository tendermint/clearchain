package client

import (
	"encoding/json"
	"fmt"

	"github.com/tendermint/clearchain/types"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-logger"
	"github.com/tendermint/go-wire"
	"github.com/tendermint/tmsp/client"
	tmsp "github.com/tendermint/tmsp/types"
)

var log = logger.New("module", "client")

// AccountsReturned defines the attributes of response's payload
type AccountsReturned struct {
	Account []*types.Account `json:"accounts"`
}

type LegalEntitiesReturned struct {
	LegalEntities []*types.LegalEntity `json:"legalEntities"`
}

var chainID string
var client tmspcli.Client

// SetChainID assigns and initializes the chain's ID
func SetChainID(id string) {
	chainID = id
}

func CreateUser(privateKey crypto.PrivKey,
	newUsersName string,
	newUsersPubKey crypto.PubKey,
	newUserCanCreateLegalEntity bool) {
	tx := &types.CreateUserTx{Address: privateKey.PubKey().Address(),
		Name:      newUsersName,
		PubKey:    newUsersPubKey,
		CanCreate: newUserCanCreateLegalEntity}

	res := sendAppendTx(privateKey, tx)

	if res.IsErr() {
		panic(fmt.Sprintf("Wrong response from server: %v", res))
	} else {
		Commit(client)
	}
}

func CreateAccount(privateKey crypto.PrivKey,
	accountID string) {
	tx := &types.CreateAccountTx{Address: privateKey.PubKey().Address(),
		AccountID: accountID}

	res := sendAppendTx(privateKey, tx)

	if res.IsErr() {
		panic(fmt.Sprintf("Wrong response from server: %v", res))
	} else {
		Commit(client)
	}

	log.Info("Created account with ID: " + accountID)
}

func CreateLegalEntity(privateKey crypto.PrivKey,
	entityID string, entityType byte, name string, parentID string) {
	tx := &types.CreateLegalEntityTx{Address: privateKey.PubKey().Address(),
		EntityID: entityID,
		Type:     entityType,
		Name:     name,
		ParentID: parentID}

	res := sendAppendTx(privateKey, tx)

	if res.IsErr() {
		panic(fmt.Sprintf("Wrong response from server: %v", res))
	} else {
		Commit(client)
	}
	log.Info("Created legal entity with ID: " + entityID)
}

//func TransferMoney(privateKey crypto.PrivKey) {
//	//	Address: privateKey.PubKey().Address(),
//	tx := &types.TransferTx{}
//
//	res := sendAppendTx(privateKey, tx)
//
//	if res.IsErr() {
//		panic(fmt.Sprintf("Wrong response from server: %v", res))
//	} else {
//		Commit(client)
//	}
//	log.Info("Created transfer entry")
//}

// GetAccounts makes a request to the ledger to returns a set of accounts
func GetAccounts(privateKey crypto.PrivKey, accountsRequested []string) (returned AccountsReturned) {
	tx := &types.AccountQueryTx{Accounts: accountsRequested,
		Address: privateKey.PubKey().Address()}

	res := sendQuery(privateKey, tx)

	err := json.Unmarshal(res.Data, &returned)
	if err != nil {
		panic(fmt.Sprintf("JSON unmarshal for message %v failed with: %v ", returned, err))
	}

	return
}

// AccountIndex makes a request to the ledger to returns all account IDs
func GetAllAccounts(privateKey crypto.PrivKey) (returned types.AccountIndex) {
	tx := &types.AccountIndexQueryTx{Address: privateKey.PubKey().Address()}

	res := sendQuery(privateKey, tx)

	err := json.Unmarshal(res.Data, &returned)
	if err != nil {
		panic(fmt.Sprintf("JSON unmarshal for message %v failed with: %v ", res, err))
	}
	return
}

func GetAllLegalEntities(privateKey crypto.PrivKey) (returned types.LegalEntityIndex) {
	tx := &types.LegalEntityIndexQueryTx{Address: privateKey.PubKey().Address()}

	res := sendQuery(privateKey, tx)

	err := json.Unmarshal(res.Data, &returned)
	if err != nil {
		panic(fmt.Sprintf("JSON unmarshal for message %v failed with: %v ", res, err))
	}
	return
}

func sendQuery(privateKey crypto.PrivKey, tx types.SignedTx) tmsp.Result {
	return sendToTendermint(privateKey, tx, client.QuerySync, true)
}

func sendAppendTx(privateKey crypto.PrivKey, tx types.SignedTx) tmsp.Result {
	return sendToTendermint(privateKey, tx, client.AppendTxSync, false)
}

func sendToTendermint(privateKey crypto.PrivKey, tx types.SignedTx, fn func(tx []byte) (res tmsp.Result), isQuery bool) tmsp.Result {
	if err := tx.SignTx(privateKey, chainID); err != nil {
		return tmsp.ErrBaseInvalidSignature.AppendLog(err.Error())
	}

	var txs []byte
	if isQuery {
		txs = []byte{tx.TxType()}
		txs = append(txs, wire.BinaryBytes(struct{ types.Tx }{tx})...)
	} else {
		txs = wire.BinaryBytes(struct{ types.Tx }{tx})
	}

	return fn(txs)
}

// StartClient is a convenience function to start the client app
func StartClient(serverAddress string) {
	var err error
	client, err = tmspcli.NewClient(serverAddress, "socket", true)
	if err != nil {
		panic("connecting to tmsp_app: " + err.Error())
	}

	log.Info("Tendermint server connection established to " + serverAddress)
}

func SetOption(key, value string) {
	res := client.SetOptionSync(key, value)
	if res.IsErr() {
		panic(fmt.Sprintf("setting %v=%v: \nlog: %v", key, value, res.Log))
	}
}

func Commit(client tmspcli.Client) {
	res := client.CommitSync()

	if res.IsErr() {
		panic(fmt.Sprintf("Commit error: %v", res.Log))
	}

	log.Info("Committed tx")
}

func AppendTx(txBytes []byte) {
	res := client.AppendTxSync(txBytes)
	if res.IsErr() {
		panic(fmt.Sprintf("AppendTx %X: %v\nlog: %v", txBytes, res, res.Log))
	}
}

func CheckTx(txBytes []byte) {
	res := client.CheckTxSync(txBytes)
	if res.IsErr() {
		panic(fmt.Sprintf("checking tx %X: %v\nlog: %v", txBytes, res, res.Log))
	}
}

func Query(txBytes []byte) (res tmsp.Result) {
	return client.QuerySync(txBytes)
}

func printKey(key []byte, title string) {
	fmt.Println(title)
	for _, v := range key {
		fmt.Print(v)
		fmt.Print(", ")
	}
	fmt.Println()

}

func GetLegalEntities(privateKey crypto.PrivKey, ids []string) (returned LegalEntitiesReturned) {
	tx := &types.LegalEntityQueryTx{Ids: ids,
		Address: privateKey.PubKey().Address()}

	res := sendQuery(privateKey, tx)
	if res.IsErr() {
		panic(fmt.Sprintf("Error in tendermint response: %v ", res.Log))
	}
	err := json.Unmarshal(res.Data, &returned)
	if err != nil {
		panic(fmt.Sprintf("JSON unmarshal for message %v failed with: %v ", res.Data, err))
	}

	return
}
