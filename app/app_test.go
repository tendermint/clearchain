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
	cust, cKey := fakeAccount(cc, types.EntityCustodian, nil)
	member, _ := fakeAccount(cc, types.EntityIndividualClearingMember, nil)
	depositMsg := types.DepositMsg{Sender: cust, Recipient: member, Amount: sdk.Coin{"USD", 700}}
	depositTx := makeTx(depositMsg, cKey)
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

	//send a create account msg	
	cc.BeginBlock(abci.RequestBeginBlock{})
	clearingHouse, chKey := fakeAccount(cc, types.EntityClearingHouse, nil)
	createAccMsg := types.CreateAccountMsg{
		Creator:     clearingHouse,
		PubKey:      crypto.GenPrivKeyEd25519().PubKey(),
		AccountType: types.EntityCustodian}
		
	 createAccTx := makeTx(createAccMsg, chKey)
	res := cc.DeliverTx(createAccTx)
	assert.EqualValues(t, sdk.CodeOK, res.Code, res.Log)
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
