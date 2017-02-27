package main

import (
	"encoding/json"
	"fmt"

	"github.com/tendermint/clearchain/types"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
	"github.com/tendermint/light-client/rpc"
	//	_ "github.com/tendermint/tendermint/rpc/core/types" // Register RPCResponse > Result types
)

func main() {
	privateKeyBytes := []byte{1, 52, 87, 91, 9, 73, 233, 187, 205, 69, 195, 81, 79, 241, 12, 155, 41, 163, 102, 240, 6, 176, 182, 105, 229, 175, 121, 183, 209, 203, 228, 212, 97, 132, 56, 120, 185, 50, 238, 73, 8, 164, 45, 36, 191, 252, 3, 158, 184, 223, 172, 212, 52, 12, 131, 52, 35, 16, 104, 32, 148, 4, 127, 175, 171}
	//ATRXWwlJ6bvNRcNRT/EMmymjZvAGsLZp5a95t9HL5NRhhDh4uTLuSQikLSS//AOeuN+s1DQMgzQjEGgglAR/r6s=

	var accountsRequested []string = []string{"1d2df1ae-accb-11e6-bbbb-00ff5244ae7f"}
	privateKey, _ := crypto.PrivKeyFromBytes(privateKeyBytes) //User{b40cbf4e-5923-4ccd-beec-e22a9117b91b "Name" 31}
	tx := &types.AccountQueryTx{Accounts: accountsRequested,
		Address: privateKey.PubKey().Address()}

	if err := tx.SignTx(privateKey, "test_chain_id"); err != nil {
		panic(err.Error())
	}

	txs := []byte{tx.TxType()}
	txs = append(txs, wire.BinaryBytes(struct{ types.Tx }{tx})...)

	httpClient := rpc.New("127.0.0.1:46657", "")
	result, err := httpClient.ABCIQuery(txs)
	if err != nil {
		panic(err.Error())
	}

	var returned *types.AccountsReturned
	err = json.Unmarshal(result.Result.Data, &returned)
	if err != nil {
		panic(fmt.Sprintf("JSON unmarshal for message %v failed with: %v ", returned, err))
	}

	fmt.Println(returned)
}
