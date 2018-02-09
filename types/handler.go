package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

// RegisterRoutes routes the message (request) to a proper handler.
func RegisterRoutes(r baseapp.Router, accts sdk.AccountMapper) {
	r.AddRoute(DepositType, DepositMsgHandler(accts))
	r.AddRoute(SettlementType, SettleMsgHandler(accts))
	r.AddRoute(WithdrawType, WithdrawMsgHandler(accts))
	r.AddRoute(CreateUserAccountType, CreateUserAccountMsgHandler(accts))
	r.AddRoute(CreateAssetAccountType, CreateAssetAccountMsgHandler(accts))
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
	if _, err := getCHOperator(ctx, d.accts, dm.Operator); err != nil {
		return err.Result()
	}
	sender, err := getAssetWithEntityType(ctx, d.accts, dm.Sender, IsCustodian)
	if err != nil {
		return err.Result()
	}
	rcpt, err := getAssetWithEntityType(ctx, d.accts, dm.Recipient, IsMember)
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
func SettleMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return settleMsgHandler{accts}.Do
}

type settleMsgHandler struct{ accts sdk.AccountMapper }

// Settlement logic
func (sh settleMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	//ensure proper message
	sm, ok := msg.(SettleMsg)
	if !ok {
		return ErrWrongMsgFormat("expected SettleMsg").Result()
	}
	// ensure proper types
	operator, err := getCHOperator(ctx, sh.accts, sm.Operator)
	if err != nil {
		return err.Result()
	}
	sender, err := getAssetWithEntityType(ctx, sh.accts, sm.Sender, IsClearingHouse)
	if err != nil {
		return err.Result()
	}
	if !BelongToSameEntity(operator, sender) {
		return ErrWrongSigner("operator and sender must belong to the same entity").Result()
	}
	rcpt, err := getAssetWithEntityType(ctx, sh.accts, sm.Recipient, IsMember)
	if err != nil {
		return err.Result()
	}
	if err := moveMoney(sh.accts, ctx, sender, rcpt, sm.Amount, false, true); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

/*

Withdraw functionality.

Sender is member
Reci is custodian
Operator is CH

*/
func WithdrawMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return withdrawMsgHandler{accts}.Do
}

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
	_, err := getCHOperator(ctx, wh.accts, wm.Operator)
	if err != nil {
		return err.Result()
	}
	sender, err := getAssetWithEntityType(ctx, wh.accts, wm.Sender, IsMember)
	if err != nil {
		return err.Result()
	}
	rcpt, err := getAssetWithEntityType(ctx, wh.accts, wm.Recipient, IsCustodian)
	if err != nil {
		return err.Result()
	}
	err = moveMoney(wh.accts, ctx, sender, rcpt, wm.Amount, true, false)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}

}

// CreateUserAccountMsgHandler returns the handler's method.
func CreateUserAccountMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return createUserAccountMsgHandler{accts}.Do
}

type createUserAccountMsgHandler struct{ accts sdk.AccountMapper }

