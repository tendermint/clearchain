package app

import (
	"encoding/json"
	"strings"

	abci "github.com/tendermint/abci/types"
	bctypes "github.com/tendermint/basecoin/types"
	"github.com/tendermint/clearchain/state"
	"github.com/tendermint/clearchain/types"
	common "github.com/tendermint/go-common"
	"github.com/tendermint/go-wire"
	eyes "github.com/tendermint/merkleeyes/client"
)

const (
	version   = "0.0.1"
	maxTxSize = 10240

	// PluginTypeByteBase defines the base plugin's byte code
	PluginTypeByteBase = 0x01
	// PluginTypeByteEyes defines the eyes plugin's byte code
	PluginTypeByteEyes = 0x02

	// PluginNameBase defines the base plugin's name
	PluginNameBase = "base"
	// PluginNameEyes defines the eyes plugin's name
	PluginNameEyes = "eyes"
)

// Ledger defines the attributes of the app
type Ledger struct {
	eyesCli    *eyes.Client
	state      *state.State
	cacheState *state.State
	plugins    *bctypes.Plugins
}

// NewLedger creates a new instance of the app
func NewLedger(eyesCli *eyes.Client) *Ledger {
	state := state.NewState(eyesCli)
	plugins := bctypes.NewPlugins()
	return &Ledger{
		eyesCli:    eyesCli,
		state:      state,
		cacheState: nil,
		plugins:    plugins,
	}
}

// Info returns app's generic information
func (app *Ledger) Info() abci.ResponseInfo {
	return abci.ResponseInfo{Data: common.Fmt("Ledger v%v", version)}
}

func (app *Ledger) RegisterPlugin(plugin bctypes.Plugin) {
	app.plugins.RegisterPlugin(plugin)
}

// SetOption modifies app's configuration
func (app *Ledger) SetOption(key string, value string) (log string) {
	PluginName, key := splitKey(key)
	if PluginName != PluginNameBase {
		// Set option on plugin
		plugin := app.plugins.GetByName(PluginName)
		if plugin == nil {
			panic("Invalid plugin name: " + PluginName)
		}
		return plugin.SetOption(app.state, key, value)
	}
	// Set option on Clearing
	switch key {
	case "chainID":
		app.state.SetChainID(value)
		return "Success"
	case "account":
		var err error
		var acc *types.Account
		wire.ReadJSONPtr(&acc, []byte(value), &err)
		if err != nil {
			panic("Error decoding acc message: " + err.Error())
		}
		app.state.SetAccount(acc.ID, acc)
		state.SetAccountInIndex(app.state, *acc)
		app.Commit()
		return "Success"
	case "user":
		var err error
		var user *types.User
		wire.ReadJSONPtr(&user, []byte(value), &err)
		if err != nil {
			panic("Error decoding user message: " + err.Error())
		}

		app.state.SetUser(user.PubKey.Address(), user)
		app.Commit()
		return "Success"
	case "legalEntity":
		var legalEntity types.LegalEntity

		err := json.Unmarshal([]byte(value), &legalEntity)

		if err != nil {
			panic("Error decoding legalEntity message: " + err.Error())
		}

		app.state.SetLegalEntity(legalEntity.ID, &legalEntity)
		state.SetLegalEntityInIndex(app.state, &legalEntity)
		app.Commit()

		return "Success"
	}
	return "Unrecognized option key " + key
}

// DeliverTx handles deliverTx
func (app *Ledger) DeliverTx(txBytes []byte) (res abci.Result) {
	return app.executeTx(txBytes, false)
}

// CheckTx handles checkTx
func (app *Ledger) CheckTx(txBytes []byte) (res abci.Result) {
	return app.executeTx(txBytes, true)
}

// Query handles queryTx
func (app *Ledger) Query(query []byte) (res abci.Result) {
	if len(query) == 0 {
		return abci.ErrEncodingError.SetLog("Query cannot be zero length")
	}
	typeByte := query[0]
	query = query[1:]
	switch typeByte {
	case types.TxTypeQueryAccount, types.TxTypeQueryAccountIndex, types.TxTypeLegalEntity, types.TxTypeQueryLegalEntityIndex:
		return app.executeQueryTx(query)
	case PluginTypeByteBase:
		return abci.OK.SetLog("This type of query not yet supported")
	case PluginTypeByteEyes:
		return app.eyesCli.QuerySync(query)
	}
	return abci.ErrBaseUnknownPlugin.SetLog(
		common.Fmt("Unknown plugin with type byte %X", typeByte))
}

// Commit handles commitTx
func (app *Ledger) Commit() (res abci.Result) {
	// Commit eyes.
	res = app.eyesCli.CommitSync()
	if res.IsErr() {
		common.PanicSanity("Error getting hash: " + res.Error())
	}
	return res
}

// InitChain initializes the chain
func (app *Ledger) InitChain(validators []*abci.Validator) {
	for _, plugin := range app.plugins.GetList() {
		plugin.InitChain(app.state, validators)
	}
}

// abci::BeginBlock
func (app *Ledger) BeginBlock(hash []byte, header *abci.Header) {
	for _, plugin := range app.plugins.GetList() {
		plugin.BeginBlock(app.state, hash, header)
	}
	app.cacheState = app.state.CacheWrap()
}

// abci::EndBlock
func (app *Ledger) EndBlock(height uint64) (res abci.ResponseEndBlock) {
	for _, plugin := range app.plugins.GetList() {
		pluginRes  := plugin.EndBlock(app.state, height)
		res.Diffs = append(res.Diffs, pluginRes.Diffs...)
	}
	return
}

func (app *Ledger) executeTx(txBytes []byte, simulate bool) (res abci.Result) {
	if len(txBytes) > maxTxSize {
		return abci.ErrBaseEncodingError.AppendLog("Tx size exceeds maximum")
	}
	// Decode tx
	var tx types.Tx
	err := wire.ReadBinaryBytes(txBytes, &tx)
	if err != nil {
		return abci.ErrBaseEncodingError.AppendLog("Error decoding tx: " + err.Error())
	}
	// Validate and exec tx
	res = state.ExecTx(app.state, app.plugins, tx, simulate, nil)
	if res.IsErr() {
		return res.PrependLog("Error in DeliverTx")
	}
	return res
}

func (app *Ledger) executeQueryTx(txBytes []byte) (res abci.Result) {
	if len(txBytes) > maxTxSize {
		return abci.ErrBaseEncodingError.AppendLog("Tx size exceeds maximum")
	}
	// Decode tx
	var tx types.Tx
	err := wire.ReadBinaryBytes(txBytes, &tx)
	if err != nil {
		return abci.ErrBaseEncodingError.AppendLog("Error decoding tx: " + err.Error())
	}
	// Validate and exec tx
	res = state.ExecQueryTx(app.state, tx)
	if res.IsErr() {
		return res.PrependLog("Error in QueryTx")
	}
	return res

}

// Splits the string at the first '/'.
// if there are none, the second string is nil.
func splitKey(key string) (prefix string, suffix string) {
	if strings.Contains(key, "/") {
		keyParts := strings.SplitN(key, "/", 2)
		return keyParts[0], keyParts[1]
	}
	return key, ""
}
