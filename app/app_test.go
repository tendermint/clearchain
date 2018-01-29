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

	cust, cKey := fakeAccount(cc, types.EntityCustodian, nil)
	member, _ := fakeAccount(cc, types.EntityIndividualClearingMember, nil)
	msg := types.DepositMsg{Sender: cust, Recipient: member, Amount: sdk.Coin{"USD", 700}}
	real := makeTx(msg, cKey)

	cc.BeginBlock(abci.RequestBeginBlock{})
	// garbage in, garbage out
	dres := cc.DeliverTx(junk)
	assert.EqualValues(t, sdk.CodeTxParse, dres.Code, dres.Log)

	// get real working
	dres = cc.DeliverTx(real)
	assert.EqualValues(t, sdk.CodeOK, dres.Code, dres.Log)

	cc.EndBlock(abci.RequestEndBlock{})
	// no data in db
	// cres := cc.Commit()
	// assert.Equal(t, 0, len(cres.Data))
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

	cc := makeTxCodec()
	bz, err := cc.MarshalBinary(tx)
	if err != nil {
		panic(err)
	}

	return bz
}

// func fakeAccount(accts sdk.AccountMapper, ctx sdk.Context, typ string, cash sdk.Coins) crypto.Address {
func fakeAccount(cc *ClearchainApp, typ string, cash sdk.Coins) (crypto.Address, crypto.PrivKey) {
	priv := crypto.GenPrivKeyEd25519()
	pub := priv.PubKey()
	addr := pub.Address()

	acct := new(types.AppAccount)
	acct.SetAddress(addr)
	acct.SetPubKey(pub)
	acct.SetCoins(cash)
	acct.Type = typ

	cc.StoreAccount(acct)

	return addr, priv
}
