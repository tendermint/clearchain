package main

import (
	"encoding/json"
	"fmt"

	"github.com/tendermint/clearchain/types"
	"github.com/tendermint/light-client/rpc"
	//	_ "github.com/tendermint/tendermint/rpc/core/types" // Register RPCResponse > Result types
)

// Tendermint core  must be running
func main() {
	
	var accountRequested string = "1d2df1ae-accb-11e6-bbbb-00ff5244ae7f"
	httpClient := rpc.NewClient("127.0.0.1:46657", "")
	var path = "/account/" + accountRequested
	result, err := httpClient.ABCIQuery(path, []byte(""), false)
	if err != nil {
		panic(err.Error())
	}

	var returned *types.AccountsReturned
	err = json.Unmarshal(result.Response.Value, &returned)
	if err != nil {
		panic(fmt.Sprintf("JSON unmarshal for message %v failed with: %v ", returned, err))
	}

	fmt.Println(returned)
}
