package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TestRegisterRoutes is an end-to-end test, making sure a normal workflow is
// supported and passing all messages through the router to simulate production code path
func TestRegisterRoutes(t *testing.T) {
	accts, ctx := fakeAccountMapper()

	cust := fakeAccount(accts, ctx, EntityCustodian, nil)
	member := fakeAccount(accts, ctx, EntityIndividualClearingMember, nil)
	member2 := fakeAccount(accts, ctx, EntityGeneralClearingMember, nil)
	operator := fakeAccount(accts, ctx, EntityClearingHouse, nil)

	router := baseapp.NewRouter()
	RegisterRoutes(router, accts)

	type args struct {
		ctx sdk.Context
		msg sdk.Msg
	}
	tests := []struct {
		name   string
		args   args
		expect sdk.CodeType
		cBal   sdk.Coins
		mBal   sdk.Coins
		mBal2  sdk.Coins
		chBal  sdk.Coins
	}{
		{
			"good deposit",
			args{ctx: ctx, msg: DepositMsg{Sender: cust, Recipient: member, Amount: sdk.Coin{"USD", 5000}}},
			sdk.CodeOK,
			sdk.Coins{{"USD", -5000}},
			sdk.Coins{{"USD", 5000}},
			nil,
			nil,
		},
		{
			"deposit2",
			args{ctx: ctx, msg: DepositMsg{Sender: cust, Recipient: member2, Amount: sdk.Coin{"USD", 7777}}},
			sdk.CodeOK,
			sdk.Coins{{"USD", -12777}},
			sdk.Coins{{"USD", 5000}},
			sdk.Coins{{"USD", 7777}},
			nil,
		},
		{
			"settlement",
			args{ctx: ctx, msg: SettleMsg{Sender: operator, Recipient: member, Amount: sdk.Coin{"USD", 3000}}},
			sdk.CodeOK,
			sdk.Coins{{"USD", -12777}},
			sdk.Coins{{"USD", 8000}},
			sdk.Coins{{"USD", 7777}},
			sdk.Coins{{"USD", -3000}},
		},
		{
			"counter settlement",
			args{ctx: ctx, msg: SettleMsg{Sender: operator, Recipient: member2, Amount: sdk.Coin{"USD", -3000}}},
			sdk.CodeOK,
			sdk.Coins{{"USD", -12777}},
			sdk.Coins{{"USD", 8000}},
			sdk.Coins{{"USD", 4777}},
			nil,
		},
		{
			"withdraw",
			args{ctx: ctx, msg: WithdrawMsg{Sender: member, Recipient: cust, Operator: operator, Amount: sdk.Coin{"USD", 5500}}},
			sdk.CodeOK,
			sdk.Coins{{"USD", -7277}},
			sdk.Coins{{"USD", 2500}},
			sdk.Coins{{"USD", 4777}},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := router.Route(tt.args.msg.Type())
			got := h(tt.args.ctx, tt.args.msg)
			assert.Equal(t, tt.expect, got.Code, got.Log)

			c := accts.GetAccount(ctx, cust)
			assert.Equal(t, tt.cBal, c.GetCoins())

			m := accts.GetAccount(ctx, member)
			assert.Equal(t, tt.mBal, m.GetCoins())

			ch := accts.GetAccount(ctx, operator)
			assert.Equal(t, tt.chBal, ch.GetCoins())
		})
	}
}

func fakeAccountMapper() (sdk.AccountMapper, sdk.Context) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	key := sdk.NewKVStoreKey("test")
	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	accts := AccountMapper(key)

	h := abci.Header{
		Height:  100,
		ChainID: "clear-chain",
	}
	ctx := sdk.NewContext(ms, h, false, []byte{1, 2, 3, 4}) // DeliverTx

	return accts, ctx
}

