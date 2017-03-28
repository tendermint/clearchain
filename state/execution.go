package state

import (
	"encoding/json"
	
	abci "github.com/tendermint/abci/types"
	bctypes "github.com/tendermint/basecoin/types"
	"github.com/tendermint/clearchain/types"
	"github.com/tendermint/go-common"
	"github.com/tendermint/go-events"
)

func transfer(state *State, tx *types.TransferTx, isCheckTx bool) abci.Result {
	// // Validate basic structure
	if res := tx.ValidateBasic(); res.IsErr() {
		return res.PrependLog("in ValidateBasic()")
	}

	// Retrieve Committer's data
	user := state.GetUser(tx.Committer.Address)
	if user == nil {
		return abci.ErrBaseUnknownAddress.AppendLog("Sender's user is unknown")
	}
	committerEntity := state.GetLegalEntity(user.EntityID)
	if committerEntity == nil {
		return abci.ErrUnauthorized.AppendLog("User's does not belong to any LegalEntity")
	}

	// Get the accounts
	senderAccount := state.GetAccount(tx.Sender.AccountID)
	if senderAccount == nil {
		return abci.ErrBaseUnknownAddress.AppendLog("Sender's account is unknown")
	}
	recipientAccount := state.GetAccount(tx.Recipient.AccountID)
	if recipientAccount == nil {
		return abci.ErrBaseUnknownAddress.AppendLog("Unknown recipient address")
	}

	// Get legal entities
	senderEntity := state.GetLegalEntity(senderAccount.EntityID)
	if committerEntity == nil {
		return abci.ErrUnauthorized.AppendLog("Sender's account does not belong to any LegalEntity")
	}
	recipientEntity := state.GetLegalEntity(recipientAccount.EntityID)
	if recipientEntity == nil {
		return abci.ErrUnauthorized.AppendLog("Recipient's account does not belong to any LegalEntity")
	}

	// Validate sender's Account
	if res := validateWalletSequence(senderAccount, tx.Sender); res.IsErr() {
		return res.PrependLog("in validateWalletSequence()")
	}

	// Generate byte-to-byte signature
	signBytes := tx.SignBytes(state.GetChainID())

	// Validate committer's permissions and signature
	if !user.VerifySignature(signBytes, tx.Committer.Signature) {
		return abci.ErrBaseInvalidSignature.AppendLog("sender's signature doesn't match")
	}
	if res := validateExecPermissions(user, committerEntity, tx); res.IsErr() {
		return res
	}
	if res := validateCommitter(user, committerEntity, senderEntity, recipientEntity, signBytes, tx); res.IsErr() {
		return res.PrependLog("in validateCommitter()")
	}

	// Validate counter signers
	if res := validateCounterSigners(state, committerEntity, tx); res.IsErr() {
		return res.PrependLog("in validateCounterSigners()")
	}

	// Apply changes
	applyChangesToInput(state, tx.Sender, senderAccount, isCheckTx)
	applyChangesToOutput(state, tx.Sender, tx.Recipient, recipientAccount, isCheckTx)

	return abci.OK

}

func createAccount(state *State, tx *types.CreateAccountTx, isCheckTx bool) abci.Result {
	// // Validate basic structure
	if res := tx.ValidateBasic(); res.IsErr() {
		return res.PrependLog("in ValidateBasic()")
	}

	// Retrieve user data
	user := state.GetUser(tx.Address)
	if user == nil {
		return abci.ErrBaseUnknownAddress.AppendLog("User is unknown")
	}
	entity := state.GetLegalEntity(user.EntityID)
	if entity == nil {
		return abci.ErrUnauthorized.AppendLog("User's does not belong to any LegalEntity")
	}
	// Validate permissions
	if !types.CanExecTx(user, tx) {
		return abci.ErrUnauthorized.AppendLog(common.Fmt(
			"User is not authorized to execute the Tx: %s", user.String()))
	}
	if !types.CanExecTx(entity, tx) {
		return abci.ErrUnauthorized.AppendLog(common.Fmt(
			"LegalEntity is not authorized to execute the Tx: %s", entity.String()))
	}
	// Generate byte-to-byte signature and validate the signature
	signBytes := tx.SignBytes(state.GetChainID())
	if !user.VerifySignature(signBytes, tx.Signature) {
		return abci.ErrBaseInvalidSignature.AppendLog("user's signature doesn't match")
	}

	// Create the new account
	if acc := state.GetAccount(tx.AccountID); acc != nil {
		return abci.ErrBaseInvalidInput.AppendLog(common.Fmt("Account already exists: %q", tx.AccountID))
	}
	// Get or create the accounts index
	if !isCheckTx {
		acc := types.NewAccount(tx.AccountID, entity.ID)
		state.SetAccount(acc.ID, acc)
		return SetAccountInIndex(state, *acc)
	}

	return abci.OK
}