// Create acc logic.
// A clearing house account is allowed to create any kind of accounts,
// including clearing house, custodian, and members accounts.
// TODO: clarify business rules and ensure this is desired
func (h createUserAccountMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	var newAcct *AppAccount
	// ensure proper message
	cm, ok := msg.(CreateUserAccountMsg)
	if !ok {
		return ErrWrongMsgFormat("expected CreateUserAccountMsg").Result()
	}
	// ensure creator exists
	creator := h.accts.GetAccount(ctx, cm.Creator)
	if creator == nil {
		return ErrInvalidAccount("couldn't find creator").Result()
	}
	// ensure proper types
	concreteCreator := creator.(*AppAccount)
	if !concreteCreator.IsActive() {
		ErrWrongSigner("the creator account must be active")
	}
	// ensure new account does not exist
	if h.accts.GetAccount(ctx, cm.PubKey.Address()) != nil {
		return ErrInvalidAccount("the account already exists").Result()
	}
	// Construct a new account
	if cm.IsAdmin {
		newAcct = NewAdminUser(cm.PubKey, cm.Creator, cm.LegalEntityName, cm.LegalEntityType)
	} else {
		newAcct = NewOpUser(cm.PubKey, cm.Creator, cm.LegalEntityName, cm.LegalEntityType)
	}
	if err := CanCreateUserAccount(concreteCreator, newAcct); err != nil {
		return ErrUnauthorized(fmt.Sprintf("can't create account: %v", err)).Result()
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
	var newAcct *AppAccount
	// ensure proper message
	cm, ok := msg.(CreateAssetAccountMsg)
	if !ok {
		return ErrWrongMsgFormat("expected CreateAssetAccountMsg").Result()
	}
	// ensure creator exists
	// no need for type checking, CreateAssetAccount
	// validates types too.
	creator := h.accts.GetAccount(ctx, cm.Creator)
	if creator == nil {
		return ErrInvalidAccount("couldn't find creator").Result()
	}
	concreteCreator := creator.(*AppAccount)
	// ensure new account does not exist
	if h.accts.GetAccount(ctx, cm.PubKey.Address()) != nil {
		return ErrInvalidAccount("the account already exists").Result()
	}
	// Construct a new account
	newAcct, err := CreateAssetAccount(concreteCreator, cm.PubKey, nil)
	if err != nil {
		return ErrWrongSigner(fmt.Sprintf("couldn't create account: %v", err)).Result()
	}
	h.accts.SetAccount(ctx, newAcct)
	return sdk.Result{}
}

// Business logic

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

func getCHOperator(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address) (*AppAccount, sdk.Error) {
	operator, err := getCHUserAccount(ctx, accts, addr, false)
	if err != nil {
		return nil, err
	}
	if !operator.IsActive() {
		return nil, ErrWrongSigner("the operator account must be active")
	}
	if operator.IsAdmin() {
		return nil, ErrWrongSigner("admins cannot perform this operation")
	}
	return operator, nil
}

func getCHUserAccount(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address, isAdmin bool) (*AppAccount, sdk.Error) {
	if isAdmin {
		return getAdminWithEntityType(ctx, accts, addr, IsClearingHouse)
	}
	return getUserWithEntityType(ctx, accts, addr, IsClearingHouse)
}

func getAdminWithEntityType(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address,
	entityTypeCheck func(LegalEntity) bool) (*AppAccount, sdk.Error) {
	account, err := getUserWithEntityType(ctx, accts, addr, entityTypeCheck)
	if err != nil {
		return nil, err
	}
	if !account.IsAdmin() {
		return nil, ErrWrongSigner("must be admin user")
	}
	return account, nil
}

func getUserWithEntityType(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address,
	entityTypeCheck func(LegalEntity) bool) (*AppAccount, sdk.Error) {
	account, err := getAccountWithEntityType(ctx, accts, addr, entityTypeCheck)
	if err != nil {
		return nil, err
	}
	if !IsUser(account) {
		return nil, ErrWrongSigner("invalid account type")
	}
	return account, nil
}

func getAssetWithEntityType(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address,
	entityTypeCheck func(LegalEntity) bool) (*AppAccount, sdk.Error) {
	account, err := getAccountWithEntityType(ctx, accts, addr, entityTypeCheck)
	if err != nil {
		return nil, err
	}
	if !IsAsset(account) {
		return nil, ErrWrongSigner("invalid account type")
	}
	return account, nil
}

// Returns the account and verifies its type.
func getAccountWithEntityType(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address,
	entityTypeCheck func(LegalEntity) bool) (*AppAccount, sdk.Error) {
	rawAccount := accts.GetAccount(ctx, addr)
	if rawAccount == nil {
		return nil, ErrInvalidAccount("account does not exist")
	}
	account := rawAccount.(*AppAccount)
	if !entityTypeCheck(account) {
		return nil, ErrWrongSigner(account.GetLegalEntityType())
	}
	return account, nil
}