func fakeAccount(accts sdk.AccountMapper, ctx sdk.Context, typ string, cash sdk.Coins) crypto.Address {
	pub := crypto.GenPrivKeyEd25519().PubKey()
	addr := pub.Address()

	acct := new(AppAccount)
	acct.SetAddress(addr)
	acct.SetPubKey(pub)
	acct.SetCoins(cash)
	acct.Type = typ

	accts.SetAccount(ctx, acct)
	return addr
}

func TestDepositMsgHandler(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	cCoins := sdk.Coins{{"EUR", 5000}, {"USD", 1000}}
	mCoins := sdk.Coins{}

	cust := fakeAccount(accts, ctx, EntityCustodian, cCoins)
	member := fakeAccount(accts, ctx, EntityIndividualClearingMember, mCoins)

	type args struct {
		ctx sdk.Context
		msg sdk.Msg
	}
	tests := []struct {
		name   string
		args   args
		expect sdk.CodeType
		cBal   sdk.Coins
		mBal   sdk.Coins
	}{
		{
			"no returns",
			args{ctx: ctx, msg: DepositMsg{Sender: member, Recipient: cust, Amount: sdk.Coin{"USD", 200}}},
			CodeWrongSigner,
			cCoins,
			nil, // sdk.Coins{}
		},
		{
			"no negative deposits",
			args{ctx: ctx, msg: DepositMsg{Sender: member, Recipient: cust, Amount: sdk.Coin{"USD", -500}}},
			CodeInvalidAmount,
			cCoins,
			nil, // sdk.Coins{}
		},
		{
			"no 0 deposits",
			args{ctx: ctx, msg: DepositMsg{Sender: member, Recipient: cust, Amount: sdk.Coin{"USD", 0}}},
			CodeInvalidAmount,
			cCoins,
			nil, // sdk.Coins{}
		},
		{
			"good deposit",
			args{ctx: ctx, msg: DepositMsg{Sender: cust, Recipient: member, Amount: sdk.Coin{"USD", 700}}},
			sdk.CodeOK,
			sdk.Coins{{"EUR", 5000}, {"USD", 300}},
			sdk.Coins{{"USD", 700}},
		},
		{
			// allow the custodian to go negative
			"overdraft",
			args{ctx: ctx, msg: DepositMsg{Sender: cust, Recipient: member, Amount: sdk.Coin{"EUR", 10000}}},
			sdk.CodeOK,
			sdk.Coins{{"EUR", -5000}, {"USD", 300}},
			sdk.Coins{{"EUR", 10000}, {"USD", 700}},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := DepositMsgHandler(accts)
			got := handler(tt.args.ctx, tt.args.msg)
			assert.Equal(t, tt.expect, got.Code, got.Log)

			c := accts.GetAccount(ctx, cust)
			assert.Equal(t, tt.cBal, c.GetCoins())

			m := accts.GetAccount(ctx, member)
			assert.Equal(t, tt.mBal, m.GetCoins())
		})
	}
}

