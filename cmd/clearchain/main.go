package main

import (	
	"fmt"	
	"flag"		
	"os"
	
	"github.com/tendermint/clearchain/app"
		
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
	common "github.com/tendermint/tmlibs/common"
	abci "github.com/tendermint/abci/types"
)

// Entry point of the Go app

func main() {
	// Read application flags
	addrPtr := flag.String("address", "0.0.0.0:46658", "Listen address")
	genFilePath := flag.String("genesis", "", "Genesis file, if any")
	flag.Parse()

	// Create Clearchain app. It creates a /data folder 
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
	db, err := dbm.NewGoLevelDB("clearchain", "data")
	//db := dbm.NewMemDB()	
	if err != nil {
		panic(err)
	}
	app := app.NewClearchainApp(app.AppName, "cc", logger, db)		
	
	// If genesis file was specified, set key-value options	
	fmt.Println("genesis filePath: " +  *genFilePath)
	if *genFilePath != "" {
		initStateFromGenesis(app, *genFilePath)
	}	
	
	// Start the listener
	fmt.Printf("Clearchain app started. Running forever on %s \n", *addrPtr)
	app.RunForever(*addrPtr)
}

// initStateFromGenesis populates the state
func initStateFromGenesis(app *app.ClearchainApp , genFilePath string) {
	stateBytes, _ := common.ReadFile(genFilePath)
	vals := []abci.Validator{}
	res := app.InitChain(abci.RequestInitChain{vals, stateBytes})
	fmt.Printf("Result from InitChain:  %s \n", res.String())	
}

