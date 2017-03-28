package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"reflect"
	"path"
	"os"
	
	"github.com/tendermint/abci/server"
	"github.com/tendermint/clearchain/app"
	common "github.com/tendermint/go-common"
	eyes "github.com/tendermint/merkleeyes/client"
)

const EyesCacheSize = 10000

func main() {

	fmt.Println("Starting Clearchain...")
	addrPtr := flag.String("address", "tcp://0.0.0.0:46658", "Listen address")
	eyesPtr := flag.String("eyes", "local", "MerkleEyes address, or 'local' for embedded")
	genFilePath := flag.String("genesis", "", "Genesis file, if any")
	flag.Parse()

	// Connect to MerkleEyes	
	var eyesCli *eyes.Client
	if *eyesPtr == "local" {
		clearchainDir := ClearchainRoot("")
		localDBPath := path.Join(clearchainDir, "dataTmp", "merkleeyes.db")
		fmt.Println("starting local MerkleEyes. Path: " +  localDBPath)
		eyesCli = eyes.NewLocalClient(localDBPath, EyesCacheSize)
	} else {
		fmt.Println("starting remote MerkleEyes")
		var err error
		eyesCli, err = eyes.NewClient(*eyesPtr)
		if err != nil {
			common.Exit("connect to MerkleEyes: " + err.Error())
		}
	}

	// Create Clearing app
	app := app.NewLedger(eyesCli)

	// If genesis file was specified, set key-value options
	fmt.Println("genesis filePath: " +  *genFilePath)
	if *genFilePath != "" {
		kvz := loadGenesis(*genFilePath)
		for _, kv := range kvz {
			log := app.SetOption(kv.Key, kv.Value)
			fmt.Println(common.Fmt("Log: %v Set %v=%v. ", log, kv.Key, kv.Value))
		}
	}

	// Start the listener
	svr, err := server.NewServer(*addrPtr, "socket", app)
	if err != nil {
		common.Exit("create listener: " + err.Error())
	}

	// Wait forever
	common.TrapSignal(func() {
		// Cleanup
		svr.Stop()
	})

}

//----------------------------------------

// KeyValue defines the attributes of a configuration variable
type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func loadGenesis(filePath string) (kvz []KeyValue) {
	kvzFromFile := []interface{}{}
	bytes, err := common.ReadFile(filePath)
	if err != nil {
		common.Exit("loading genesis file: " + err.Error())
	}
	err = json.Unmarshal(bytes, &kvzFromFile)
	if err != nil {
		common.Exit("parsing genesis file: " + err.Error())
	}
	if len(kvzFromFile)%2 != 0 {
		common.Exit("genesis cannot have an odd number of items.  Format = [key1, value1, key2, value2, ...]")
	}
	for i := 0; i < len(kvzFromFile); i += 2 {
		keyIfc := kvzFromFile[i]
		valueIfc := kvzFromFile[i+1]
		var key, value string
		key, ok := keyIfc.(string)
		if !ok {
			common.Exit(common.Fmt(
				"genesis had invalid key %v of type %v", keyIfc, reflect.TypeOf(keyIfc)))
		}
		if v, ok := valueIfc.(string); ok {
			value = v
		} else {
			valueBytes, err := json.Marshal(valueIfc)
			if err != nil {
				common.Exit(common.Fmt(
					"genesis had invalid value %v: %v", v, err.Error()))
			}
			value = string(valueBytes)
		}
		kvz = append(kvz, KeyValue{key, value})
	}
	return kvz
}

func ClearchainRoot(rootDir string) string {
	if rootDir == "" {
		rootDir = os.Getenv("CCHOME")
	}
	if rootDir == "" {
		rootDir = os.Getenv("HOME") + "/.clearchain"
	}
	return rootDir
}
