package app

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"

	"github.com/tendermint/clearchain/types"
)

// TODO: init state
// TODO: query
func TestApp_DepositMsg(t *testing.T) {
	cc := NewClearchainApp("depositMsg", "cc")
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
	cc.EndBlock(abci.RequestEndBlock{})

	// TODO: not working yet...
	// cres := cc.Commit()
	// assert.NotEqual(t, 0, len(cres.Data))
}

func TestApp_FreezeOperator(t *testing.T) {
	cc := NewClearchainApp("freezeOperator", "cc")

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
	cc := NewClearchainApp("freezeOperator", "cc")

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
	cc.StoreAccount(acct)
	return addr
}

func fakeOpAccount(cc *ClearchainApp, typ string, entityName string) (crypto.Address, crypto.PrivKey) {
	priv := crypto.GenPrivKeyEd25519()
	pub := priv.PubKey()
	addr := pub.Address()
	acct := types.NewOpUser(pub, nil, entityName, typ)
	cc.StoreAccount(acct)
	return addr, priv
}

func fakeAdminAccount(cc *ClearchainApp, typ string, entityName string) (crypto.Address, crypto.PrivKey) {
	priv := crypto.GenPrivKeyEd25519()
	pub := priv.PubKey()
	addr := pub.Address()
	acct := types.NewAdminUser(pub, nil, entityName, typ)
	cc.StoreAccount(acct)
	return addr, priv
}
