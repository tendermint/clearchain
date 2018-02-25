package app

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	common "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"github.com/tendermint/clearchain/types"
)

// TODO: init state
// TODO: query
func TestApp_DepositMsg(t *testing.T) {

	cc := newTestClearchainApp("depositMsg", "cc")
	junk := []byte("khekfhewgfsug")

	cc.BeginBlock(abci.RequestBeginBlock{})
	// send a deposit msg
	chOpAddr, chOpPrivKey := fakeOpAccount(cc, types.EntityClearingHouse, "CH")
	custAssetAddr := fakeAssetAccount(cc, nil, types.EntityCustodian, "CUST")
	memberAssetAddr := fakeAssetAccount(cc, nil, types.EntityIndividualClearingMember, "ICM")
	depositMsg := types.DepositMsg{Operator: chOpAddr, Sender: custAssetAddr,
		Recipient: memberAssetAddr, Amount: sdk.Coin{"USD", 700}}
	depositTx := makeTx(depositMsg, chOpPrivKey)
	// garbage in, garbage out
	dres := cc.DeliverTx(junk)
	assert.EqualValues(t, sdk.CodeTxParse, dres.Code, dres.Log)
	// get real working
	dres = cc.DeliverTx(depositTx)
	assert.EqualValues(t, sdk.CodeOK, dres.Code, dres.Log)
	cc.Commit()
	cc.EndBlock(abci.RequestEndBlock{})

	// Query data to verify the deposit
	res := cc.Query(abci.RequestQuery{Data: memberAssetAddr, Path: "/cc/key"})
	codec := types.MakeTxCodec()
	var foundAcc types.AppAccount
	err := codec.UnmarshalBinary(res.GetValue(), &foundAcc)
	assert.Nil(t, err)
	assert.NotNil(t, foundAcc)
	assert.Equal(t, int64(700), foundAcc.Coins.AmountOf("USD"))
}

func TestApp_FreezeOperator(t *testing.T) {
	cc := newTestClearchainApp("freezeOperator", "cc")

	cc.BeginBlock(abci.RequestBeginBlock{})
	// send a deposit msg
	chAdmAddr, chAdmPrivKey := fakeAdminAccount(cc, types.EntityClearingHouse, "CH")
	chOpAddr, chOpPrivKey := fakeOpAccount(cc, types.EntityClearingHouse, "CH")
	custAssetAddr := fakeAssetAccount(cc, nil, types.EntityCustodian, "CUST")
	memberAssetAddr := fakeAssetAccount(cc, nil, types.EntityIndividualClearingMember, "ICM")
	depositMsg := types.DepositMsg{Operator: chOpAddr, Sender: custAssetAddr,
		Recipient: memberAssetAddr, Amount: sdk.Coin{"USD", 700}}
	freezeOpMsg := types.FreezeOperatorMsg{types.BaseFreezeAccountMsg{Admin: chAdmAddr, Target: chOpAddr}}
	freezeOperatorTx := makeTx(freezeOpMsg, chAdmPrivKey)
	dres := cc.DeliverTx(freezeOperatorTx)
	assert.EqualValues(t, sdk.CodeOK, dres.Code, dres.Log)
	depositTx := makeTx(depositMsg, chOpPrivKey)
	// get real working
	dres = cc.DeliverTx(depositTx)
	assert.EqualValues(t, types.CodeInactiveAccount, dres.Code, dres.Log)
	cc.EndBlock(abci.RequestEndBlock{})
}

func TestApp_FreezeAdmin(t *testing.T) {
	cc := newTestClearchainApp("freezeOperator", "cc")

	cc.BeginBlock(abci.RequestBeginBlock{})
	// send a deposit msg
	chAdmAddr, chAdmPrivKey := fakeAdminAccount(cc, types.EntityClearingHouse, "CH")
	memAdmAddr, memAdmPrivKey := fakeAdminAccount(cc, types.EntityGeneralClearingMember, "GCM")
	memOpAddr, _ := fakeOpAccount(cc, types.EntityGeneralClearingMember, "GCM")
	// the CH admin freezes the member admin
	freezeAdmMsg := types.FreezeAdminMsg{types.BaseFreezeAccountMsg{Admin: chAdmAddr, Target: memAdmAddr}}
	freezeAdmTx := makeTx(freezeAdmMsg, chAdmPrivKey)
	dres := cc.DeliverTx(freezeAdmTx)
	assert.EqualValues(t, sdk.CodeOK, dres.Code, dres.Log)
	// the member's admin can no longer freeze its own operators
	freezeOpMsg := types.FreezeOperatorMsg{types.BaseFreezeAccountMsg{Admin: memAdmAddr, Target: memOpAddr}}
	freezeOperatorTx := makeTx(freezeOpMsg, memAdmPrivKey)
	dres = cc.DeliverTx(freezeOperatorTx)
	assert.EqualValues(t, types.CodeInactiveAccount, dres.Code, dres.Log)
	cc.EndBlock(abci.RequestEndBlock{})
}

