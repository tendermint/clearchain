package app

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/abci/server"

	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	abci "github.com/tendermint/abci/types"

	"github.com/tendermint/clearchain/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

const AppName = "ClearchainApp"

// ClearchainApp is basic application
type ClearchainApp struct {
	*baseapp.BaseApp
	accts sdk.AccountMapper
}

func NewClearchainApp(appname, storeKey string, logger log.Logger, db dbm.DB) *ClearchainApp {
	// var app = &ClearchainApp{}

	// make multistore with various keys
	mainKey := sdk.NewKVStoreKey(storeKey)
	// ibcKey = sdk.NewKVStoreKey("ibc")

	bApp := baseapp.NewBaseApp(appname, logger, db)
	mountMultiStore(bApp, mainKey)
	err := bApp.LoadLatestVersion(mainKey)
	if err != nil {
		panic(err)
	}

	// register routes on new application
	accts := types.AccountMapper(mainKey)
	types.RegisterRoutes(bApp.Router(), accts)

	// set up ante and tx parsing
	setAnteHandler(bApp, accts)
	initBaseAppTxDecoder(bApp)

	ccApp := &ClearchainApp{
		BaseApp: bApp,
		accts:   accts,
	}
	ccApp.SetInitChainer(ccApp.initChainer)

	return ccApp
}

// RunForever starts the abci server
func (app *ClearchainApp) RunForever(addrPtr string) {
	srv, err := server.NewServer(addrPtr, "socket", app)
	if err != nil {
		panic(err)
	}
	srv.Start()
	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})
}

func mountMultiStore(bApp *baseapp.BaseApp,
	keys ...*sdk.KVStoreKey) {

	// create substore for every key
	for _, key := range keys {
		bApp.MountStore(key, sdk.StoreTypeIAVL)
	}
}

func setAnteHandler(bApp *baseapp.BaseApp, accts sdk.AccountMapper) {
	// this checks auth, but may take fee is future, check for compatibility
	bApp.SetAnteHandler(
		auth.NewAnteHandler(accts))
}

func initBaseAppTxDecoder(bApp *baseapp.BaseApp) {
	cdc := types.MakeTxCodec()
	bApp.SetTxDecoder(func(txBytes []byte) (sdk.Tx, sdk.Error) {
		var tx = sdk.StdTx{}
		// StdTx.Msg is an interface whose concrete
		// types are registered in app/msgs.go.
		err := cdc.UnmarshalBinary(txBytes, &tx)
		if err != nil {
			return nil, sdk.ErrTxParse("").TraceCause(err, "")
		}
		return tx, nil
	})
}

// custom logic for clearchain initialization
func (app *ClearchainApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	genesisState := new(types.GenesisState)
	err := json.Unmarshal(stateJSON, genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}
	for _, gAdUsr := range genesisState.AdminUsers {
		acc, err := gAdUsr.ToAdminUser()
		if err != nil {
			panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
			//	return sdk.ErrGenesisParse("").TraceCause(err, "")
		}
		app.accts.SetAccount(ctx, acc)
		fmt.Println("***** Set Admin user *****")
		fmt.Printf("Entity name: %v \n", acc.EntityName)
		fmt.Printf("Entity type: %v \n", acc.EntityType)
		fmt.Println("*****")
	}
	fmt.Println("Genesis file loaded successfully!")
	return abci.ResponseInitChain{}
}
