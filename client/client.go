package client

import (
	"encoding/json"
	"fmt"

	"github.com/tendermint/clearchain/types"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
	"github.com/tendermint/tmsp/client"
	tendermintTypes "github.com/tendermint/tmsp/types"
)

// AccountsReturned defines the attributes of response's payload
type AccountsReturned struct {
	Account []*types.Account `json:"accounts"`
}

var chainID string
var client tmspcli.Client

// SetChainID assigns and initializes the chain's ID
func SetChainID(id string) {
	chainID = id
}

// GetAccounts makes a request to the ledger to returns a set of accounts
func GetAccounts(privateKey crypto.PrivKey, accountsRequested []string) AccountsReturned {
	publicKey := privateKey.PubKey()
	tx := &types.AccountQueryTx{Accounts: accountsRequested,
		Address: publicKey.Address()}
	tx.Signature = privateKey.Sign(tx.SignBytes(chainID))

	txs := []byte{types.TxTypeQueryAccount}
	txs = append(txs, wire.BinaryBytes(struct{ types.Tx }{tx})...)

	res := client.QuerySync(txs)

	var accountsReturned AccountsReturned

	err := json.Unmarshal(res.Data, &accountsReturned)
	if err != nil {
		panic(err)
	}
	return accountsReturned
}

// StartClient is a convenience function to start the client app
func StartClient(serverAddress string) {
	var err error
	client, err = tmspcli.NewClient(serverAddress, "socket", true)
	if err != nil {
		panic("connecting to tmsp_app: " + err.Error())
	}
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
		panic(fmt.Sprintf("committing: %v", res.Log))
	}
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

func Query(txBytes []byte) (res tendermintTypes.Result) {
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
