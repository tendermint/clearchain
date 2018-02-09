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

	_, chOpPriv := fakeUser(accts, ctx, EntityClearingHouse)
	op := chOpPriv.PubKey().Address()
	_, cust := fakeAsset(accts, ctx, nil, EntityCustodian)
	_, member := fakeAsset(accts, ctx, nil, EntityIndividualClearingMember)
	_, ch := fakeAsset(accts, ctx, nil, EntityClearingHouse)
	_, member2 := fakeAsset(accts, ctx, nil, EntityGeneralClearingMember)

	router := baseapp.NewRouter()
	RegisterRoutes(router, accts)

	type args struct {
		ctx sdk.Context
		msg sdk.Msg
	}
	tests := []struct {
		name       string
		args       args
		expect     sdk.CodeType
		custBal    sdk.Coins
		memberBal  sdk.Coins
		member2Bal sdk.Coins
		chBal      sdk.Coins
	}{
		{
			"good deposit", args{ctx: ctx, msg: DepositMsg{Operator: op, Sender: cust, Recipient: member, Amount: sdk.Coin{"USD", 5000}}},
			sdk.CodeOK, sdk.Coins{{"USD", -5000}}, sdk.Coins{{"USD", 5000}}, nil, nil,
		},
		{
			"deposit2",
			args{ctx: ctx, msg: DepositMsg{Operator: op, Sender: cust, Recipient: member2, Amount: sdk.Coin{"USD", 7777}}},
			sdk.CodeOK,
			sdk.Coins{{"USD", -12777}},
			sdk.Coins{{"USD", 5000}},
			sdk.Coins{{"USD", 7777}},
			nil,
		},
		{
			"settlement",
			args{ctx: ctx, msg: SettleMsg{Operator: op, Sender: ch, Recipient: member, Amount: sdk.Coin{"USD", 3000}}},
			sdk.CodeOK,
			sdk.Coins{{"USD", -12777}},
			sdk.Coins{{"USD", 8000}},
			sdk.Coins{{"USD", 7777}},
			sdk.Coins{{"USD", -3000}},
		},
		{
			"counter settlement",
			args{ctx: ctx, msg: SettleMsg{Operator: op, Sender: ch, Recipient: member2, Amount: sdk.Coin{"USD", -3000}}},
			sdk.CodeOK,
			sdk.Coins{{"USD", -12777}},
			sdk.Coins{{"USD", 8000}},
			sdk.Coins{{"USD", 4777}},
			nil,
		},
		{
			"withdraw",
			args{ctx: ctx, msg: WithdrawMsg{Operator: op, Sender: member, Recipient: cust, Amount: sdk.Coin{"USD", 5500}}},
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
			assert.Equal(t, tt.custBal, c.GetCoins())

			m := accts.GetAccount(ctx, member)
			assert.Equal(t, tt.memberBal, m.GetCoins())

			m2 := accts.GetAccount(ctx, member2)
			assert.Equal(t, tt.member2Bal, m2.GetCoins())

			chAsset := accts.GetAccount(ctx, ch)
			assert.Equal(t, tt.chBal, chAsset.GetCoins())
		})
	}
}
func Test_depositMsgHandler_Do(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	cCoins := sdk.Coins{{"EUR", 5000}, {"USD", 1000}}
	mCoins := sdk.Coins{}

	_, chAdm := fakeAdmin(accts, ctx, EntityClearingHouse)
	chAdmAddr := chAdm.PubKey().Address()
	_, chOpPriv := fakeUser(accts, ctx, EntityClearingHouse)
	_, cust := fakeAsset(accts, ctx, cCoins, EntityCustodian)
	_, member := fakeAsset(accts, ctx, mCoins, EntityIndividualClearingMember)
	chOp := chOpPriv.PubKey().Address()
	inactiveOp, _ := fakeInactiveUser(accts, ctx, EntityClearingHouse)

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
			"admin ain't allowed", args{ctx: ctx, msg: DepositMsg{Operator: chAdmAddr, Sender: cust,
				Recipient: member, Amount: sdk.Coin{"USD", 700}}}, CodeWrongSigner, cCoins, nil,
		},
		{
			"inactive operator", args{ctx: ctx, msg: DepositMsg{Operator: inactiveOp.Address, Sender: cust,
				Recipient: member, Amount: sdk.Coin{"USD", 700}}}, CodeWrongSigner, cCoins, nil,
		},
		{
			"no returns", args{ctx: ctx, msg: DepositMsg{Operator: chOp, Sender: member,
				Recipient: cust, Amount: sdk.Coin{"USD", 200}}}, CodeWrongSigner, cCoins, nil, // sdk.Coins{}
		},
		{
			"good deposit", args{ctx: ctx, msg: DepositMsg{Operator: chOp, Sender: cust, Recipient: member,
				Amount: sdk.Coin{"USD", 700}}}, sdk.CodeOK, sdk.Coins{{"EUR", 5000}, {"USD", 300}}, sdk.Coins{{"USD", 700}},
		},
		{
			// allow the custodian to go negative
			"overdraft", args{ctx: ctx, msg: DepositMsg{Operator: chOp, Sender: cust, Recipient: member, Amount: sdk.Coin{"EUR", 10000}}},
			sdk.CodeOK, sdk.Coins{{"EUR", -5000}, {"USD", 300}}, sdk.Coins{{"EUR", 10000}, {"USD", 700}},
		},
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

