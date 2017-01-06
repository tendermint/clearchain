package state

import (
	"encoding/json"

	bctypes "github.com/tendermint/basecoin/types"
	"github.com/tendermint/clearchain/types"
	"github.com/tendermint/go-common"
	"github.com/tendermint/go-events"
	tmsp "github.com/tendermint/tmsp/types"
)

func transfer(state *State, tx *types.TransferTx, isCheckTx bool) tmsp.Result {
	// // Validate basic structure
	if res := tx.ValidateBasic(); res.IsErr() {
		return res.PrependLog("in ValidateBasic()")
	}

	// Retrieve Sender's data
	user := state.GetUser(tx.Sender.Address)
	if user == nil {
		return tmsp.ErrBaseUnknownAddress.AppendLog("Sender's user is unknown")
	}
	entity := state.GetLegalEntity(user.EntityID)
	if entity == nil {
		return tmsp.ErrUnauthorized.AppendLog("User's does not belong to any LegalEntity")
	}

	// Get the accounts
	senderAccount := state.GetAccount(tx.Sender.AccountID)
	if senderAccount == nil {
		return tmsp.ErrBaseUnknownAddress.AppendLog("Sender's account is unknown")
	}
	recipientAccount := state.GetAccount(tx.Recipient.AccountID)
	if recipientAccount == nil {
		return tmsp.ErrBaseUnknownAddress.AppendLog("Unknown recipient address")
	}

	// Validate sender's Account
	if res := validateWalletSequence(senderAccount, tx.Sender); res.IsErr() {
		return res.PrependLog("in validateWalletSequence()")
	}

	// Generate byte-to-byte signature
	signBytes := tx.SignBytes(state.GetChainID())

	// Validate sender's permissions and signature
	if res := validateSender(senderAccount, entity, user, signBytes, tx); res.IsErr() {
		return res.PrependLog("in validateSender()")
	}
	// Validate counter signers
	if res := validateCounterSigners(state, senderAccount, entity, signBytes, tx); res.IsErr() {
		return res.PrependLog("in validateCounterSigners()")
	}

	// Apply changes
	applyChangesToInput(state, tx.Sender, senderAccount, isCheckTx)
	applyChangesToOutput(state, tx.Sender, tx.Recipient, recipientAccount, isCheckTx)

	return tmsp.OK

}

func createAccount(state *State, tx *types.CreateAccountTx, isCheckTx bool) tmsp.Result {
	// // Validate basic structure
	if res := tx.ValidateBasic(); res.IsErr() {
		return res.PrependLog("in ValidateBasic()")
	}

	// Retrieve user data
	user := state.GetUser(tx.Address)
	if user == nil {
		return tmsp.ErrBaseUnknownAddress.AppendLog("User is unknown")
	}
	entity := state.GetLegalEntity(user.EntityID)
	if entity == nil {
		return tmsp.ErrUnauthorized.AppendLog("User's does not belong to any LegalEntity")
	}
	// Validate permissions
	if !types.CanExecTx(user, tx) {
		return tmsp.ErrUnauthorized.AppendLog(common.Fmt(
			"User is not authorized to execute the Tx: %s", user.String()))
	}
	if !types.CanExecTx(entity, tx) {
		return tmsp.ErrUnauthorized.AppendLog(common.Fmt(
			"LegalEntity is not authorized to execute the Tx: %s", entity.String()))
	}
	// Generate byte-to-byte signature and validate the signature
	signBytes := tx.SignBytes(state.GetChainID())
	if !user.VerifySignature(signBytes, tx.Signature) {
		return tmsp.ErrBaseInvalidSignature.AppendLog("Verification failed, user's signature doesn't match")
	}

	// Create the new account
	if acc := state.GetAccount(tx.AccountID); acc != nil {
		return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Account already exists: %q", tx.AccountID))
	}
	// Get or create the accounts index
	accountIndex := GetOrMakeAccountIndex(state)
	if accountIndex.Has(tx.AccountID) {
		return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Account already exists in the account index: %q", tx.AccountID))
	}
	acc := types.NewAccount(tx.AccountID, entity.ID)
	accountIndex.Add(tx.AccountID)
	if !isCheckTx {
		state.SetAccount(acc.ID, acc)
		state.SetAccountIndex(accountIndex)
	}

	return tmsp.OK
}

