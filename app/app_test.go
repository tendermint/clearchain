package app

import (
	"bytes"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"

	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	common "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"github.com/tendermint/clearchain/types"
)

// TODO: init state
func TestApp_DepositMsg(t *testing.T) {

	cc := newTestClearchainApp()
	junk := []byte("khekfhewgfsug")
	cc.BeginBlock(abci.RequestBeginBlock{})
	// send a deposit msg
	// garbage in, garbage out
	ctx := cc.NewContext(false, abci.Header{})
	chOpAddr, chOpPrivKey := fakeOpAccount(cc, ctx, types.EntityClearingHouse, "CH")
	custAssetAddr := fakeAssetAccount(cc, ctx, sdk.Coins{}, types.EntityCustodian, "CUST")
	memberAssetAddr := fakeAssetAccount(cc, ctx, sdk.Coins{}, types.EntityIndividualClearingMember, "ICM")
	depositMsg := types.DepositMsg{Operator: chOpAddr, Sender: custAssetAddr,
		Recipient: memberAssetAddr, Amount: sdk.Coin{"USD", 700}}
	depositTx := makeTx(cc.cdc, depositMsg, chOpPrivKey)
	dres := cc.DeliverTx(junk)
	assert.EqualValues(t, sdk.CodeTxDecode, dres.Code, dres.Log)
	// get real working
	dres = cc.DeliverTx(depositTx)
	assert.EqualValues(t, sdk.CodeOK, dres.Code, dres.Log)
	cc.Commit()
	cc.EndBlock(abci.RequestEndBlock{})

	// Query data to verify the deposit
	res := cc.Query(abci.RequestQuery{Data: memberAssetAddr, Path: "/main/key"})
	codec := cc.cdc
	var foundAcc *types.AppAccount
	err := codec.UnmarshalBinary(res.GetValue(), &foundAcc)
	assert.Nil(t, err)
	assert.NotNil(t, foundAcc)
	assert.Equal(t, int64(700), foundAcc.Coins.AmountOf("USD"))
}

func TestApp_FreezeOperator(t *testing.T) {
	cc := newTestClearchainApp()

	cc.BeginBlock(abci.RequestBeginBlock{})
	// send a deposit msg
	ctx := cc.NewContext(false, abci.Header{})
	chAdmAddr, chAdmPrivKey := fakeAdminAccount(cc, ctx, types.EntityClearingHouse, "CH")
	chOpAddr, chOpPrivKey := fakeOpAccount(cc, ctx, types.EntityClearingHouse, "CH")
	custAssetAddr := fakeAssetAccount(cc, ctx, sdk.Coins{}, types.EntityCustodian, "CUST")
	memberAssetAddr := fakeAssetAccount(cc, ctx, sdk.Coins{}, types.EntityIndividualClearingMember, "ICM")
	depositMsg := types.DepositMsg{Operator: chOpAddr, Sender: custAssetAddr,
		Recipient: memberAssetAddr, Amount: sdk.Coin{"USD", 700}}
	freezeOpMsg := types.FreezeOperatorMsg{types.BaseFreezeAccountMsg{Admin: chAdmAddr, Target: chOpAddr}}
	freezeOperatorTx := makeTx(cc.cdc, freezeOpMsg, chAdmPrivKey)
	dres := cc.DeliverTx(freezeOperatorTx)
	assert.EqualValues(t, sdk.CodeOK, dres.Code, dres.Log)
	depositTx := makeTx(cc.cdc, depositMsg, chOpPrivKey)
	// get real working
	dres = cc.DeliverTx(depositTx)
	assert.EqualValues(t, types.CodeInactiveAccount, dres.Code, dres.Log)
	cc.EndBlock(abci.RequestEndBlock{})
}

func TestApp_FreezeAdmin(t *testing.T) {
	cc := newTestClearchainApp()

	cc.BeginBlock(abci.RequestBeginBlock{})
	// send a deposit msg
	ctx := cc.NewContext(false, abci.Header{})
	chAdmAddr, chAdmPrivKey := fakeAdminAccount(cc, ctx, types.EntityClearingHouse, "CH")
	memAdmAddr, memAdmPrivKey := fakeAdminAccount(cc, ctx, types.EntityGeneralClearingMember, "GCM")
	memOpAddr, _ := fakeOpAccount(cc, ctx, types.EntityGeneralClearingMember, "GCM")
	// the CH admin freezes the member admin
	freezeAdmMsg := types.FreezeAdminMsg{types.BaseFreezeAccountMsg{Admin: chAdmAddr, Target: memAdmAddr}}
	freezeAdmTx := makeTx(cc.cdc, freezeAdmMsg, chAdmPrivKey)
	dres := cc.DeliverTx(freezeAdmTx)
	assert.EqualValues(t, sdk.CodeOK, dres.Code, dres.Log)
	// the member's admin can no longer freeze its own operators
	freezeOpMsg := types.FreezeOperatorMsg{types.BaseFreezeAccountMsg{Admin: memAdmAddr, Target: memOpAddr}}
	freezeOperatorTx := makeTx(cc.cdc, freezeOpMsg, memAdmPrivKey)
	dres = cc.DeliverTx(freezeOperatorTx)
	assert.EqualValues(t, types.CodeInactiveAccount, dres.Code, dres.Log)
	cc.EndBlock(abci.RequestEndBlock{})
}

