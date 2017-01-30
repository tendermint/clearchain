package client

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
//	abcicli "github.com/tendermint/abci/client"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/clearchain/types"
	"github.com/tendermint/go-common"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-logger"
	"github.com/tendermint/go-rpc/client"
	"github.com/tendermint/go-rpc/types"
	"github.com/tendermint/go-wire"
	_ "github.com/tendermint/tendermint/rpc/core/types" // Register RPCResponse > Result types
)

var log = logger.New("module", "client")
var chainID string

//var client abcicli.Client

var ws *rpcclient.WSClient

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

	res := sendDeliverTxSync(privateKey, tx)

	if res.IsErr() {
		panic(fmt.Sprintf("Wrong response from server: %v", res))
	} else {
		Commit()
	}
}

func CreateAccount(privateKey crypto.PrivKey,
	accountID string) {
	tx := &types.CreateAccountTx{Address: privateKey.PubKey().Address(),
		AccountID: accountID}

	res := sendDeliverTxSync(privateKey, tx)

	if res.IsErr() {
		panic(fmt.Sprintf("Wrong response from server: %v", res))
	} else {
		Commit()
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

	res := sendDeliverTxSync(privateKey, tx)

	if res.IsErr() {
		panic(fmt.Sprintf("Wrong response from server: %v", res))
	} else {
		Commit()
	}
	log.Info("Created legal entity with ID: " + entityID)
}

//Creates money transfer entry in blockchain
func TransferMoney(privateKey crypto.PrivKey, senderID string, recipientID string, counterSignerAddresses [][]byte, amount int64, currency string) {
	senderAccount := GetAccounts(privateKey, []string{senderID}).Account[0]
	newSequenceID := senderAccount.GetWallet(currency).Sequence + 1

	sender := types.TxTransferSender{AccountID: senderID,
		Amount:   amount,
		Currency: currency,
		Sequence: newSequenceID,
		Address:  privateKey.PubKey().Address(),
	}

	recipient := types.TxTransferRecipient{AccountID: recipientID}

	counterSigners := make([]types.TxTransferCounterSigner, len(counterSignerAddresses))
	for i, address := range counterSignerAddresses {
		privKey, err := crypto.PrivKeyFromBytes(address)
		if err != nil {
			panic(fmt.Sprintf("counterSigner signing failed with: %v", err.Error()))
		}
		counterSigners[i] = types.TxTransferCounterSigner{Address: privKey.PubKey().Address()}

		err = counterSigners[i].SignTx(privKey, chainID)
		if err != nil {
			panic(fmt.Sprintf("counterSigner signing failed with: %v", err.Error()))
		}

	}

	tx := &types.TransferTx{
		Sender:         sender,
		Recipient:      recipient,
		CounterSigners: counterSigners,
	}

	res := sendDeliverTxSync(privateKey, tx)

	if res.IsErr() {
		panic(fmt.Sprintf("Wrong response from server: %v", res))
	} else {
		Commit()
	}
	log.Info("Created transfer entry")
}

// GetAccounts makes a request to the ledger to returns a set of accounts
func GetAccounts(privateKey crypto.PrivKey, accountsRequested []string) (returned types.AccountsReturned) {
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

func sendQuery(privateKey crypto.PrivKey, tx types.SignedTx) abci.Result {
	//	return sendToTendermint(privateKey, tx, client.QuerySync, true)

	txBytes, result := getTXBytes(privateKey, tx, true)
	if result.IsErr() {
		return result
	}

	request := rpctypes.NewRPCRequest("RPCQueryRequestID", "query", common.Arr(txBytes))
	reqBytes := wire.JSONBytes(request)
	err := ws.WriteMessage(websocket.TextMessage, reqBytes)
	if err != nil {
		panic("writing websocket request: " + err.Error())
	}

	res, ok := <-ws.ResultsCh
	if !ok {
		panic("websocket read error: " + err.Error())
	}
	
	returned := abci.OK
	returned.SetData(res)
	
	return returned
}

func sendDeliverTxSync(privateKey crypto.PrivKey, tx types.SignedTx) abci.Result {
	txBytes, result := getTXBytes(privateKey, tx, false)
	if result.IsErr() {
		return result
	}

	request := rpctypes.NewRPCRequest("RPCRequestID", "broadcast_tx_sync", common.Arr(txBytes))
	reqBytes := wire.JSONBytes(request)
	err := ws.WriteMessage(websocket.TextMessage, reqBytes)
	if err != nil {
		panic("writing websocket request: " + err.Error())
	}

	return abci.OK
}

func sendToTendermint(privateKey crypto.PrivKey, tx types.SignedTx, fn func(tx []byte) (res abci.Result), isQuery bool) abci.Result {
	txBytes, result := getTXBytes(privateKey, tx, isQuery)

	if result.IsErr() {
		return result
	}

	return fn(txBytes)
}

func getTXBytes(privateKey crypto.PrivKey, tx types.SignedTx, isQuery bool) (txs []byte, result abci.Result) {
	if err := tx.SignTx(privateKey, chainID); err != nil {
		return nil, abci.ErrBaseInvalidSignature.AppendLog(err.Error())
	}

	if isQuery {
		txs = []byte{tx.TxType()}
		txs = append(txs, wire.BinaryBytes(struct{ types.Tx }{tx})...)
	} else {
		txs = wire.BinaryBytes(struct{ types.Tx }{tx})
	}

	return txs, abci.OK
}

// StartClient is a convenience function to start the client app
func StartClient(serverAddress string) {
	//serverAddress := "192.168.99.100:46657"
	ws = rpcclient.NewWSClient(serverAddress, "/websocket")
	_, err := ws.Start()
	if err != nil {
		panic(err.Error())
	}

//	var err error
//	client, err = abcicli.NewClient(serverAddress, "socket", true)
//	if err != nil {
//		panic("connecting to abci_app: " + err.Error())
//	}

	log.Info("Tendermint server connection established to " + serverAddress)
}

func Commit() {
//	res := client.CommitSync()
//
//	if res.IsErr() {
//		panic(fmt.Sprintf("Commit error: %v", res.Log))
//	}
//
//	log.Info("Committed tx")
}

func GetLegalEntities(privateKey crypto.PrivKey, ids []string) (returned types.LegalEntitiesReturned) {
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
