package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/tendermint/clearchain/client"
	"github.com/tendermint/clearchain/types"
	"github.com/tendermint/go-crypto"
)

var serverAddress = "tcp://127.0.0.1:46658"
var chainID = "test_chain_id"
var privateKeyInBase64 = "ATRXWwlJ6bvNRcNRT/EMmymjZvAGsLZp5a95t9HL5NRhhDh4uTLuSQikLSS//AOeuN+s1DQMgzQjEGgglAR/r6s="

var privateKey crypto.PrivKey

const (
	ServerAddressKey = "serverAddress"
	ChainIDKey       = "chainID"
	PrivateKey       = "privateKey"
	Help             = "help"
)

func main() {
	handleCommandLine()

	client.SetChainID(chainID)
	client.StartClient(serverAddress)

	startWebserver()
}

func handleCommandLine() {
	flag.String(ServerAddressKey, "", "TMSP address to Tendermint server")
	flag.String(ChainIDKey, "", "Blockchain ID")
	flag.String(PrivateKey, "", "Base64 encoded privateKey for message signing")
	flag.Bool(Help, false, "Prints command line argument usage")

	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("!!!! Running in test mode with hardwired Tendermint server connection config !!!! ")
	}
	flag.Visit(flagHandler)

	var err error
	privateKey, err = crypto.PrivKeyFromBytes(client.Decode(privateKeyInBase64))
	if err != nil {
		panic(fmt.Sprintf("Error during building privateKey from %v, %v", privateKeyInBase64, err))
	}
}

func flagHandler(currentFlag *flag.Flag) {
	fmt.Println("Setting", currentFlag.Name, "to ", currentFlag.Value)

	switch currentFlag.Name {
	case ServerAddressKey:
		serverAddress = currentFlag.Value.String()
		return
	case ChainIDKey:
		chainID = currentFlag.Value.String()
		return
	case PrivateKey:
		privateKeyInBase64 = currentFlag.Value.String()
		return
	case Help:
		flag.Usage()
		return
	default:
		panic(fmt.Sprintf(":( Unimplemented flag: %v", currentFlag.Name))
	}
}

func startWebserver() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/", defaultHandler)

	http.ListenAndServe(":8080", nil)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Unimplemented request for %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: Security issue: No autentication, authorization is there to limit access to this code.

	accountsRequested := client.GetAllAccounts(privateKey).Accounts
	accounts := client.GetAccounts(privateKey, accountsRequested).Account
	
	legalEntityIds := client.GetAllLegalEntities(privateKey).Ids
	legalEntities := client.GetLegalEntities(privateKey, legalEntityIds).LegalEntities
	
	jsonBytes, err := json.Marshal(struct {
		LegalEntities []*types.LegalEntity `json:"legalEntities"`
		Account       []*types.Account     `json:"accounts"`
	}{legalEntities, accounts})

	if err != nil {
		panic(fmt.Sprintf("JSON marshalling error for message: %v, %v", accounts, err))
	}

	w.Write(jsonBytes)
}