func createLegalEntity(state *State, tx *types.CreateLegalEntityTx, isCheckTx bool) tmsp.Result {
	// // Validate basic structure
	if res := tx.ValidateBasic(); res.IsErr() {
		return res.PrependLog("in ValidateBasic()")
	}

	// Retrieve user data
	user := state.GetUser(tx.Address)
	if user == nil {
		return tmsp.ErrBaseUnknownAddress.AppendLog("User is unknown")
	}
	entity := state.GetLegalEntity(user.EntityID)
	if entity == nil {
		return tmsp.ErrUnauthorized.AppendLog("User's does not belong to any LegalEntity")
	}
	// Validate permissions
	if !types.CanExecTx(user, tx) {
		return tmsp.ErrUnauthorized.AppendLog(common.Fmt(
			"User is not authorized to execute the Tx: %s", user.String()))
	}
	if !types.CanExecTx(entity, tx) {
		return tmsp.ErrUnauthorized.AppendLog(common.Fmt(
			"LegalEntity is not authorized to execute the Tx: %s", entity.String()))
	}
	// Generate byte-to-byte signature and validate the signature
	signBytes := tx.SignBytes(state.GetChainID())
	if !user.VerifySignature(signBytes, tx.Signature) {
		return tmsp.ErrBaseInvalidSignature.AppendLog("Verification failed, user's signature doesn't match")
	}

	// Create new legal entity
	if ent := state.GetLegalEntity(tx.EntityID); ent != nil {
		return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("LegalEntity already exists: %q", tx.EntityID))
	}
	makeNewEntity(state, user, tx, isCheckTx)

	return tmsp.OK
}

func createUser(state *State, tx *types.CreateUserTx, isCheckTx bool) tmsp.Result {
	// // Validate basic structure
	if res := tx.ValidateBasic(); res.IsErr() {
		return res.PrependLog("in ValidateBasic()")
	}

	// Retrieve user data
	creator := state.GetUser(tx.Address)
	if creator == nil {
		return tmsp.ErrBaseUnknownAddress.AppendLog("User is unknown")
	}
	entity := state.GetLegalEntity(creator.EntityID)
	if entity == nil {
		return tmsp.ErrUnauthorized.AppendLog("User's does not belong to any LegalEntity")
	}

	// Validate permissions
	if !types.CanExecTx(creator, tx) {
		return tmsp.ErrUnauthorized.AppendLog(common.Fmt(
			"User is not authorized to execute the Tx: %s", creator.String()))
	}
	if !types.CanExecTx(entity, tx) {
		return tmsp.ErrUnauthorized.AppendLog(common.Fmt(
			"LegalEntity is not authorized to execute the Tx: %s", entity.String()))
	}
	// Generate byte-to-byte signature and validate the signature
	signBytes := tx.SignBytes(state.GetChainID())
	if !creator.VerifySignature(signBytes, tx.Signature) {
		return tmsp.ErrBaseInvalidSignature.AppendLog("Verification failed, user's signature doesn't match")
	}
	// Create new user
	if usr := state.GetUser(tx.PubKey.Address()); usr != nil {
		return tmsp.ErrBaseDuplicateAddress.AppendLog(common.Fmt("User already exists: %q", tx.PubKey.Address()))
	}
	makeNewUser(state, creator, tx, isCheckTx)

	return tmsp.OK
}