func createLegalEntity(state *State, tx *types.CreateLegalEntityTx, isCheckTx bool) abci.Result {
	// // Validate basic structure
	if res := tx.ValidateBasic(); res.IsErr() {
		return res.PrependLog("in ValidateBasic()")
	}

	// Retrieve user data
	user := state.GetUser(tx.Address)
	if user == nil {
		return abci.ErrBaseUnknownAddress.AppendLog("User is unknown")
	}
	entity := state.GetLegalEntity(user.EntityID)
	if entity == nil {
		return abci.ErrUnauthorized.AppendLog("User's does not belong to any LegalEntity")
	}
	// Validate permissions
	if !types.CanExecTx(user, tx) {
		return abci.ErrUnauthorized.AppendLog(common.Fmt(
			"User is not authorized to execute the Tx: %s", user.String()))
	}
	if !types.CanExecTx(entity, tx) {
		return abci.ErrUnauthorized.AppendLog(common.Fmt(
			"LegalEntity is not authorized to execute the Tx: %s", entity.String()))
	}
	// Generate byte-to-byte signature and validate the signature
	signBytes := tx.SignBytes(state.GetChainID())
	if !user.VerifySignature(signBytes, tx.Signature) {
		return abci.ErrBaseInvalidSignature.AppendLog("user's signature doesn't match")
	}

	// Create new legal entity
	if ent := state.GetLegalEntity(tx.EntityID); ent != nil {
		return abci.ErrBaseInvalidInput.AppendLog(common.Fmt("LegalEntity already exists: %q", tx.EntityID))
	}
	if !isCheckTx {
		legalEntity := types.NewLegalEntityByType(tx.Type, tx.EntityID, tx.Name, user.PubKey.Address(), tx.ParentID)
		state.SetLegalEntity(legalEntity.ID, legalEntity)
		return SetLegalEntityInIndex(state, legalEntity)
	}

	return abci.OK
}

func createUser(state *State, tx *types.CreateUserTx, isCheckTx bool) abci.Result {
	// // Validate basic structure
	if res := tx.ValidateBasic(); res.IsErr() {
		return res.PrependLog("in ValidateBasic()")
	}

	// Retrieve user data
	creator := state.GetUser(tx.Address)
	if creator == nil {
		return abci.ErrBaseUnknownAddress.AppendLog("User is unknown")
	}
	entity := state.GetLegalEntity(creator.EntityID)
	if entity == nil {
		return abci.ErrUnauthorized.AppendLog("User's does not belong to any LegalEntity")
	}

	// Validate permissions
	if !types.CanExecTx(creator, tx) {
		return abci.ErrUnauthorized.AppendLog(common.Fmt(
			"User is not authorized to execute the Tx: %s", creator.String()))
	}
	if !types.CanExecTx(entity, tx) {
		return abci.ErrUnauthorized.AppendLog(common.Fmt(
			"LegalEntity is not authorized to execute the Tx: %s", entity.String()))
	}
	// Generate byte-to-byte signature and validate the signature
	signBytes := tx.SignBytes(state.GetChainID())
	if !creator.VerifySignature(signBytes, tx.Signature) {
		return abci.ErrBaseInvalidSignature.AppendLog("user's signature doesn't match")
	}
	// Create new user
	if usr := state.GetUser(tx.PubKey.Address()); usr != nil {
		return abci.ErrBaseDuplicateAddress.AppendLog(common.Fmt("User already exists: %q", tx.PubKey.Address()))
	}
	makeNewUser(state, creator, tx, isCheckTx)

	return abci.OK
}

