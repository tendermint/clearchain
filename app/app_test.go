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
func TestApp(t *testing.T) {
	cc := NewClearchainApp()
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
