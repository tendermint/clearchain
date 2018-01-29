package types

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

func RegisterRoutes(r baseapp.Router, accts sdk.AccountMapper) {
	r.AddRoute(DepositType, DepositMsgHandler(accts))
	r.AddRoute(SettlementType, SettleMsgHandler(accts))
	r.AddRoute(WithdrawType, WithDrawMsgHandler(accts))
	r.AddRoute(CreateAccountType, CreateAccountMsgHandler(accts))
}

/*

Sender == Custodian
Rec == Member
*/
func DepositMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return depositMsgHandler{accts}.Do
}

type depositMsgHandler struct {
	accts sdk.AccountMapper
}

func (d depositMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	// TODO: ensure auth actually checks the sigs

	// ensure proper message
	dm, ok := msg.(DepositMsg)
	if !ok {
		return sdk.ErrTxParse("Expected DepositMsg").Result()
	}

	// ensure proper types
	sender, err := getAccountWithType(ctx, d.accts, dm.Sender, IsCustodian)
	if err != nil {
		return err.Result()
	}
	rcpt, err := getAccountWithType(ctx, d.accts, dm.Recipient, IsMember)
	if err != nil {
		return err.Result()
	}

	err = moveMoney(d.accts, ctx, sender, rcpt, dm.Amount, false, true)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

/*
Sender == CH
Rec == member
*/
func SettleMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return settleMsgHandler{accts}.Do
}

type settleMsgHandler struct {
	accts sdk.AccountMapper
}

func (sh settleMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {

	// ensure proper message
	sm, ok := msg.(SettleMsg)
	if !ok {
		return sdk.ErrTxParse("Expected SettleMsg").Result()
	}

	// ensure proper types
	sender, err := getAccountWithType(ctx, sh.accts, sm.Sender, IsClearingHouse)
	if err != nil {
		return err.Result()
	}
	rcpt, err := getAccountWithType(ctx, sh.accts, sm.Recipient, IsMember)
	if err != nil {
		return err.Result()
	}

	err = moveMoney(sh.accts, ctx, sender, rcpt, sm.Amount, false, true)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}

}

/*
Sender == member
Reci == custodian
Operator == CH

*/
func WithDrawMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return withdrawMsgHandler{accts}.Do
}

type withdrawMsgHandler struct {
	accts sdk.AccountMapper
}

func (wh withdrawMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	// ensure proper message
	wm, ok := msg.(WithdrawMsg)
	if !ok {
		return sdk.ErrTxParse("Expected WithdrawMsg").Result()
	}

	// ensure proper types
	sender, err := getAccountWithType(ctx, wh.accts, wm.Sender, IsMember)
	if err != nil {
		return err.Result()
	}
	rcpt, err := getAccountWithType(ctx, wh.accts, wm.Recipient, IsCustodian)
	if err != nil {
		return err.Result()
	}
	_, err = getAccountWithType(ctx, wh.accts, wm.Operator, IsClearingHouse)
	if err != nil {
		return err.Result()
	}

	err = moveMoney(wh.accts, ctx, sender, rcpt, wm.Amount, true, false)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}

}

/*
Creates a new account
*/
func CreateAccountMsgHandler(accts sdk.AccountMapper) sdk.Handler {
	return createAccountMsgHandler{accts}.Do
}

type createAccountMsgHandler struct {
	accts sdk.AccountMapper
}

func (h createAccountMsgHandler) Do(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	// ensure proper message
	cm, ok := msg.(CreateAccountMsg)
	if !ok {
		return sdk.ErrTxParse("Expected CreateAccountMsg").Result()
	}

	// ensure proper types
	creator, err := getAccountWithType(ctx, h.accts, cm.Creator, IsClearingHouse)
	if err != nil {
		return err.Result()
	}
	// A clearing house account is allowed to create any kind of accounts,
	// including clearing house, custodian, and members accounts.
	// TODO: clarify business rules and ensure this is desired
	if rawAccount := h.accts.GetAccount(ctx, cm.PubKey.Address()); rawAccount != nil {
		return ErrInvalidAccount("the account already exists").Result()
	}

	acct := createAccount(creator.GetAddress(), cm.PubKey, cm.AccountType)
	h.accts.SetAccount(ctx, acct)

	return sdk.Result{}
}

//*********************************** helper methods *********************************************

func moveMoney(accts sdk.AccountMapper, ctx sdk.Context, sender *AppAccount, recipient *AppAccount,
	amount sdk.Coin, senderMustBePositive bool, recipientMustBePositive bool) sdk.Error {

	// now make the transfer
	transfer := sdk.Coins{amount}
	sender.Coins = sender.Coins.Minus(transfer)
	if senderMustBePositive && !sender.Coins.IsNotNegative() {
		return sdk.ErrInsufficientFunds("sender balance negative")
	}
	recipient.Coins = recipient.Coins.Plus(transfer)
	if recipientMustBePositive && !recipient.Coins.IsNotNegative() {
		return sdk.ErrInsufficientFunds("recipient balance negative")
	}

	// and save the result
	accts.SetAccount(ctx, sender)
	accts.SetAccount(ctx, recipient)
	return nil
}

func getAccountWithType(ctx sdk.Context, accts sdk.AccountMapper, addr crypto.Address,
	typeCheck func(*AppAccount) bool) (*AppAccount, sdk.Error) {

	rawAccount := accts.GetAccount(ctx, addr)
	if rawAccount == nil {
		return nil, sdk.ErrUnrecognizedAddress(addr)
	}
	account := rawAccount.(*AppAccount)
	if !typeCheck(account) {
		return nil, ErrWrongSigner(account.Type)
	}

	return account, nil
}

func createAccount(creator crypto.Address, newAccPubKey crypto.PubKey, typ string) *AppAccount {
	acct := new(AppAccount)
	acct.SetAddress(newAccPubKey.Address())
	acct.SetPubKey(newAccPubKey)
	acct.SetCoins(nil)
	acct.Type = typ
	// TODO:
	// acct.SetCreator(creator)
	return acct
}
