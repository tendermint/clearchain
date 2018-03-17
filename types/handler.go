package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

// RegisterRoutes routes the message (request) to a proper handler.
func RegisterRoutes(r baseapp.Router, accts sdk.AccountMapper) {
	r.AddRoute(DepositType, DepositMsgHandler(accts)).
		AddRoute(SettlementType, SettleMsgHandler(accts)).
		AddRoute(WithdrawType, WithdrawMsgHandler(accts)).
		AddRoute(CreateOperatorType, CreateOperatorMsgHandler(accts)).
		AddRoute(CreateAdminType, CreateAdminMsgHandler(accts)).
		AddRoute(CreateAssetAccountType, CreateAssetAccountMsgHandler(accts)).
		AddRoute(FreezeOperatorType, FreezeOperatorMsgHandler(accts)).
		AddRoute(FreezeAdminType, FreezeAdminMsgHandler(accts))
}

/*

Deposit functionality.

Sender is Custodian
Rec is Member
*/
func DepositMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return depositMsgHandler{accts}.Do
}

type depositMsgHandler struct{ accts sdk.AccountMapper }

// Deposit logic
func (d depositMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	// ensure proper message
	dm, ok := msg.(DepositMsg)
	if !ok {
		return ErrWrongMsgFormat("expected DepositMsg").Result()
	}
	// ensure proper types
	if _, err := getCHActiveOperator(ctx, d.accts, dm.Operator); err != nil {
		return err.Result()
	}
	sender, err := getActiveAssetWithEntityType(ctx, d.accts, dm.Sender, IsCustodian)
	if err != nil {
		return err.Result()
	}
	rcpt, err := getActiveAssetWithEntityType(ctx, d.accts, dm.Recipient, IsMember)
	if err != nil {
		return err.Result()
	}
	// Exchange cash
	if err := moveMoney(d.accts, ctx, sender, rcpt, dm.Amount, false, true); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

/*
Settlement funcionality.

Operator is CH
Sender is CH
Rec is member
*/
func SettleMsgHandler(accts sdk.AccountMapper) sdk.Handler { return settleMsgHandler{accts}.Do }

type settleMsgHandler struct{ accts sdk.AccountMapper }

// Settlement logic
func (sh settleMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	//ensure proper message
	sm, ok := msg.(SettleMsg)
	if !ok {
		return ErrWrongMsgFormat("expected SettleMsg").Result()
	}
	// ensure proper types
	operator, err := getCHActiveOperator(ctx, sh.accts, sm.Operator)
	if err != nil {
		return err.Result()
	}
	sender, err := getActiveAssetWithEntityType(ctx, sh.accts, sm.Sender, IsClearingHouse)
	if err != nil {
		return err.Result()
	}
	if !BelongToSameEntity(operator, sender) {
		return ErrWrongSigner("operator and sender must belong to the same entity").Result()
	}
	rcpt, err := getActiveAssetWithEntityType(ctx, sh.accts, sm.Recipient, IsMember)
	if err != nil {
		return err.Result()
	}
	if err := moveMoney(sh.accts, ctx, sender, rcpt, sm.Amount, false, true); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// WithdrawMsgHandler implements the withdraw functionality.
//
// Sender is member
// Reci is custodian
// Operator is CH
//
func WithdrawMsgHandler(accts sdk.AccountMapper) sdk.Handler { return withdrawMsgHandler{accts}.Do }

type withdrawMsgHandler struct {
	accts sdk.AccountMapper
}

// Withdraw logic
func (wh withdrawMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	// ensure proper message
	wm, ok := msg.(WithdrawMsg)
	if !ok {
		return ErrWrongMsgFormat("expected WithdrawMsg").Result()
	}
	// ensure proper types
	_, err := getCHActiveOperator(ctx, wh.accts, wm.Operator)
	if err != nil {
		return err.Result()
	}
	sender, err := getActiveAssetWithEntityType(ctx, wh.accts, wm.Sender, IsMember)
	if err != nil {
		return err.Result()
	}
	rcpt, err := getActiveAssetWithEntityType(ctx, wh.accts, wm.Recipient, IsCustodian)
	if err != nil {
		return err.Result()
	}
	err = moveMoney(wh.accts, ctx, sender, rcpt, wm.Amount, true, false)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}

}

// CreateOperatorMsgHandler returns the handler's method.
func CreateOperatorMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return createOperatorMsgHandler{accts}.Do
}

type createOperatorMsgHandler struct{ accts sdk.AccountMapper }

func (h createOperatorMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	// ensure proper message
	cm, ok := msg.(CreateOperatorMsg)
	if !ok {
		return ErrWrongMsgFormat("expected CreateOperatorMsg").Result()
	}
	newAcct, err := validateAdminAndCreateOperator(ctx, h.accts, cm.Creator, cm.PubKey)
	if err != nil {
		return err.Result()
	}
	h.accts.SetAccount(ctx, newAcct)
	return sdk.Result{}
}