func Test_settleMsgHandler_Do(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	clhCoins := sdk.Coins{{"EUR", 5000}, {"USD", 1000}}
	mCoins := sdk.Coins{{"USD", 1000}}

	chAdmAcc, _ := fakeAdmin(accts, ctx, EntityClearingHouse)
	chOpAcc, _ := fakeUser(accts, ctx, EntityClearingHouse)
	chAdm := chAdmAcc.Address
	chOp := chOpAcc.Address
	_, clh := fakeAssetWithEntityName(accts, ctx, clhCoins, chOpAcc.GetLegalEntityType(), EntityClearingHouse)
	_, member := fakeAssetWithEntityName(accts, ctx, mCoins, "ICM", EntityIndividualClearingMember)

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
			"admins cannot settle",
			args{ctx: ctx, msg: SettleMsg{Operator: chAdm, Sender: member, Recipient: clh, Amount: sdk.Coin{"USD", 200}}},
			CodeWrongSigner, clhCoins, mCoins,
		},
		{
			"no returns",
			args{ctx: ctx, msg: SettleMsg{Operator: chOp, Sender: member, Recipient: clh, Amount: sdk.Coin{"USD", 200}}},
			CodeWrongSigner, clhCoins, mCoins,
		},
		{
			"negative good settle",
			args{ctx: ctx, msg: SettleMsg{Operator: chOp, Sender: clh, Recipient: member, Amount: sdk.Coin{"USD", -500}}},
			sdk.CodeOK, sdk.Coins{{"EUR", 5000}, {"USD", 1500}}, sdk.Coins{{"USD", 500}},
		},

		{
			"positive good settle",
			args{ctx: ctx, msg: SettleMsg{Operator: chOp, Sender: clh, Recipient: member, Amount: sdk.Coin{"USD", 500}}},
			sdk.CodeOK, sdk.Coins{{"EUR", 5000}, {"USD", 1000}}, sdk.Coins{{"USD", 1000}},
		},
		{
			// allow the clearing house to go negative
			"overdraft",
			args{ctx: ctx, msg: SettleMsg{Operator: chOp, Sender: clh, Recipient: member, Amount: sdk.Coin{"EUR", 10000}}},
			sdk.CodeOK, sdk.Coins{{"EUR", -5000}, {"USD", 1000}}, sdk.Coins{{"EUR", 10000}, {"USD", 1000}},
		},
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

func Test_withdrawMsgHandler_Do(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	mCoins := sdk.Coins{{"EUR", 5000}, {"USD", 1000}}
	custCoins := sdk.Coins{}

	_, cust := fakeAsset(accts, ctx, custCoins, EntityCustodian)
	_, member := fakeAsset(accts, ctx, mCoins, EntityIndividualClearingMember)
	chAdm, _ := fakeAdmin(accts, ctx, EntityClearingHouse)
	opAcc, _ := fakeUser(accts, ctx, EntityClearingHouse)
	admPub := chAdm.Address
	operator := opAcc.Address
	tests := []struct {
		name   string
		msg    sdk.Msg
		expect sdk.CodeType
		cBal   sdk.Coins
		mBal   sdk.Coins
	}{
		{
			"no returns", WithdrawMsg{Sender: cust, Recipient: member, Operator: operator, Amount: sdk.Coin{"USD", 200}},
			CodeWrongSigner, nil, mCoins,
		},
		{
			"good Withdraw", WithdrawMsg{Sender: member, Recipient: cust, Operator: operator, Amount: sdk.Coin{"USD", 500}},
			sdk.CodeOK, sdk.Coins{{"USD", 500}}, sdk.Coins{{"EUR", 5000}, {"USD", 500}},
		},
		{
			"admins can't withdraw", WithdrawMsg{Sender: member, Recipient: cust, Operator: admPub, Amount: sdk.Coin{"USD", 500}},
			CodeWrongSigner, sdk.Coins{{"USD", 500}}, sdk.Coins{{"EUR", 5000}, {"USD", 500}},
		},
		{
			"invalid address", WithdrawMsg{Sender: member, Recipient: cust, Operator: crypto.GenPrivKeyEd25519().PubKey().Address(), Amount: sdk.Coin{"USD", 500}},
			CodeInvalidAccount, sdk.Coins{{"USD", 500}}, sdk.Coins{{"EUR", 5000}, {"USD", 500}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := WithdrawMsgHandler(accts)
			got := handler(ctx, tt.msg)
			assert.Equal(t, tt.expect, got.Code, got.Log)

			c := accts.GetAccount(ctx, cust)
			assert.Equal(t, tt.cBal, c.GetCoins())

			m := accts.GetAccount(ctx, member)
			assert.Equal(t, tt.mBal, m.GetCoins())
		})
	}
}