// ExecTx actually executes a Tx
func ExecTx(state *State, pgz *bctypes.Plugins, tx types.Tx,
	isCheckTx bool, evc events.Fireable) abci.Result {

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
		return abci.ErrBaseEncodingError.SetLog("Unknown tx type")
	}
}

 func accountQuery(state *State, accountID string) (res abci.ResponseQuery) {
 	account := state.GetAccount(accountID)
 	if account == nil {
 		res.Code = abci.CodeType_BaseInvalidInput
		res.Log = common.Fmt("Invalid account_id: %q", accountID)
		return
 	}	
 	
 	data, err := json.Marshal(types.AccountsReturned{Account: []*types.Account{account}})
 	if err != nil {
	 	res.Code = abci.CodeType_InternalError
		res.Log = common.Fmt("Couldn't make the response: %v", err)
		return 		
 	}
 	
 	res.Code = abci.CodeType_OK
 	res.Value = data
 	return 
 }

 func accountIndexQuery(state *State) (res abci.ResponseQuery) {
 	
 	// Check that the account index exists
 	accountIndex := state.GetAccountIndex()
 	if accountIndex == nil {
 		res.Code = abci.CodeType_InternalError
		res.Log = "AccountIndex has not yet been initialized"
		return 	 		
 	}
 	
 	data, err := json.Marshal(accountIndex)
 	if err != nil {
 		res.Code = abci.CodeType_InternalError
		res.Log = common.Fmt("Couldn't make the response: %v", err)
		return 	  		
 	}
 	
 	res.Code = abci.CodeType_OK
 	res.Value = data
 	return 
 }

// ExecQuery handles queries.
func ExecQuery(state *State, resource string, object string) abci.ResponseQuery {

	 switch  {
		 case resource == "account" && len(object) > 0 :
		 	return accountQuery(state, object)
	
		 case resource == "account" && len(object) == 0 :
		 	return accountIndexQuery(state)
	
		 case resource == "legal_entity" && len(object) > 0 :
		 	return legalEntityQuery(state, object)
	
		 case resource == "legal_entity" && len(object) == 0 :
		 	return legalEntityIndexQueryTx(state)
	
		 default:			
			return  abci.ResponseQuery {
				Code : abci.CodeType_BaseEncodingError,
				Log : common.Fmt("Unknown resource and object: %v/%v", resource, object),
			}		 	
	 }
	
}

//--------------------------------------------------------------------------------

func validateWalletSequence(acc *types.Account, in types.TxTransferSender) abci.Result {
	wal := acc.GetWallet(in.Currency)
	// Wallet does not exist, Sequence must be 1
	if wal == nil {
		if in.Sequence != 1 {
			return abci.ErrBaseInvalidSequence.AppendLog(common.Fmt("Invalid sequence: got: %v, want: 1", in.Sequence))
		}
		return abci.OK
	}
	if in.Sequence != wal.Sequence+1 {
		return abci.ErrBaseInvalidSequence.AppendLog(common.Fmt("Invalid sequence: got: %v, want: %v", in.Sequence, wal.Sequence+1))
	}
	return abci.OK
}

func validateCommitter(u *types.User, committerEntity, senderEntity, recipientEntity *types.LegalEntity, signBytes []byte, tx *types.TransferTx) abci.Result {
	// TODO: apply business rules
	return abci.OK
}

// Validate countersignatures
func validateCounterSigners(state *State, entity *types.LegalEntity, tx *types.TransferTx) abci.Result {
	var users = make(map[string]bool)

	// Make sure users are not duplicated
	users[string(tx.Committer.Address)] = true

	for _, in := range tx.CounterSigners {
		// Users must not be duplicated either
		if _, ok := users[string(in.Address)]; ok {
			return abci.ErrBaseDuplicateAddress
		}
		users[string(in.Address)] = true

		// User must exist
		user := state.GetUser(in.Address)
		if user == nil {
			return abci.ErrBaseUnknownAddress
		}

		// Validate the permissions
		if res := validateExecPermissions(user, entity, tx); res.IsErr() {
			return res
		}
		// Verify the signature
		if !user.VerifySignature(in.SignBytes(state.GetChainID()), in.Signature) {
			return abci.ErrBaseInvalidSignature.AppendLog(common.Fmt("countersigner's signature doesn't match, user: %s", user))
		}
	}

	return abci.OK
}

