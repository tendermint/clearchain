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
			sdk.CodeOK, sdk.Coins{{"USD", -5000}}, sdk.Coins{{"USD", 5000}}, sdk.Coins{}, sdk.Coins{},
		},
		{
			"deposit2",
			args{ctx: ctx, msg: DepositMsg{Operator: op, Sender: cust, Recipient: member2, Amount: sdk.Coin{"USD", 7777}}},
			sdk.CodeOK,
			sdk.Coins{{"USD", -12777}},
			sdk.Coins{{"USD", 5000}},
			sdk.Coins{{"USD", 7777}},
			sdk.Coins{},
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
			sdk.Coins{},
		},
		{
			"withdraw",
			args{ctx: ctx, msg: WithdrawMsg{Operator: op, Sender: member, Recipient: cust, Amount: sdk.Coin{"USD", 5500}}},
			sdk.CodeOK,
			sdk.Coins{{"USD", -7277}},
			sdk.Coins{{"USD", 2500}},
			sdk.Coins{{"USD", 4777}},
			sdk.Coins{},
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
	_, inactiveAssetAddr := fakeInactiveAssetWithEntityName(accts, ctx, cCoins, EntityGeneralClearingMember, EntityGeneralClearingMember)

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
			"admins ain't allowed", args{ctx: ctx, msg: DepositMsg{Operator: chAdmAddr, Sender: cust,
				Recipient: member, Amount: sdk.Coin{"USD", 700}}}, CodeWrongSigner, cCoins, sdk.Coins{},
		},
		{
			"inactive operator", args{ctx: ctx, msg: DepositMsg{Operator: inactiveOp.Address, Sender: cust,
				Recipient: member, Amount: sdk.Coin{"USD", 700}}}, CodeInactiveAccount, cCoins, sdk.Coins{},
		},
		{
			"no returns", args{ctx: ctx, msg: DepositMsg{Operator: chOp, Sender: member,
				Recipient: cust, Amount: sdk.Coin{"USD", 200}}}, CodeWrongSigner, cCoins, sdk.Coins{}, // sdk.Coins{}
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
		{
			// want an active asset for the recipient
			"overdraft", args{ctx: ctx, msg: DepositMsg{Operator: chOp, Sender: cust, Recipient: inactiveAssetAddr, Amount: sdk.Coin{"EUR", 10000}}},
			CodeInactiveAccount, sdk.Coins{{"EUR", -5000}, {"USD", 300}}, sdk.Coins{{"EUR", 10000}, {"USD", 700}},
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
	_, clh := fakeAssetWithEntityName(accts, ctx, clhCoins, chOpAcc.LegalEntityType(), EntityClearingHouse)
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
			CodeWrongSigner, sdk.Coins{}, mCoins,
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
	existing, _ := fakeAsset(accts, ctx, nil, "CUST")
	tests := []struct {
		name   string
		msg    sdk.Msg
		expect sdk.CodeType
	}{
		{"CH admin can create CH asset", CreateAssetAccountMsg{Creator: chAdm, PubKey: mkKey()}, sdk.CodeOK},
		{"CH op cannot create CH asset", CreateAssetAccountMsg{Creator: chOp, PubKey: mkKey()}, CodeWrongSigner},
		{"CUST admin can create CUST asset", CreateAssetAccountMsg{Creator: custAdm, PubKey: mkKey()}, sdk.CodeOK},
		{"CUST op cannot create CUST asset", CreateAssetAccountMsg{Creator: custOp, PubKey: mkKey()}, CodeWrongSigner},
		{"account already exists", CreateAssetAccountMsg{Creator: custAdm, PubKey: existing.PubKey}, CodeInvalidAccount},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := CreateAssetAccountMsgHandler(accts)
			got := handler(ctx, tt.msg)
			assert.Equal(t, tt.expect, got.Code, got.Log)

			concreteMsg := tt.msg.(CreateAssetAccountMsg)
			newAcc := accts.GetAccount(ctx, concreteMsg.PubKey.Address())
			if tt.name == "account already exists" {
				assert.True(t, newAcc != nil)
			} else if tt.expect == sdk.CodeOK {
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

func Test_validateAdminAndCreateOperator(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	newAccPub := crypto.GenPrivKeyEd25519().PubKey()
	inactiveAdm, _ := fakeInactiveAdmin(accts, ctx, EntityClearingHouse)
	admin, _ := fakeAdmin(accts, ctx, EntityClearingHouse)
	opCreated := NewOpUser(newAccPub, admin.Address, admin.LegalEntityName(), admin.LegalEntityType())
	operator, _ := fakeUser(accts, ctx, EntityClearingHouse)
	type args struct {
		creatorAddr crypto.Address
		pub         crypto.PubKey
	}
	tests := []struct {
		name  string
		args  args
		want  *AppAccount
		want1 sdk.CodeType
	}{
		{"inactive admin cannot create", args{inactiveAdm.Address, newAccPub}, nil, CodeInactiveAccount},
		{"operator cannot create operator", args{operator.Address, newAccPub}, nil, CodeWrongSigner},
		{"admin can create operator", args{admin.Address, newAccPub}, opCreated, sdk.CodeOK},
		{"account already exists", args{admin.Address, operator.PubKey}, nil, CodeInvalidAccount},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := validateAdminAndCreateOperator(ctx, accts, tt.args.creatorAddr, tt.args.pub)
			if got == nil || tt.want == nil {
				assert.True(t, got == tt.want)
			} else {
				assert.True(t, accountEqual(tt.want, got))
			}
			if got1 == nil {
				assert.Equal(t, tt.want1, sdk.CodeOK)
			} else {
				assert.Equal(t, tt.want1, got1.ABCICode(), got1.ABCILog())
			}
		})
	}
}

func Test_validateCHAdminAndCreateXEntityAdmin(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	ent := BaseLegalEntity{EntityName: "ICM", EntityType: EntityIndividualClearingMember}
	newAccPub := crypto.GenPrivKeyEd25519().PubKey()
	admin, _ := fakeAdmin(accts, ctx, EntityClearingHouse)
	inactiveAdmin, _ := fakeInactiveAdmin(accts, ctx, EntityClearingHouse)
	adminCreated := NewAdminUser(newAccPub, admin.Address, ent.LegalEntityName(), ent.LegalEntityType())
	icmAdmin, _ := fakeAdmin(accts, ctx, EntityIndividualClearingMember)
	chOperator, _ := fakeUser(accts, ctx, EntityClearingHouse)
	custOperator, _ := fakeUser(accts, ctx, EntityCustodian)
	type args struct {
		creatorAddr crypto.Address
		pub         crypto.PubKey
		ent         LegalEntity
	}
	tests := []struct {
		name  string
		args  args
		want  *AppAccount
		want1 sdk.CodeType
	}{
		{"foreign operator cannot create", args{custOperator.Address, newAccPub, ent}, nil, CodeWrongSigner},
		{"foreign admin cannot create", args{icmAdmin.Address, newAccPub, ent}, nil, CodeWrongSigner},
		{"operator cannot create admin", args{chOperator.Address, newAccPub, ent}, nil, CodeWrongSigner},
		{"admin can create admin", args{admin.Address, newAccPub, ent}, adminCreated, sdk.CodeOK},
		{"existing account", args{admin.Address, icmAdmin.PubKey, icmAdmin.BaseLegalEntity}, nil, CodeInvalidAccount},
		{"inactive admin cannot create", args{inactiveAdmin.Address, newAccPub, ent}, nil, CodeInactiveAccount},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := validateCHAdminAndCreateXEntityAdmin(ctx, accts, tt.args.creatorAddr, tt.args.pub, tt.args.ent)
			if got == nil || tt.want == nil {
				assert.True(t, got == tt.want)
			} else {
				assert.True(t, accountEqual(tt.want, got))
			}
			if got1 == nil {
				assert.Equal(t, tt.want1, sdk.CodeOK)
			} else {
				assert.Equal(t, tt.want1, got1.ABCICode(), got1.ABCILog())
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

	accts := NewAccountMapper(key)
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

func fakeInactiveAdmin(accts sdk.AccountMapper, ctx sdk.Context, typ string) (*AppAccount, crypto.PrivKey) {
	acct, priv := makeAdminUser(typ, typ)
	acct.Active = false
	accts.SetAccount(ctx, acct)
	return acct, priv
}

func fakeAdminWithEntityName(accts sdk.AccountMapper, ctx sdk.Context, entname, typ string) (*AppAccount, crypto.PrivKey) {
	acct, priv := makeAdminUser(entname, typ)
	accts.SetAccount(ctx, acct)
	return acct, priv
}

func fakeAsset(accts sdk.AccountMapper, ctx sdk.Context, cash sdk.Coins, typ string) (*AppAccount, crypto.Address) {
	return fakeAssetWithEntityName(accts, ctx, cash, typ, typ)
}

func fakeInactiveAssetWithEntityName(accts sdk.AccountMapper, ctx sdk.Context,
	cash sdk.Coins, entname, typ string) (*AppAccount, crypto.Address) {
	acct, addr := makeAssetAccount(cash, entname, typ)
	acct.Active = false
	accts.SetAccount(ctx, acct)
	return acct, addr
}

func fakeAssetWithEntityName(accts sdk.AccountMapper, ctx sdk.Context, cash sdk.Coins, entname, typ string) (*AppAccount, crypto.Address) {
	acct, addr := makeAssetAccount(cash, entname, typ)
	accts.SetAccount(ctx, acct)
	return acct, addr
}

func Test_createOperatorMsgHandler_Do(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	newPub := crypto.GenPrivKeyEd25519().PubKey()
	admin, _ := fakeAdminWithEntityName(accts, ctx, "member", EntityIndividualClearingMember)
	operator, _ := fakeUser(accts, ctx, EntityCustodian)
	asset, _ := fakeAsset(accts, ctx, nil, EntityGeneralClearingMember)
	admAddr := admin.Address
	opAddr := operator.Address
	assetAddr := asset.Address
	tests := []struct {
		name string
		msg  BaseCreateUserMsg
		want sdk.CodeType
	}{
		{"admin can create", BaseCreateUserMsg{Creator: admAddr, PubKey: newPub}, sdk.CodeOK},
		{"already exists", BaseCreateUserMsg{Creator: admAddr, PubKey: admin.PubKey}, CodeInvalidAccount},
		{"operator cannot create", BaseCreateUserMsg{Creator: opAddr, PubKey: newPub}, CodeWrongSigner},
		{"asset cannot create", BaseCreateUserMsg{Creator: assetAddr, PubKey: newPub}, CodeWrongSigner},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := CreateOperatorMsg{tt.msg}
			h := createOperatorMsgHandler{
				accts: accts,
			}
			assert.Equal(t, tt.want, h.Do(ctx, msg).Code)
		})
	}
}

func Test_createAdminMsgHandler_Do(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	newPub := crypto.GenPrivKeyEd25519().PubKey()
	chAdmin, _ := fakeAdminWithEntityName(accts, ctx, "clearing house", EntityClearingHouse)
	nonchAdmin, _ := fakeAdminWithEntityName(accts, ctx, "member", EntityIndividualClearingMember)
	operator, _ := fakeUser(accts, ctx, EntityCustodian)
	asset, _ := fakeAsset(accts, ctx, nil, EntityGeneralClearingMember)
	chAdmAddr := chAdmin.Address
	nonchAdmAddr := nonchAdmin.Address
	opAddr := operator.Address
	assetAddr := asset.Address
	tests := []struct {
		name string
		msg  CreateAdminMsg
		want sdk.CodeType
	}{
		{"non CH admin cannot create", CreateAdminMsg{BaseCreateUserMsg{
			Creator: nonchAdmAddr, PubKey: newPub}, BaseLegalEntity{}}, CodeWrongSigner},
		{"CH admin can create", CreateAdminMsg{BaseCreateUserMsg{
			Creator: chAdmAddr, PubKey: newPub}, BaseLegalEntity{}}, sdk.CodeOK},
		{"operator cannot create", CreateAdminMsg{BaseCreateUserMsg{
			Creator: opAddr, PubKey: newPub}, BaseLegalEntity{}}, CodeWrongSigner},
		{"asset cannot create", CreateAdminMsg{BaseCreateUserMsg{
			Creator: assetAddr, PubKey: newPub}, BaseLegalEntity{}}, CodeWrongSigner},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := createAdminMsgHandler{
				accts: accts,
			}
			got := h.Do(ctx, tt.msg)
			assert.Equal(t, tt.want, got.Code, got.Log)
		})
	}
}

func Test_freezeOperatorMsgHandler_Do(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	admin, _ := fakeAdminWithEntityName(accts, ctx, "qweasdzxc", EntityIndividualClearingMember)
	admAddr := admin.Address
	inactiveAdmin, _ := fakeAdminWithEntityName(accts, ctx, "qweasdzxc", EntityIndividualClearingMember)
	inactiveAdmin.Active = false
	accts.SetAccount(ctx, inactiveAdmin)
	inadmAddr := inactiveAdmin.Address
	operator, _ := fakeUserWithEntityName(accts, ctx, admin.LegalEntityName(), EntityIndividualClearingMember)
	opAddr := operator.Address
	inactiveOperator, _ := fakeUserWithEntityName(accts, ctx, admin.LegalEntityName(), EntityIndividualClearingMember)
	inactiveOperator.Active = false
	accts.SetAccount(ctx, inactiveOperator)
	inopAddr := inactiveOperator.Address
	foreignOperator, _ := fakeUser(accts, ctx, EntityIndividualClearingMember)
	fopAddr := foreignOperator.Address
	asset, _ := fakeAsset(accts, ctx, nil, EntityGeneralClearingMember)
	assetAddr := asset.Address
	tests := []struct {
		name string
		msg  BaseFreezeAccountMsg
		want sdk.CodeType
	}{
		{"inactive admin", BaseFreezeAccountMsg{Admin: inadmAddr, Target: opAddr}, CodeInactiveAccount},
		{"inactive op", BaseFreezeAccountMsg{Admin: admAddr, Target: inopAddr}, CodeInactiveAccount},
		{"op can't freeze", BaseFreezeAccountMsg{Admin: opAddr}, CodeWrongSigner},
		{"foreign operator", BaseFreezeAccountMsg{Admin: admAddr, Target: fopAddr}, CodeWrongSigner},
		{"invalid account", BaseFreezeAccountMsg{Admin: admAddr, Target: assetAddr}, CodeWrongSigner},
		{"ok", BaseFreezeAccountMsg{Admin: admAddr, Target: opAddr}, sdk.CodeOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := freezeOperatorMsgHandler{accts: accts}
			got := h.Do(ctx, FreezeOperatorMsg{tt.msg})
			assert.Equal(t, tt.want, got.Code, got.Log)
			if got.Code == sdk.CodeOK {
				acct := accts.GetAccount(ctx, tt.msg.Target)
				ca := acct.(*AppAccount)
				assert.False(t, ca.IsActive())
			}
		})
	}
}