// ExecTx actually executes a Tx
func ExecTx(state *State, pgz *bctypes.Plugins, tx types.Tx,
	isCheckTx bool, evc events.Fireable) tmsp.Result {

	// Execute transaction
	switch tx := tx.(type) {
	case *types.TransferTx:
		return transfer(state, tx, isCheckTx)

	case *types.CreateAccountTx:
		return createAccount(state, tx, isCheckTx)

	case *types.CreateLegalEntityTx:
		return createLegalEntity(state, tx, isCheckTx)

	case *types.CreateUserTx:
		return createUser(state, tx, isCheckTx)

	default:
		return tmsp.ErrBaseEncodingError.SetLog("Unknown tx type")
	}
}

func accountQuery(state *State, tx *types.AccountQueryTx) tmsp.Result {
	// Validate basic
	if res := tx.ValidateBasic(); res.IsErr() {
		return res.PrependLog("in ValidateBasic()")
	}

	user := state.GetUser(tx.Address)
	if user == nil {
		return tmsp.ErrBaseUnknownAddress.AppendLog(common.Fmt("Address is unknown: %v", tx.Address))
	}
	accounts := make([]*types.Account, len(tx.Accounts))
	for i, accountID := range tx.Accounts {
		account := state.GetAccount(accountID)
		if account == nil {
			return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid account_id: %q", accountID))
		}
		accounts[i] = account
	}

	// Generate byte-to-byte signature
	signBytes := tx.SignBytes(state.GetChainID())
	if !user.VerifySignature(signBytes, tx.Signature) {
		return tmsp.ErrUnauthorized.AppendLog("Verification failed, signature doesn't match")
	}
	data, err := json.Marshal(struct {
		Account []*types.Account `json:"accounts"`
	}{accounts})
	if err != nil {
		return tmsp.ErrInternalError.AppendLog(common.Fmt("Couldn't make the response: %v", err))
	}
	return tmsp.OK.SetData(data)
}

func accountIndexQuery(state *State, tx *types.AccountIndexQueryTx) tmsp.Result {
	// Validate basic
	if res := tx.ValidateBasic(); res.IsErr() {
		return res.PrependLog("in ValidateBasic()")
	}

	user := state.GetUser(tx.Address)
	if user == nil {
		return tmsp.ErrBaseUnknownAddress.AppendLog(common.Fmt("Address is unknown: %v", tx.Address))
	}

	// Check that the account index exists
	accountIndex := state.GetAccountIndex()
	if accountIndex == nil {
		return tmsp.ErrInternalError.AppendLog("AccountIndex has not yet been initialized")
	}

	// Generate byte-to-byte signature
	signBytes := tx.SignBytes(state.GetChainID())
	if !user.VerifySignature(signBytes, tx.Signature) {
		return tmsp.ErrUnauthorized.AppendLog("Verification failed, signature doesn't match")
	}
	data, err := json.Marshal(accountIndex)
	if err != nil {
		return tmsp.ErrInternalError.AppendLog(common.Fmt("Couldn't make the response: %v", err))
	}
	return tmsp.OK.SetData(data)
}

// ExecQueryTx handles queries.
func ExecQueryTx(state *State, tx types.Tx) tmsp.Result {
	
	// Execute transaction
	switch tx := tx.(type) {
	case *types.AccountQueryTx:
		return accountQuery(state, tx)

	case *types.AccountIndexQueryTx:
		return accountIndexQuery(state, tx)

	case *types.LegalEntityQueryTx:
		return legalEntityQuery(state, tx)

	case *types.LegalEntityIndexQueryTx:
		return legalEntityIndexQueryTx(state, tx)

	default:
		return tmsp.ErrBaseEncodingError.SetLog("Unknown tx type")
	}
}

//--------------------------------------------------------------------------------