// CreateAdminMsgHandler returns the handler's method.
func CreateAdminMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return createAdminMsgHandler{accts}.Do
}

type createAdminMsgHandler struct{ accts sdk.AccountMapper }

func (h createAdminMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	// ensure proper message
	cm, ok := msg.(CreateAdminMsg)
	if !ok {
		return ErrWrongMsgFormat("expected CreateAdminMsg").Result()
	}
	newAcct, err := validateCHAdminAndCreateXEntityAdmin(ctx, h.accts, cm.Creator, cm.PubKey, cm.BaseLegalEntity)
	if err != nil {
		return err.Result()
	}
	h.accts.SetAccount(ctx, newAcct)
	return sdk.Result{}
}

// CreateAssetAccountMsgHandler returns the handler's method.
func CreateAssetAccountMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return createAssetAccountMsgHandler{accts}.Do
}

type createAssetAccountMsgHandler struct{ accts sdk.AccountMapper }

// Create asset account logic.
// Admins can create asset accounts for their own entity only.
// TODO: clarify business rules and ensure this is desired
func (h createAssetAccountMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	// ensure proper message
	cm, ok := msg.(CreateAssetAccountMsg)
	if !ok {
		return ErrWrongMsgFormat("expected CreateAssetAccountMsg").Result()
	}
	// ensure creator exists
	// no need for type checking, CreateAssetAccount
	// validates types too.
	creator, err := getActiveAdmin(ctx, h.accts, cm.Creator)
	if err != nil {
		return err.Result()
	}
	// ensure new account does not exist
	if h.accts.GetAccount(ctx, cm.PubKey.Address()) != nil {
		return ErrInvalidAccount("the account already exists").Result()
	}
	// Construct a new account
	newAcct := NewAssetAccount(cm.PubKey, sdk.Coins{}, creator.Address, creator.LegalEntityName(), creator.LegalEntityType())
	h.accts.SetAccount(ctx, newAcct)
	return sdk.Result{}
}

// FreezeOperatorMsgHandler returns the handler's method.
func FreezeOperatorMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return freezeOperatorMsgHandler{accts}.Do
}

type freezeOperatorMsgHandler struct{ accts sdk.AccountMapper }

// Freeze operator's message logic.
// Admins can freeze their own entity's operator accounts.
func (h freezeOperatorMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	// ensure proper message
	cm, ok := msg.(FreezeOperatorMsg)
	if !ok {
		return ErrWrongMsgFormat("expected FreezeOperatorMsg").Result()
	}
	// ensure admin exists
	admin, err := getActiveAdmin(ctx, h.accts, cm.Admin)
	if err != nil {
		return err.Result()
	}
	// ensure operator exists
	operator, err := getActiveOperator(ctx, h.accts, cm.Target)
	if err != nil {
		return err.Result()
	}
	if !BelongToSameEntity(admin, operator) {
		return ErrWrongSigner("admin and operator do not belong to the same entity").Result()
	}
	// Construct a new account
	operator.Active = false
	h.accts.SetAccount(ctx, operator)
	return sdk.Result{}
}

// FreezeAdminMsgHandler returns the handler's method.
func FreezeAdminMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return freezeAdminMsgHandler{accts}.Do
}

type freezeAdminMsgHandler struct{ accts sdk.AccountMapper }

// Freeze admin's message logic.
// Clearing house Admins can freeze any other admin.
func (h freezeAdminMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	// ensure proper message
	cm, ok := msg.(FreezeAdminMsg)
	if !ok {
		return ErrWrongMsgFormat("expected FreezeAdminMsg").Result()
	}
	// ensure clearing house admin exists
	if _, err := getCHActiveAdmin(ctx, h.accts, cm.Admin); err != nil {
		return err.Result()
	}
	// ensure target admin exists
	admin, err := getActiveAdmin(ctx, h.accts, cm.Target)
	if err != nil {
		return err.Result()
	}
	// Construct a new account
	admin.Active = false
	h.accts.SetAccount(ctx, admin)
	return sdk.Result{}
}

// Business logic

func validateAdminAndCreateOperator(ctx sdk.Context, accts sdk.AccountMapper,
	creatorAddr crypto.Address, pub crypto.PubKey) (*AppAccount, sdk.Error) {
	creator, err := getActiveAdmin(ctx, accts, creatorAddr)
	if err != nil {
		return nil, err
	}
	// ensure new account does not exist
	if accts.GetAccount(ctx, pub.Address()) != nil {
		return nil, ErrInvalidAccount("couldn't create the account, it already exists")
	}
	return NewOpUser(pub, creator.GetAddress(), creator.LegalEntityName(), creator.LegalEntityType()), nil
}