func Test_freezeAdminMsgHandler_Do(t *testing.T) {
	accts, ctx := fakeAccountMapper()
	chAdmin, _ := fakeAdmin(accts, ctx, EntityClearingHouse)
	chAdmAddr := chAdmin.Address
	inactiveChAdmin, _ := fakeAdminWithEntityName(accts, ctx, chAdmin.EntityName, EntityClearingHouse)
	inactiveChAdmin.Active = false
	accts.SetAccount(ctx, inactiveChAdmin)
	inChAdmAddr := inactiveChAdmin.Address
	operator, _ := fakeUser(accts, ctx, EntityIndividualClearingMember)
	opAddr := operator.Address
	foreignAdmin, _ := fakeAdmin(accts, ctx, EntityIndividualClearingMember)
	foaAddr := foreignAdmin.Address
	foreignInactiveAdmin, _ := fakeAdmin(accts, ctx, EntityIndividualClearingMember)
	foreignInactiveAdmin.Active = false
	accts.SetAccount(ctx, foreignInactiveAdmin)
	foaInAddr := foreignInactiveAdmin.Address
	asset, _ := fakeAsset(accts, ctx, nil, EntityGeneralClearingMember)
	assetAddr := asset.Address
	tests := []struct {
		name string
		msg  BaseFreezeAccountMsg
		want sdk.CodeType
	}{
		{"can't freeze op", BaseFreezeAccountMsg{Admin: inChAdmAddr, Target: foaAddr}, CodeInactiveAccount},
		{"inactive admin", BaseFreezeAccountMsg{Admin: inChAdmAddr, Target: foaAddr}, CodeInactiveAccount},
		{"inactive target", BaseFreezeAccountMsg{Admin: chAdmAddr, Target: foaInAddr}, CodeInactiveAccount},
		{"op can't freeze", BaseFreezeAccountMsg{Admin: opAddr}, CodeWrongSigner},
		{"invalid account", BaseFreezeAccountMsg{Admin: chAdmAddr, Target: assetAddr}, CodeWrongSigner},
		{"can't freeze op", BaseFreezeAccountMsg{Admin: chAdmAddr, Target: opAddr}, CodeWrongSigner},
		{"ok", BaseFreezeAccountMsg{Admin: chAdmAddr, Target: foaAddr}, sdk.CodeOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := freezeAdminMsgHandler{accts: accts}
			got := h.Do(ctx, FreezeAdminMsg{tt.msg})
			assert.Equal(t, tt.want, got.Code, got.Log)
			if got.Code == sdk.CodeOK {
				acct := accts.GetAccount(ctx, tt.msg.Target)
				ca := acct.(*AppAccount)
				assert.False(t, ca.IsActive())
			}
		})
	}
}
