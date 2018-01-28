package app

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/abci/server"
	"github.com/tendermint/clearchain/types"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
	cmn "github.com/tendermint/tmlibs/common"
)

const appName = "ClearchainApp"

type ClearchainApp struct {
	*bam.BaseApp
	router          bam.Router
	cdc             *wire.Codec
	multiStore      sdk.CommitMultiStore
	capKeyMainStore *sdk.KVStoreKey
	capKeyIBCStore  *sdk.KVStoreKey
	accountMapper   sdk.AccountMapper
}

func NewClearchainApp() *ClearchainApp {
	var app = &ClearchainApp{}
	app.initCapKeys()  // ./init_capkeys.go
	app.initBaseApp()  // ./init_baseapp.go
	app.initStores()   // ./init_stores.go
	app.initHandlers() // ./init_handlers.go
	return app
}

func (app *ClearchainApp) RunForever() {
	srv, err := server.NewServer("0.0.0.0:46658", "socket", app)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	srv.Start()
	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})
}

func (app *ClearchainApp) loadStores() {
	if err := app.LoadLatestVersion(app.capKeyMainStore); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (app *ClearchainApp) initCapKeys() {
	app.capKeyMainStore = sdk.NewKVStoreKey("main")
	app.capKeyIBCStore = sdk.NewKVStoreKey("ibc")

}

func (app *ClearchainApp) initBaseApp() {
	bapp := baseapp.NewBaseApp(appName)
	app.BaseApp = bapp
	app.router = bapp.Router()
	app.initBaseAppTxDecoder()
}

func (app *ClearchainApp) initBaseAppTxDecoder() {
	cdc := makeTxCodec()
	app.BaseApp.SetTxDecoder(func(txBytes []byte) (sdk.Tx, sdk.Error) {
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

// initCapKeys, initBaseApp, initStores, initHandlers.
func (app *ClearchainApp) initStores() {
	app.mountStores()
	app.initAccountMapper()
}

// Initialize root stores.
func (app *ClearchainApp) mountStores() {

	// Create MultiStore mounts.
	app.BaseApp.MountStore(app.capKeyMainStore, sdk.StoreTypeIAVL)
	app.BaseApp.MountStore(app.capKeyIBCStore, sdk.StoreTypeIAVL)
}

// Initialize the AccountMapper.
func (app *ClearchainApp) initAccountMapper() {

	var accountMapper = auth.NewAccountMapper(
		app.capKeyMainStore, // target store
		&types.AppAccount{}, // prototype
	)

	// Register all interfaces and concrete types that
	// implement those interfaces, here.
	cdc := accountMapper.WireCodec()
	auth.RegisterWireBaseAccount(cdc)

	// Make accountMapper's WireCodec() inaccessible.
	app.accountMapper = accountMapper.Seal()
}

func (app *ClearchainApp) initHandlers() {
	app.initDefaultAnteHandler()
	app.initRouterHandlers()
}

func (app *ClearchainApp) initDefaultAnteHandler() {

	// Deducts fee from payer.
	// Verifies signatures and nonces.
	// Sets Signers to ctx.
	app.BaseApp.SetDefaultAnteHandler(
		auth.NewAnteHandler(app.accountMapper))
}

func (app *ClearchainApp) initRouterHandlers() {
	// All handlers must be added here.
	// The order matters.
	app.router.AddRoute("deposit", types.DepositMsgHandler(app.accountMapper))
	app.router.AddRoute("settle", types.SettleMsgHandler(app.accountMapper))
	app.router.AddRoute("settle", types.WithDrawMsgHandler(app.accountMapper))
}

func makeTxCodec() (cdc *wire.Codec) {
	cdc = wire.NewCodec()

	// Register crypto.[PubKey,PrivKey,Signature] types.
	crypto.RegisterWire(cdc)

	// Register clearchain types.
	types.RegisterWire(cdc)

	return
}
