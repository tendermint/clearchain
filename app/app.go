package app

import (
	"encoding/json"
	"strings"

	bctypes "github.com/tendermint/basecoin/types"
	sm "github.com/tendermint/clearchain/state"
	"github.com/tendermint/clearchain/types"
	common "github.com/tendermint/go-common"
	"github.com/tendermint/go-wire"
	"github.com/tendermint/governmint/gov"
	eyes "github.com/tendermint/merkleeyes/client"
	tmsp "github.com/tendermint/tmsp/types"
)

const (
	version   = "0.0.1"
	maxTxSize = 10240

	// PluginTypeByteBase defines the base plugin's byte code
	PluginTypeByteBase = 0x01
	// PluginTypeByteEyes defines the eyes plugin's byte code
	PluginTypeByteEyes = 0x02
	// PluginTypeByteGov defines the gov plugin's byte code
	PluginTypeByteGov = 0x03

	// PluginNameBase defines the base plugin's name
	PluginNameBase = "base"
	// PluginNameEyes defines the eyes plugin's name
	PluginNameEyes = "eyes"
	// PluginNameGov defines the gov plugin's name
	PluginNameGov = "gov"
)

// Ledger defines the attributes of the app
type Ledger struct {
	eyesCli    *eyes.Client
	govMint    *gov.Governmint
	state      *sm.State
	cacheState *sm.State
	plugins    *bctypes.Plugins
}

// NewLedger creates a new instance of the app
func NewLedger(eyesCli *eyes.Client) *Ledger {
	govMint := gov.NewGovernmint()
	state := sm.NewState(eyesCli)
	plugins := bctypes.NewPlugins()
	plugins.RegisterPlugin(PluginTypeByteGov, PluginNameGov, govMint)
	return &Ledger{
		eyesCli:    eyesCli,
		govMint:    govMint,
		state:      state,
		cacheState: nil,
		plugins:    plugins,
	}
}

// Info returns app's generic information
func (app *Ledger) Info() string {
	return common.Fmt("Ledger v%v", version)
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
		accountIndex := sm.GetOrMakeAccountIndex(app.state)
		accountIndex.Add(acc.ID)
		app.state.SetAccountIndex(accountIndex)

		return "Success"
	case "user":
		var err error
		var user *types.User
		wire.ReadJSONPtr(&user, []byte(value), &err)
		if err != nil {
			panic("Error decoding user message: " + err.Error())
		}

		app.state.SetUser(user.PubKey.Address(), user)
		return "Success"
	case "legalEntity":
		var legalEntity types.LegalEntity

		err := json.Unmarshal([]byte(value), &legalEntity)

		if err != nil {
			panic("Error decoding legalEntity message: " + err.Error())
		}

		app.state.SetLegalEntity(legalEntity.ID, &legalEntity)
		legalEntities := app.state.GetLegalEntityIndex()
		if legalEntities == nil {
			legalEntities = &types.LegalEntityIndex{Ids: []string{}}
		}
		legalEntities.Add(legalEntity.ID)
		app.state.SetLegalEntityIndex(legalEntities)
		return "Success"
	}
	return "Unrecognized option key " + key
}

// AppendTx handles appendTx
func (app *Ledger) AppendTx(txBytes []byte) (res tmsp.Result) {
	return app.executeTx(txBytes, false)
}

// CheckTx handles checkTx
func (app *Ledger) CheckTx(txBytes []byte) (res tmsp.Result) {
	return app.executeTx(txBytes, true)
}

// Query handles queryTx
func (app *Ledger) Query(query []byte) (res tmsp.Result) {
	if len(query) == 0 {
		return tmsp.ErrEncodingError.SetLog("Query cannot be zero length")
	}
	typeByte := query[0]
	query = query[1:]
	switch typeByte {
	case types.TxTypeQueryAccount, types.TxTypeQueryAccountIndex, types.TxTypeLegalEntity, types.TxTypeQueryLegalEntityIndex:
		return app.executeQueryTx(query)
	case PluginTypeByteBase:
		return tmsp.OK.SetLog("This type of query not yet supported")
	case PluginTypeByteEyes:
		return app.eyesCli.QuerySync(query)
	}
	return tmsp.ErrBaseUnknownPlugin.SetLog(
		common.Fmt("Unknown plugin with type byte %X", typeByte))
}

// Commit handles commitTx
func (app *Ledger) Commit() (res tmsp.Result) {
	// Commit eyes.
	res = app.eyesCli.CommitSync()
	if res.IsErr() {
		common.PanicSanity("Error getting hash: " + res.Error())
	}
	return res
}

// InitChain initializes the chain
func (app *Ledger) InitChain(validators []*tmsp.Validator) {
	for _, plugin := range app.plugins.GetList() {
		plugin.Plugin.InitChain(app.state, validators)
	}
}

// TMSP::BeginBlock
func (app *Ledger) BeginBlock(height uint64) {
	for _, plugin := range app.plugins.GetList() {
		plugin.Plugin.BeginBlock(app.state, height)
	}
	app.cacheState = app.state.CacheWrap()
}

// TMSP::EndBlock
func (app *Ledger) EndBlock(height uint64) (diffs []*tmsp.Validator) {
	for _, plugin := range app.plugins.GetList() {
		moreDiffs := plugin.Plugin.EndBlock(app.state, height)
		diffs = append(diffs, moreDiffs...)
	}
	return
}

func (app *Ledger) executeTx(txBytes []byte, simulate bool) (res tmsp.Result) {
	if len(txBytes) > maxTxSize {
		return tmsp.ErrBaseEncodingError.AppendLog("Tx size exceeds maximum")
	}
	// Decode tx
	var tx types.Tx
	err := wire.ReadBinaryBytes(txBytes, &tx)
	if err != nil {
		return tmsp.ErrBaseEncodingError.AppendLog("Error decoding tx: " + err.Error())
	}
	// Validate and exec tx
	res = sm.ExecTx(app.state, app.plugins, tx, simulate, nil)
	if res.IsErr() {
		return res.PrependLog("Error in AppendTx")
	}
	return res
}

func (app *Ledger) executeQueryTx(txBytes []byte) (res tmsp.Result) {
	if len(txBytes) > maxTxSize {
		return tmsp.ErrBaseEncodingError.AppendLog("Tx size exceeds maximum")
	}
	// Decode tx
	var tx types.Tx
	err := wire.ReadBinaryBytes(txBytes, &tx)
	if err != nil {
		return tmsp.ErrBaseEncodingError.AppendLog("Error decoding tx: " + err.Error())
	}
	// Validate and exec tx
	res = sm.ExecQueryTx(app.state, tx)
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