func validateCHAdminAndCreateXEntityAdmin(ctx sdk.Context, accts sdk.AccountMapper,
	creatorAddr crypto.Address, pub crypto.PubKey, ent LegalEntity) (*AppAccount, sdk.Error) {
	if _, err := getCHActiveAdmin(ctx, accts, creatorAddr); err != nil {
		return nil, err
	}
	// ensure new account does not exist
	if accts.GetAccount(ctx, pub.Address()) != nil {
		return nil, ErrInvalidAccount("couldn't create the account, it already exists")
	}
	return NewAdminUser(pub, creatorAddr, ent.LegalEntityName(), ent.LegalEntityType()), nil
}

// Transfers money from the sender to the  recipient
func moveMoney(accts sdk.AccountMapper, ctx sdk.Context, sender *AppAccount, recipient *AppAccount,
	amount sdk.Coin, senderMustBePositive bool, recipientMustBePositive bool) sdk.Error {
	transfer := sdk.Coins{amount}
	// first verify funds
	sender.Coins = sender.Coins.Minus(transfer)
	if senderMustBePositive && !sender.Coins.IsNotNegative() {
		return ErrInvalidAmount("sender  has insufficient funds")
	}
	// transfer may be negative
	recipient.Coins = recipient.Coins.Plus(transfer)
	if recipientMustBePositive && !recipient.Coins.IsNotNegative() {
		return ErrInvalidAmount("recipient has insufficient funds")
	}
	// now make the transfer and save the result
	accts.SetAccount(ctx, sender)
	accts.SetAccount(ctx, recipient)
	return nil
}

// Auxiliary functions

func getCHActiveAdmin(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address) (*AppAccount, sdk.Error) {
	return getUserAccountWithGetterAndEntityType(ctx, accts, addr, getActiveAdmin, IsClearingHouse)
}

func getCHActiveOperator(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address) (*AppAccount, sdk.Error) {
	return getUserAccountWithGetterAndEntityType(ctx, accts, addr, getActiveOperator, IsClearingHouse)
}

func getUserAccountWithGetterAndEntityType(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address,
	accGetter func(sdk.Context, sdk.AccountMapper, crypto.Address) (*AppAccount, sdk.Error),
	entityTypeCheck func(LegalEntity) bool) (*AppAccount, sdk.Error) {
	account, err := accGetter(ctx, accts, addr)
	if err != nil {
		return nil, err
	}
	if !entityTypeCheck(account) {
		return nil, ErrWrongSigner(account.LegalEntityType())
	}
	return account, nil
}

func getActiveAssetWithEntityType(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address,
	entityTypeCheck func(LegalEntity) bool) (*AppAccount, sdk.Error) {
	account, err := getActiveAsset(ctx, accts, addr)
	if err != nil {
		return nil, err
	}
	if !entityTypeCheck(account) {
		return nil, ErrWrongSigner(account.LegalEntityType())
	}
	return account, nil
}

func getActiveAdmin(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address) (*AppAccount, sdk.Error) {
	return getUser(ctx, accts, addr, true, true)
}

func getActiveOperator(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address) (*AppAccount, sdk.Error) {
	return getUser(ctx, accts, addr, true, false)
}

func getUser(ctx sdk.Context, accts sdk.AccountMapper,
	addr crypto.Address, wantActive, wantAdmin bool) (*AppAccount, sdk.Error) {
	rawAccount := accts.GetAccount(ctx, addr)
	if rawAccount == nil {
		return nil, ErrInvalidAccount("account does not exist")
	}
	account := rawAccount.(*AppAccount)
	if !IsUser(account) {
		return nil, ErrWrongSigner("invalid account type")
	}
	if wantActive && !account.Active {
		return nil, ErrInactiveUser(fmt.Sprintf("%v", addr))
	}
	if !wantActive && account.Active {
		return nil, ErrInactiveUser(fmt.Sprintf("%v", addr))
	}
	if wantAdmin && !account.IsAdmin() {
		return nil, ErrWrongSigner("must be admin")
	}
	if !wantAdmin && account.IsAdmin() {
		return nil, ErrWrongSigner("must not be admin")
	}
	return account, nil
}

func getActiveAsset(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address) (*AppAccount, sdk.Error) {
	return getAsset(ctx, accts, addr, true)
}

func getAsset(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address, wantActive bool) (*AppAccount, sdk.Error) {
	rawAccount := accts.GetAccount(ctx, addr)
	if rawAccount == nil {
		return nil, ErrInvalidAccount("account does not exist")
	}
	account := rawAccount.(*AppAccount)
	if !IsAsset(account) {
		return nil, ErrWrongSigner("invalid account type")
	}
	if wantActive && !account.Active {
		return nil, ErrInactiveUser(fmt.Sprintf("%v", addr))
	}
	if !wantActive && account.Active {
		return nil, ErrInactiveUser(fmt.Sprintf("%v", addr))
	}
	return account, nil
}