func TestSettleMsgHandler(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	clhCoins := sdk.Coins{{"EUR", 5000}, {"USD", 1000}}
	mCoins := sdk.Coins{{"USD", 1000}}

	clh := fakeAccount(accts, ctx, EntityClearingHouse, clhCoins)
	member := fakeAccount(accts, ctx, EntityIndividualClearingMember, mCoins)

	type args struct {
		ctx sdk.Context
		msg sdk.Msg
	}
	tests := []struct {
		name   string
		args   args
		expect sdk.CodeType
		cBal   sdk.Coins
		mBal   sdk.Coins
	}{
		{
			"no returns",
			args{ctx: ctx, msg: SettleMsg{Sender: member, Recipient: clh, Amount: sdk.Coin{"USD", 200}}},
			CodeWrongSigner,
			clhCoins,
			mCoins,
		},
		{
			"no 0 settle",
			args{ctx: ctx, msg: SettleMsg{Sender: clh, Recipient: member, Amount: sdk.Coin{"USD", 0}}},
			CodeInvalidAmount,
			clhCoins,
			mCoins,
		},
		{
			"negative good settle",
			args{ctx: ctx, msg: SettleMsg{Sender: clh, Recipient: member, Amount: sdk.Coin{"USD", -500}}},
			sdk.CodeOK,
			sdk.Coins{{"EUR", 5000}, {"USD", 1500}},
			sdk.Coins{{"USD", 500}},
		},

		{
			"positive good settle",
			args{ctx: ctx, msg: SettleMsg{Sender: clh, Recipient: member, Amount: sdk.Coin{"USD", 500}}},
			sdk.CodeOK,
			sdk.Coins{{"EUR", 5000}, {"USD", 1000}},
			sdk.Coins{{"USD", 1000}},
		},
		{
			// allow the clearing house to go negative
			"overdraft",
			args{ctx: ctx, msg: SettleMsg{Sender: clh, Recipient: member, Amount: sdk.Coin{"EUR", 10000}}},
			sdk.CodeOK,
			sdk.Coins{{"EUR", -5000}, {"USD", 1000}},
			sdk.Coins{{"EUR", 10000}, {"USD", 1000}},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := SettleMsgHandler(accts)
			got := handler(tt.args.ctx, tt.args.msg)
			assert.Equal(t, tt.expect, got.Code, got.Log)

			c := accts.GetAccount(ctx, clh)
			assert.Equal(t, tt.cBal, c.GetCoins())

			m := accts.GetAccount(ctx, member)
			assert.Equal(t, tt.mBal, m.GetCoins())
		})
	}
}

func TestWithDrawMsgHandler(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	mCoins := sdk.Coins{{"EUR", 5000}, {"USD", 1000}}
	custCoins := sdk.Coins{}

	cust := fakeAccount(accts, ctx, EntityCustodian, custCoins)
	member := fakeAccount(accts, ctx, EntityIndividualClearingMember, mCoins)
	operator := fakeAccount(accts, ctx, EntityClearingHouse, nil)

	type args struct {
		ctx sdk.Context
		msg sdk.Msg
	}
	tests := []struct {
		name   string
		args   args
		expect sdk.CodeType
		cBal   sdk.Coins
		mBal   sdk.Coins
	}{
		{
			"no returns",
			args{ctx: ctx, msg: WithdrawMsg{Sender: cust, Recipient: member, Operator: operator, Amount: sdk.Coin{"USD", 200}}},
			CodeWrongSigner,
			nil,
			mCoins,
		},
		{
			"no 0 Withdraw",
			args{ctx: ctx, msg: WithdrawMsg{Sender: member, Recipient: cust, Operator: operator, Amount: sdk.Coin{"USD", 0}}},
			CodeInvalidAmount,
			nil,
			mCoins,
		},
		{
			"no negative Withdraw",
			args{ctx: ctx, msg: WithdrawMsg{Sender: member, Recipient: cust, Operator: operator, Amount: sdk.Coin{"USD", -500}}},
			CodeInvalidAmount,
			nil,
			mCoins,
		},

		{
			"good Withdraw",
			args{ctx: ctx, msg: WithdrawMsg{Sender: member, Recipient: cust, Operator: operator, Amount: sdk.Coin{"USD", 500}}},
			sdk.CodeOK,
			sdk.Coins{{"USD", 500}},
			sdk.Coins{{"EUR", 5000}, {"USD", 500}},
		},
		{
			"overdraft",
			args{ctx: ctx, msg: WithdrawMsg{Sender: member, Recipient: cust, Operator: operator, Amount: sdk.Coin{"EUR", 10000}}},
			sdk.CodeInsufficientFunds,
			sdk.Coins{{"USD", 500}},
			sdk.Coins{{"EUR", 5000}, {"USD", 500}},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := WithDrawMsgHandler(accts)
			got := handler(tt.args.ctx, tt.args.msg)
			assert.Equal(t, tt.expect, got.Code, got.Log)

			c := accts.GetAccount(ctx, cust)
			assert.Equal(t, tt.cBal, c.GetCoins())

			m := accts.GetAccount(ctx, member)
			assert.Equal(t, tt.mBal, m.GetCoins())
		})
	}
}