func validateWalletSequence(acc *types.Account, in types.TxTransferSender) tmsp.Result {
	wal := acc.GetWallet(in.Currency)
	// Wallet does not exist, Sequence must be 1
	if wal == nil {
		if in.Sequence != 1 {
			return tmsp.ErrBaseInvalidSequence.AppendLog(common.Fmt("Invalid sequence: got: %v, want: 1", in.Sequence))
		}
		return tmsp.OK
	}
	if in.Sequence != wal.Sequence+1 {
		return tmsp.ErrBaseInvalidSequence.AppendLog(common.Fmt("Invalid sequence: got: %v, want: %v", in.Sequence, wal.Sequence+1))
	}
	return tmsp.OK
}

func validateSender(acc *types.Account, entity *types.LegalEntity, u *types.User, signBytes []byte, tx *types.TransferTx) tmsp.Result {
	if res := validatePermissions(u, entity, acc, tx); res.IsErr() {
		return res
	}
	if !u.VerifySignature(signBytes, tx.Sender.Signature) {
		return tmsp.ErrBaseInvalidSignature.AppendLog("Verification failed, sender's signature doesn't match")
	}
	return tmsp.OK
}

// Validate countersignatures
func validateCounterSigners(state types.UserGetter, acc *types.Account, entity *types.LegalEntity, signBytes []byte, tx *types.TransferTx) tmsp.Result {
	var users = make(map[string]bool)

	// Make sure users are not duplicated
	users[string(tx.Sender.Address)] = true

	for _, in := range tx.CounterSigners {
		// Users must not be duplicated either
		if _, ok := users[string(in.Address)]; ok {
			return tmsp.ErrBaseDuplicateAddress
		}
		users[string(in.Address)] = true

		// User must exist
		user := state.GetUser(in.Address)
		if user == nil {
			return tmsp.ErrBaseUnknownAddress
		}

		// Validate the permissions
		if res := validatePermissions(user, entity, acc, tx); res.IsErr() {
			return res
		}
		// Verify the signature
		if !user.VerifySignature(signBytes, in.Signature) {
			return tmsp.ErrBaseInvalidSignature.AppendLog(common.Fmt("Verification failed, countersigner's signature doesn't match, user: %s", user))
		}
	}

	return tmsp.OK
}

func validatePermissions(u *types.User, e *types.LegalEntity, a *types.Account, tx types.Tx) tmsp.Result {
	// Verify user belongs to the legal entity
	if !a.BelongsTo(u.EntityID) {
		return tmsp.ErrUnauthorized.AppendLog(common.Fmt(
			"Access forbidden for user %s to account %s", u.Name, a.String()))
	}
	// Valdate permissions
	if !types.CanExecTx(u, tx) {
		return tmsp.ErrUnauthorized.AppendLog(common.Fmt(
			"User is not authorized to execute the Tx: %s", u.String()))
	}
	if !types.CanExecTx(e, tx) {
		return tmsp.ErrUnauthorized.AppendLog(common.Fmt(
			"LegalEntity is not authorized to execute the Tx: %s", e.String()))
	}
	return tmsp.OK
}

// Apply changes to inputs
func applyChangesToInput(state types.AccountSetter, in types.TxTransferSender, acc *types.Account, isCheckTx bool) {
	wal := acc.GetWallet(in.Currency)
	if wal == nil {
		acc.Wallets = append(acc.Wallets, types.Wallet{
			Currency: in.Currency,
			Balance:  -in.Amount,
			Sequence: 1})
	} else {
		wal.Balance -= in.Amount
		wal.Sequence++
	}
	if !isCheckTx {
		state.SetAccount(in.AccountID, acc)
	}
}

// Apply changes to outputs
func applyChangesToOutput(state types.AccountSetter, in types.TxTransferSender, out types.TxTransferRecipient, acc *types.Account, isCheckTx bool) {
	wal := acc.GetWallet(in.Currency)
	if wal == nil {
		acc.Wallets = append(acc.Wallets, types.Wallet{
			Currency: in.Currency,
			Balance:  in.Amount,
			Sequence: 1})

	} else {
		wal.Balance += in.Amount
		wal.Sequence++
	}
	if !isCheckTx {
		state.SetAccount(out.AccountID, acc)
	}
}