//Test_Genesis is an end-to-end test that verifies the complete process of loading a genesis file.
// It makes the app read an external genesis file and then verifies that all accounts were created by using the Query interface
func Test_Genesis(t *testing.T) {

	codec := types.MakeTxCodec()
	app := newTestClearchainApp("loadFromGenesis", "cc")
	absPathFileOk, _ := filepath.Abs("test/genesis_ok_test.1.json")
	pubBytes, _ := hex.DecodeString("328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f874")
	publicKey1, _ := crypto.PubKeyFromBytes(pubBytes)
	pubBytes, _ = hex.DecodeString("328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f875")
	publicKey2, _ := crypto.PubKeyFromBytes(pubBytes)
	pubBytes, _ = hex.DecodeString("328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f876")
	publicKey3, _ := crypto.PubKeyFromBytes(pubBytes)

	adminCreated1 := types.NewAdminUser(publicKey1, nil, "ClearChain", "ch")
	adminCreated2 := types.NewAdminUser(publicKey2, nil, "ClearingHouse", "ch")
	adminCreated3 := types.NewAdminUser(publicKey3, nil, "Admin", "gcm")

	stateBytes, _ := common.ReadFile(absPathFileOk)
	vals := []abci.Validator{}
	app.BeginBlock(abci.RequestBeginBlock{})
	app.InitChain(abci.RequestInitChain{vals, stateBytes})
	app.Commit()
	app.EndBlock(abci.RequestEndBlock{})

	expectedAccounts := []*types.AppAccount{adminCreated1, adminCreated2, adminCreated3}
	for _, expAcc := range expectedAccounts {
		// Query the existing data
		res := app.Query(abci.RequestQuery{Data: expAcc.GetAddress(), Path: "/cc/key"})
		assert.NotNil(t, res.GetValue())
		var foundAcc types.AppAccount
		err := codec.UnmarshalBinary(res.GetValue(), &foundAcc)
		assert.Nil(t, err)
		assert.Equal(t, hex.EncodeToString(expAcc.Address), hex.EncodeToString(foundAcc.Address))
		assert.True(t, expAcc.PubKey.Equals(foundAcc.PubKey))
		assert.True(t, foundAcc.Coins.IsZero())
		assert.True(t, foundAcc.Creator == nil)
		assert.Equal(t, expAcc.EntityName, foundAcc.EntityName)
		assert.Equal(t, expAcc.EntityType, foundAcc.EntityType)
		assert.Equal(t, expAcc.AccountType, foundAcc.AccountType)
		assert.Equal(t, types.AccountUser, foundAcc.AccountType)
		assert.Equal(t, expAcc.Active, foundAcc.Active)
		assert.True(t, foundAcc.Active)
		assert.Equal(t, expAcc.Admin, foundAcc.Admin)
		assert.True(t, foundAcc.Admin)
	}
}

func makeTx(msg sdk.Msg, keys ...crypto.PrivKey) []byte {
	tx := sdk.StdTx{Msg: msg}

	sigs := make([]sdk.StdSignature, len(keys))
	for i, k := range keys {
		sigs[i] = sdk.StdSignature{
			PubKey:    k.PubKey(),
			Sequence:  0,
			Signature: k.Sign(tx.GetSignBytes()),
		}
	}
	tx.Signatures = sigs

	cc := types.MakeTxCodec()
	bz, err := cc.MarshalBinary(tx)
	if err != nil {
		panic(err)
	}

	return bz
}

func fakeAssetAccount(cc *ClearchainApp, cash sdk.Coins, typ string, entityName string) crypto.Address {
	pub := crypto.GenPrivKeyEd25519().PubKey()
	addr := pub.Address()
	acct := types.NewAssetAccount(pub, cash, nil, entityName, typ)
	var ctx = cc.NewContext(false, abci.Header{})
	cc.accts.SetAccount(ctx, acct)
	return addr
}

func fakeOpAccount(cc *ClearchainApp, typ string, entityName string) (crypto.Address, crypto.PrivKey) {
	priv := crypto.GenPrivKeyEd25519()
	pub := priv.PubKey()
	addr := pub.Address()
	acct := types.NewOpUser(pub, nil, entityName, typ)
	var ctx = cc.NewContext(false, abci.Header{})
	cc.accts.SetAccount(ctx, acct)
	return addr, priv
}

func fakeAdminAccount(cc *ClearchainApp, typ string, entityName string) (crypto.Address, crypto.PrivKey) {
	priv := crypto.GenPrivKeyEd25519()
	pub := priv.PubKey()
	addr := pub.Address()
	acct := types.NewAdminUser(pub, nil, entityName, typ)
	var ctx = cc.NewContext(false, abci.Header{})
	cc.accts.SetAccount(ctx, acct)
	return addr, priv
}

// newTestClearchainApp a ClearchainApp with an in-memory datastore
func newTestClearchainApp(appname, storeKey string) *ClearchainApp {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "app")
	db := dbm.NewMemDB()
	return NewClearchainApp(appname, storeKey, logger, db)
}