//Test_Genesis is an end-to-end test that verifies the complete process of loading a genesis file.
// It makes the app read an external genesis file and then verifies that all accounts were created by using the Query interface
func Test_Genesis(t *testing.T) {

	app := newTestClearchainApp()
	codec := app.cdc
	absPathFileOk, err := filepath.Abs("test/genesis_ok_ch_admin_test.json")
	pubBytes, _ := hex.DecodeString("01328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f874")
	publicKey1, _ := crypto.PubKeyFromBytes(pubBytes)

	adminCreated1 := types.NewAdminUser(publicKey1, nil, "ClearChain", "ch")
	stateBytes, _ := common.ReadFile(absPathFileOk)
	vals := []abci.Validator{}
	res := app.Query(abci.RequestQuery{Data: []byte("nothing"), Path: "/main/key"})
	assert.Equal(t, 0, len(res.Value))
	app.BeginBlock(abci.RequestBeginBlock{})
	app.InitChain(abci.RequestInitChain{Validators: vals, AppStateBytes: stateBytes})
	app.Commit()
	app.EndBlock(abci.RequestEndBlock{})

	expAcc := adminCreated1
	// Query the existing data
	res = app.Query(abci.RequestQuery{Data: expAcc.GetAddress().Bytes(), Path: "/main/key"})
	assert.NotNil(t, res.GetValue())
	var foundAcc *types.AppAccount
	err = codec.UnmarshalBinary(res.Value, &foundAcc)
	if err != nil {
		panic(err)
	}
	assert.Nil(t, err)
	assert.True(t, bytes.Equal(expAcc.Address.Bytes(), foundAcc.Address.Bytes()))
	assert.True(t, expAcc.PubKey.Equals(foundAcc.PubKey))
	assert.True(t, foundAcc.Coins.IsZero())
	assert.Equal(t, len(foundAcc.Creator), 0)
	assert.Equal(t, expAcc.EntityName, foundAcc.EntityName)
	assert.Equal(t, expAcc.EntityType, foundAcc.EntityType)
	assert.Equal(t, expAcc.AccountType, foundAcc.AccountType)
	assert.Equal(t, types.AccountUser, foundAcc.AccountType)
	assert.Equal(t, expAcc.Active, foundAcc.Active)
	assert.True(t, foundAcc.Active)
	assert.Equal(t, expAcc.Admin, foundAcc.Admin)
	assert.True(t, foundAcc.Admin)
}

func makeTx(cdc *wire.Codec, msg sdk.Msg, keys ...crypto.PrivKey) []byte {
	sigs := make([]sdk.StdSignature, len(keys))
	for i, k := range keys {
		sig := k.Sign(sdk.StdSignBytes("", []int64{0}, sdk.StdFee{}, msg))
		sigs[i] = sdk.StdSignature{
			PubKey:    k.PubKey(),
			Signature: sig,
			Sequence:  0,
		}
	}
	tx := sdk.NewStdTx(msg, sdk.StdFee{}, sigs)

	bz, err := cdc.MarshalBinary(tx)
	if err != nil {
		panic(err)
	}

	return bz
}

func fakeAssetAccount(cc *ClearchainApp, ctx sdk.Context, cash sdk.Coins, typ string, entityName string) sdk.Address {
	pub := crypto.GenPrivKeyEd25519().PubKey()
	addr := pub.Address()
	acct := types.NewAssetAccount(pub, cash, nil, entityName, typ)
	cc.accountMapper.SetAccount(ctx, acct)
	return addr
}

func fakeOpAccount(cc *ClearchainApp, ctx sdk.Context, typ string, entityName string) (sdk.Address, crypto.PrivKey) {
	priv := crypto.GenPrivKeyEd25519()
	pub := priv.PubKey()
	addr := pub.Address()
	acct := types.NewOpUser(pub, []byte(""), entityName, typ)
	cc.accountMapper.SetAccount(ctx, acct)
	return addr, priv.Wrap()
}

func fakeAdminAccount(cc *ClearchainApp, ctx sdk.Context, typ string, entityName string) (sdk.Address, crypto.PrivKey) {
	priv := crypto.GenPrivKeyEd25519()
	pub := priv.PubKey()
	addr := pub.Address()
	acct := types.NewAdminUser(pub, []byte(""), entityName, typ)
	cc.accountMapper.SetAccount(ctx, acct)
	return addr, priv.Wrap()
}

// newTestClearchainApp a ClearchainApp with an in-memory datastore
func newTestClearchainApp() *ClearchainApp {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "app")
	db := dbm.NewMemDB()
	return NewClearchainApp(logger, db)
}