func makeNewEntity(state types.LegalEntitySetter, user *types.User, tx *types.CreateLegalEntityTx, isCheckTx bool) {
	ent := types.NewLegalEntityByType(tx.Type, tx.EntityID, tx.Name, user.PubKey.Address())
	if ent == nil {
		common.PanicSanity(common.Fmt("Unexpected TxType: %x", tx.Type))
	}
	if !isCheckTx {
		state.SetLegalEntity(ent.ID, ent)
	}
}

func makeNewUser(state types.UserSetter, creator *types.User, tx *types.CreateUserTx, isCheckTx bool) {
	perms := creator.Permissions
	if !tx.CanCreate {
		perms = perms.Clear(types.PermCreateUserTx.Add(types.PermCreateLegalEntityTx))
	}
	user := types.NewUser(tx.PubKey, tx.Name, creator.EntityID, perms)
	if user == nil {
		common.PanicSanity(common.Fmt("Unexpected nil User"))
	}
	if !isCheckTx {
		state.SetUser(tx.PubKey.Address(), user)
	}
}

func GetOrMakeAccountIndex(state types.AccountIndexGetter) *types.AccountIndex {
	if index := state.GetAccountIndex(); index != nil {
		return index
	}
	return types.NewAccountIndex()
}

func legalEntityQuery(state *State, tx *types.LegalEntityQueryTx) tmsp.Result {
	// Validate basic
	if res := tx.ValidateBasic(); res.IsErr() {
		return res.PrependLog("in ValidateBasic()")
	}

	user := state.GetUser(tx.Address)
	if user == nil {
		return tmsp.ErrBaseUnknownAddress.AppendLog(common.Fmt("Address is unknown: %v", tx.Address))
	}
	legalEntities := make([]*types.LegalEntity, len(tx.Ids))
	for i, id := range tx.Ids {
		legalEntity := state.GetLegalEntity(id)
		if legalEntity == nil {
			return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid legalEntity id: %q", id))
		}
		legalEntities[i] = legalEntity
	}

	// Generate byte-to-byte signature
	signBytes := tx.SignBytes(state.GetChainID())
	if !user.VerifySignature(signBytes, tx.Signature) {
		return tmsp.ErrUnauthorized.AppendLog("Verification failed, signature doesn't match")
	}
	data, err := json.Marshal(struct {
		LegalEntities []*types.LegalEntity `json:"legalEntities"`
	}{legalEntities})
	if err != nil {
		return tmsp.ErrInternalError.AppendLog(common.Fmt("Couldn't make the response: %v", err))
	}
	return tmsp.OK.SetData(data)
}

func legalEntityIndexQueryTx(state *State, tx *types.LegalEntityIndexQueryTx) tmsp.Result {
	// Validate basic
	if res := tx.ValidateBasic(); res.IsErr() {
		return res.PrependLog("in ValidateBasic()")
	}

	user := state.GetUser(tx.Address)
	if user == nil {
		return tmsp.ErrBaseUnknownAddress.AppendLog(common.Fmt("Address is unknown: %v", tx.Address))
	}

	// Check that the account index exists
	legalEntities := state.GetLegalEntityIndex()
	if legalEntities == nil {
		return tmsp.ErrInternalError.AppendLog("LegalEntities has not yet been initialized")
	}

	// Generate byte-to-byte signature
	signBytes := tx.SignBytes(state.GetChainID())
	if !user.VerifySignature(signBytes, tx.Signature) {
		return tmsp.ErrUnauthorized.AppendLog("Verification failed, signature doesn't match")
	}
	data, err := json.Marshal(legalEntities)
	if err != nil {
		return tmsp.ErrInternalError.AppendLog(common.Fmt("Couldn't make the response: %v", err))
	}
	return tmsp.OK.SetData(data)
}