func Test_createUserAccountMsgHandler_Do(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	createAccount := func(typ, entityName string, isAdm bool) crypto.Address {
		if isAdm {
			ac, _ := fakeAdminWithEntityName(accts, ctx, entityName, typ)
			return ac.PubKey.Address()
		}
		ac, _ := fakeUserWithEntityName(accts, ctx, entityName, typ)
		return ac.PubKey.Address()
	}
	mkKey := func() crypto.PubKey {
		return crypto.GenPrivKeyEd25519().PubKey()
	}
	chOp := createAccount(EntityClearingHouse, "CH", false)
	chAdm := createAccount(EntityClearingHouse, "CH", true)
	custOp := createAccount(EntityCustodian, "CUST", false)
	custAdm := createAccount(EntityCustodian, "CUST", true)
	existingCustOp, _ := fakeUserWithEntityName(accts, ctx, "CUST", EntityCustodian)
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
			"CH admin can create CH admin", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: chAdm, PubKey: mkKey(), LegalEntityType: EntityClearingHouse,
				IsAdmin: true, LegalEntityName: "CH"}}, sdk.CodeOK,
		},
		{
			"CH admin can create CH op", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: chAdm, PubKey: mkKey(), LegalEntityType: EntityClearingHouse,
				IsAdmin: false, LegalEntityName: "CH"}}, sdk.CodeOK,
		},
		{
			"CH op cannot create CH admin", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: chOp, PubKey: mkKey(), LegalEntityType: EntityClearingHouse,
				IsAdmin: true, LegalEntityName: "CH"}}, sdk.CodeUnauthorized,
		},
		{
			"CH admin can create CUST admin", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: chAdm, PubKey: mkKey(), LegalEntityType: EntityCustodian,
				IsAdmin: true, LegalEntityName: "CUST"}}, sdk.CodeOK,
		},
		{
			"CH admin can create CUST op", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: chAdm, PubKey: mkKey(), LegalEntityType: EntityCustodian,
				IsAdmin: false, LegalEntityName: "CUST"}}, sdk.CodeUnauthorized,
		},
		{
			"CH op cannot create CUST admin", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: chOp, PubKey: mkKey(), LegalEntityType: EntityCustodian,
				IsAdmin: true, LegalEntityName: "CUST"}}, sdk.CodeUnauthorized,
		},
		{
			"CH op cannot create CUST op", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: chOp, PubKey: mkKey(), LegalEntityType: EntityCustodian,
				IsAdmin: false, LegalEntityName: "CUST"}}, sdk.CodeUnauthorized,
		},
		{
			"CUST admin can create CUST op", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: custAdm, PubKey: mkKey(), LegalEntityType: EntityCustodian,
				IsAdmin: false, LegalEntityName: "CUST"}}, sdk.CodeOK,
		},
		{
			"CUST op cannot create CUST op", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: custOp, PubKey: mkKey(), LegalEntityType: EntityCustodian,
				IsAdmin: false, LegalEntityName: "CUST"}}, sdk.CodeUnauthorized,
		},
		{
			"CUST admin cannot create member op", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: custAdm, PubKey: mkKey(), LegalEntityType: EntityIndividualClearingMember,
				IsAdmin: false, LegalEntityName: "ICM"}}, sdk.CodeUnauthorized,
		},
		{
			"CUST admin cannot create CH op", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: custAdm, PubKey: mkKey(), LegalEntityType: EntityClearingHouse,
				IsAdmin: false, LegalEntityName: "CH"}}, sdk.CodeUnauthorized,
		},
		{
			"CUST admin cannot create CH admin", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: custAdm, PubKey: mkKey(), LegalEntityType: EntityClearingHouse,
				IsAdmin: true, LegalEntityName: "CH"}}, sdk.CodeUnauthorized,
		},
		{
			"CUST admin cannot create CH admin", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: mkKey().Address(), PubKey: mkKey(), LegalEntityType: EntityClearingHouse,
				IsAdmin: true, LegalEntityName: "CH"}}, CodeInvalidAccount,
		},
		{
			"account already exists", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: custAdm, PubKey: existingCustOp.PubKey, LegalEntityType: EntityCustodian,
				IsAdmin: true, LegalEntityName: "CUST"}}, CodeInvalidAccount,
		},
		{
			"nil creator", args{ctx: ctx, msg: CreateUserAccountMsg{
				Creator: nil, PubKey: mkKey(), LegalEntityType: EntityClearingHouse,
				IsAdmin: true, LegalEntityName: "CH"}}, CodeInvalidAccount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := CreateUserAccountMsgHandler(accts)
			got := handler(tt.args.ctx, tt.args.msg)
			assert.Equal(t, tt.expect, got.Code, got.Log)

			newAcc := accts.GetAccount(ctx, tt.args.msg.(CreateUserAccountMsg).PubKey.Address())
			if tt.expect == sdk.CodeOK || tt.name == "account already exists" {
				assert.True(t, newAcc != nil)
			} else {
				assert.True(t, newAcc == nil)
			}

		})
	}
}

