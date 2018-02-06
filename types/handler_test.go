package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
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

func TestWithdrawMsgHandler(t *testing.T) {
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
			"good Withdraw",
			args{ctx: ctx, msg: WithdrawMsg{Sender: member, Recipient: cust, Operator: operator, Amount: sdk.Coin{"USD", 500}}},
			sdk.CodeOK,
			sdk.Coins{{"USD", 500}},
			sdk.Coins{{"EUR", 5000}, {"USD", 500}},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := WithdrawMsgHandler(accts)
			got := handler(tt.args.ctx, tt.args.msg)
			assert.Equal(t, tt.expect, got.Code, got.Log)

			c := accts.GetAccount(ctx, cust)
			assert.Equal(t, tt.cBal, c.GetCoins())

			m := accts.GetAccount(ctx, member)
			assert.Equal(t, tt.mBal, m.GetCoins())
		})
	}
}

func TestCreateAccountMsgHandler(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	createAccount := func(typ, entityName string, isAdm bool) crypto.Address {
		if isAdm {
			return fakeAdminAccountWithEntityName(accts, ctx, typ, sdk.Coins{}, entityName)
		}
		return fakeAccountWithEntityName(accts, ctx, typ, sdk.Coins{}, entityName)
	}
	mkKey := func() crypto.PubKey {
		return crypto.GenPrivKeyEd25519().PubKey()
	}
	chOp := createAccount(EntityClearingHouse, "CH", false)
	chAdm := createAccount(EntityClearingHouse, "CH", true)
	custOp := createAccount(EntityCustodian, "CUST", false)
	custAdm := createAccount(EntityCustodian, "CUST", true)
	existingCustOp := makeAccountWithEntityName(EntityCustodian, nil, false, "CUST")
	accts.SetAccount(ctx, existingCustOp)

	type args struct {
		ctx sdk.Context
		msg sdk.Msg
	}
	tests := []struct {
		name   string
		args   args
		expect sdk.CodeType
	}{
		{
			"CH admin can create CH admin", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: chAdm, PubKey: mkKey(), AccountType: EntityClearingHouse,
				IsAdmin: true, LegalEntityName: "CH"}}, sdk.CodeOK,
		},
		{
			"CH admin can create CH op", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: chAdm, PubKey: mkKey(), AccountType: EntityClearingHouse,
				IsAdmin: false, LegalEntityName: "CH"}}, sdk.CodeOK,
		},
		{
			"CH op cannot create CH admin", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: chOp, PubKey: mkKey(), AccountType: EntityClearingHouse,
				IsAdmin: true, LegalEntityName: "CH"}}, sdk.CodeUnauthorized,
		},
		{
			"CH admin can create CUST admin", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: chAdm, PubKey: mkKey(), AccountType: EntityCustodian,
				IsAdmin: true, LegalEntityName: "CUST"}}, sdk.CodeOK,
		},
		{
			"CH admin can create CUST op", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: chAdm, PubKey: mkKey(), AccountType: EntityCustodian,
				IsAdmin: false, LegalEntityName: "CUST"}}, sdk.CodeUnauthorized,
		},
		{
			"CH op cannot create CUST admin", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: chOp, PubKey: mkKey(), AccountType: EntityCustodian,
				IsAdmin: true, LegalEntityName: "CUST"}}, sdk.CodeUnauthorized,
		},
		{
			"CH op cannot create CUST op", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: chOp, PubKey: mkKey(), AccountType: EntityCustodian,
				IsAdmin: false, LegalEntityName: "CUST"}}, sdk.CodeUnauthorized,
		},
		{
			"CUST admin can create CUST op", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: custAdm, PubKey: mkKey(), AccountType: EntityCustodian,
				IsAdmin: false, LegalEntityName: "CUST"}}, sdk.CodeOK,
		},
		{
			"CUST op cannot create CUST op", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: custOp, PubKey: mkKey(), AccountType: EntityCustodian,
				IsAdmin: false, LegalEntityName: "CUST"}}, sdk.CodeUnauthorized,
		},
		{
			"CUST admin cannot create member op", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: custAdm, PubKey: mkKey(), AccountType: EntityIndividualClearingMember,
				IsAdmin: false, LegalEntityName: "ICM"}}, sdk.CodeUnauthorized,
		},
		{
			"CUST admin cannot create CH op", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: custAdm, PubKey: mkKey(), AccountType: EntityClearingHouse,
				IsAdmin: false, LegalEntityName: "CH"}}, sdk.CodeUnauthorized,
		},
		{
			"CUST admin cannot create CH admin", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: custAdm, PubKey: mkKey(), AccountType: EntityClearingHouse,
				IsAdmin: true, LegalEntityName: "CH"}}, sdk.CodeUnauthorized,
		},
		{
			"CUST admin cannot create CH admin", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: mkKey().Address(), PubKey: mkKey(), AccountType: EntityClearingHouse,
				IsAdmin: true, LegalEntityName: "CH"}}, CodeInvalidAccount,
		},
		{
			"account already exists", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: custAdm, PubKey: existingCustOp.PubKey, AccountType: EntityCustodian,
				IsAdmin: true, LegalEntityName: "CUST"}}, CodeInvalidAccount,
		},
		{
			"nil creator", args{ctx: ctx, msg: CreateAccountMsg{
				Creator: nil, PubKey: mkKey(), AccountType: EntityClearingHouse,
				IsAdmin: true, LegalEntityName: "CH"}}, CodeInvalidAccount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := CreateAccountMsgHandler(accts)
			got := handler(tt.args.ctx, tt.args.msg)
			assert.Equal(t, tt.expect, got.Code, got.Log)

			newAcc := accts.GetAccount(ctx, tt.args.msg.(CreateAccountMsg).PubKey.Address())
			if tt.expect == sdk.CodeOK || tt.name == "account already exists" {
				assert.True(t, newAcc != nil)
			} else {
				assert.True(t, newAcc == nil)
			}

		})
	}
}

//---------------- helpers --------------------

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
	acct := makeAccount(typ, cash, false)
	accts.SetAccount(ctx, acct)
	return acct.Address
}

func fakeAccountWithEntityName(accts sdk.AccountMapper, ctx sdk.Context, typ string, cash sdk.Coins, entityName string) crypto.Address {
	acct := makeAccountWithEntityName(typ, cash, false, entityName)
	accts.SetAccount(ctx, acct)
	return acct.Address
}

func fakeAdminAccount(accts sdk.AccountMapper, ctx sdk.Context, typ string, cash sdk.Coins) crypto.Address {
	acct := makeAccount(typ, cash, true)
	accts.SetAccount(ctx, acct)
	return acct.Address
}

func fakeAdminAccountWithEntityName(accts sdk.AccountMapper, ctx sdk.Context, typ string, cash sdk.Coins, entityName string) crypto.Address {
	acct := makeAccountWithEntityName(typ, cash, true, entityName)
	accts.SetAccount(ctx, acct)
	return acct.Address
}

func makeAccount(typ string, cash sdk.Coins, isAdmin bool) *AppAccount {
	pub := crypto.GenPrivKeyEd25519().PubKey()
	addr := pub.Address()
	return NewAppAccount(pub, cash, typ, nil, isAdmin, string(addr))
}

func makeAccountWithEntityName(typ string, cash sdk.Coins, isAdmin bool, entityName string) *AppAccount {
	acct := makeAccount(typ, cash, isAdmin)
	acct.LegalEntityName = entityName
	return acct
}
