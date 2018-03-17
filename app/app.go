package app

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/abci/server"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/clearchain/types"
	wire "github.com/tendermint/go-wire"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

// AppName defines the name of the app.
const AppName = "ClearchainApp"

// ClearchainApp is basic application
type ClearchainApp struct {
	*baseapp.BaseApp
	cdc             *wire.Codec
	capKeyMainStore *sdk.KVStoreKey
	capKeyIBCStore  *sdk.KVStoreKey
	accountMapper   sdk.AccountMapper
}

// NewClearchainApp creates a new ClearchainApp type.
func NewClearchainApp(logger log.Logger, db dbm.DB) *ClearchainApp {
	var app = &ClearchainApp{
		BaseApp:         baseapp.NewBaseApp(AppName, logger, db),
		cdc:             types.MakeCodec(),
		capKeyMainStore: sdk.NewKVStoreKey("main"),
		capKeyIBCStore:  sdk.NewKVStoreKey("ibc"),
	}
	// define the account mapper
	app.accountMapper = auth.NewAccountMapperSealed(
		app.capKeyMainStore, // target store
		&types.AppAccount{}, // prototype
	)
	// add handlers and register routes
	types.RegisterRoutes(app.Router(), app.accountMapper)

	// initialise BaseApp
	app.SetTxDecoder(app.txDecoder)
	app.SetInitChainer(app.initChainer)
	// Multi-store feature is currently broken
	// https://github.com/cosmos/cosmos-sdk/issues/532
	app.MountStoresIAVL(app.capKeyMainStore)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper))
	err := app.LoadLatestVersion(app.capKeyMainStore)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
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

// custom logic for transaction decoding
func (app *ClearchainApp) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	// StdTx.Msg is an interface. The concrete types are registered by MakeCodec.
	var tx = sdk.StdTx{}
	err := app.cdc.UnmarshalBinary(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxParse("").TraceCause(err, "")
	}
	return tx, nil
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

	genChAdmin := genesisState.ClearingHouseAdmin
	if (genChAdmin != types.GenesisAccount{}) {
		acc, err := genChAdmin.ToClearingHouseAdmin()
		if err != nil {
			panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
			//	return sdk.ErrGenesisParse("").TraceCause(err, "")
		}
		app.accountMapper.SetAccount(ctx, acc)
		fmt.Println("***** Set Ch Admin *****")
		fmt.Printf("Entity name: %v \n", acc.EntityName)
		fmt.Printf("Entity type: %v \n", acc.EntityType)
		fmt.Printf("Public key: %v \n", genChAdmin.PubKeyHexa)
		fmt.Println("*****")
	}

	fmt.Println("Genesis file loaded successfully!")
	return abci.ResponseInitChain{}
}