func validateExecPermissions(u *types.User, e *types.LegalEntity, tx types.Tx) abci.Result {
	// Valdate exec permissions
	if !types.CanExecTx(u, tx) {
		return abci.ErrUnauthorized.AppendLog(common.Fmt(
			"User is not authorized to execute the Tx: %s", u.String()))
	}
	if !types.CanExecTx(e, tx) {
		return abci.ErrUnauthorized.AppendLog(common.Fmt(
			"LegalEntity is not authorized to execute the Tx: %s", e.String()))
	}
	return abci.OK
}

// Apply changes to inputs
func applyChangesToInput(state types.AccountSetter, in types.TxTransferSender, account *types.Account, isCheckTx bool) {
	applyChanges(account, in.Currency, in.Amount, false)

	if !isCheckTx {
		state.SetAccount(account.ID, account)
	}
}

// Apply changes to outputs
func applyChangesToOutput(state types.AccountSetter, in types.TxTransferSender, out types.TxTransferRecipient, account *types.Account, isCheckTx bool) {
	applyChanges(account, in.Currency, in.Amount, true)

	if !isCheckTx {
		state.SetAccount(account.ID, account)
	}
}

func applyChanges(account *types.Account, currency string, amount int64, isBuy bool) {

	wal := account.GetWallet(currency)

	if wal == nil {
		wal = &types.Wallet{Currency: currency}
	}

	if isBuy {
		wal.Balance += amount
	} else {
		wal.Balance += -amount
	}

	wal.Sequence++

	account.SetWallet(*wal)
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

//Returns existing AccountIndex from store or creates new empty one
func GetOrMakeAccountIndex(state types.AccountIndexGetter) *types.AccountIndex {
	if index := state.GetAccountIndex(); index != nil {
		return index
	}
	return types.NewAccountIndex()
}

//Sets Account in AccountIndex in store
func SetAccountInIndex(state *State, account types.Account) abci.Result {
	accountIndex := GetOrMakeAccountIndex(state)
	if accountIndex.Has(account.ID) {
		return abci.ErrBaseInvalidInput.AppendLog(common.Fmt("Account already exists in the account index: %q", account.ID))
	}
	accountIndex.Add(account.ID)
	state.SetAccountIndex(accountIndex)
	return abci.OK
}

//Sets LegalEntity in LegalEntityIndex in store
func SetLegalEntityInIndex(state *State, legalEntity *types.LegalEntity) abci.Result {
	legalEntities := state.GetLegalEntityIndex()

	if legalEntities == nil {
		legalEntities = &types.LegalEntityIndex{Ids: []string{}}
	}

	if legalEntities.Has(legalEntity.ID) {
		return abci.ErrBaseInvalidInput.AppendLog(common.Fmt("LegalEntity already exists in the LegalEntity index: %q", legalEntity.ID))
	}
	legalEntities.Add(legalEntity.ID)

	state.SetLegalEntityIndex(legalEntities)

	return abci.OK
}

 func legalEntityQuery(state *State, entityID string)  (res abci.ResponseQuery) {
 	
 	legalEntity := state.GetLegalEntity(entityID)
 	if legalEntity == nil {
 		res.Code = abci.CodeType_BaseInvalidInput
		res.Log = common.Fmt("Invalid legalEntity id: %q", entityID)
		return 			
 	}
 	data, err := json.Marshal(types.LegalEntitiesReturned{LegalEntities: []*types.LegalEntity{legalEntity}})
 	if err != nil {
 		res.Code = abci.CodeType_InternalError
		res.Log = common.Fmt("Couldn't make the response: %v", err)
		return 	  		
 	}
 	
 	res.Code = abci.CodeType_OK
 	res.Value = data
 	return 
 }

 func legalEntityIndexQueryTx(state *State) (res abci.ResponseQuery) {
 	
 	// Check that the account index exists
 	legalEntities := state.GetLegalEntityIndex()
 	if legalEntities == nil {
 		res.Code = abci.CodeType_InternalError
		res.Log = "LegalEntities has not yet been initialized"
		return 			
 	}
 	
 	data, err := json.Marshal(legalEntities)
 	if err != nil {
 		res.Code = abci.CodeType_InternalError
		res.Log = common.Fmt("Couldn't make the response: %v", err)
		return 	  		
 	}
 	
 	res.Code = abci.CodeType_OK
 	res.Value = data
 	return 
 }