func Test_createAssetAccountMsgHandler_Do(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	createAccount := func(typ, entityName string, isAdm bool) crypto.Address {
		if isAdm {
			ac, _ := fakeAdminWithEntityName(accts, ctx, entityName, typ)
			return ac.PubKey.Address()
		}
		ac, _ := fakeUserWithEntityName(accts, ctx, entityName, typ)
		return ac.PubKey.Address()
	}
	mkKey := func() crypto.PubKey {
		return crypto.GenPrivKeyEd25519().PubKey()
	}
	chOp := createAccount(EntityClearingHouse, "CH", false)
	chAdm := createAccount(EntityClearingHouse, "CH", true)
	custOp := createAccount(EntityCustodian, "CUST", false)
	custAdm := createAccount(EntityCustodian, "CUST", true)
	tests := []struct {
		name   string
		msg    sdk.Msg
		expect sdk.CodeType
	}{
		{"CH admin can create CH asset", CreateAssetAccountMsg{Creator: chAdm, PubKey: mkKey()}, sdk.CodeOK},
		{"CH op cannot create CH asset", CreateAssetAccountMsg{Creator: chOp, PubKey: mkKey()}, CodeWrongSigner},
		{"CUST admin can create CUST asset", CreateAssetAccountMsg{Creator: custAdm, PubKey: mkKey()}, sdk.CodeOK},
		{"CUST op cannot create CUST asset", CreateAssetAccountMsg{Creator: custOp, PubKey: mkKey()}, CodeWrongSigner},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := CreateAssetAccountMsgHandler(accts)
			got := handler(ctx, tt.msg)
			assert.Equal(t, tt.expect, got.Code, got.Log)

			concreteMsg := tt.msg.(CreateAssetAccountMsg)
			newAcc := accts.GetAccount(ctx, concreteMsg.PubKey.Address())
			if tt.expect == sdk.CodeOK || tt.name == "account already exists" {
				creator := accts.GetAccount(ctx, concreteMsg.Creator).(*AppAccount)
				assert.True(t, newAcc != nil)
				concAcc := newAcc.(*AppAccount)
				assert.Equal(t, concAcc.BaseLegalEntity, creator.BaseLegalEntity)
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

func fakeUser(accts sdk.AccountMapper, ctx sdk.Context, typ string) (*AppAccount, crypto.PrivKey) {
	return fakeUserWithEntityName(accts, ctx, typ, typ)
}

func fakeInactiveUser(accts sdk.AccountMapper, ctx sdk.Context, typ string) (*AppAccount, crypto.PrivKey) {
	acct, priv := makeUser(typ, typ)
	acct.Active = false
	accts.SetAccount(ctx, acct)
	return acct, priv
}

func fakeUserWithEntityName(accts sdk.AccountMapper, ctx sdk.Context, entname, typ string) (*AppAccount, crypto.PrivKey) {
	acct, priv := makeUser(entname, typ)
	accts.SetAccount(ctx, acct)
	return acct, priv
}

func fakeAdmin(accts sdk.AccountMapper, ctx sdk.Context, typ string) (*AppAccount, crypto.PrivKey) {
	return fakeAdminWithEntityName(accts, ctx, typ, typ)
}

func fakeAdminWithEntityName(accts sdk.AccountMapper, ctx sdk.Context, entname, typ string) (*AppAccount, crypto.PrivKey) {
	acct, priv := makeAdminUser(entname, typ)
	accts.SetAccount(ctx, acct)
	return acct, priv
}

func fakeAsset(accts sdk.AccountMapper, ctx sdk.Context, cash sdk.Coins, typ string) (*AppAccount, crypto.Address) {
	return fakeAssetWithEntityName(accts, ctx, cash, typ, typ)
}

func fakeAssetWithEntityName(accts sdk.AccountMapper, ctx sdk.Context, cash sdk.Coins, entname, typ string) (*AppAccount, crypto.Address) {
	acct, addr := makeAssetAccount(cash, entname, typ)
	accts.SetAccount(ctx, acct)
	return acct, addr
}
