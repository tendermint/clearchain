package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	crypto "github.com/tendermint/go-crypto"
)

// EntityType string identifiers
const (
	AccountUser  = "user"
	AccountAsset = "asset"
)

// ensure AppAccount implements the sdk.Account interface
var (
	_ sdk.Account = (*AppAccount)(nil)
	_ LegalEntity = (*AppAccount)(nil)
	_ UserAccount = (*AppAccount)(nil)
)

// UserAccount is the interface that wraps the basic
// accessor methods to set and get user accounts attributes.
type UserAccount interface {
	GetAccountType() string
	IsAdmin() bool
	IsActive() bool
}

// AppAccount defines the properties of an AppAccount.
type AppAccount struct {
	auth.BaseAccount
	BaseLegalEntity
	Creator     crypto.Address
	AccountType string
	Active      bool
	Admin       bool
}

// NewAppAccount constructs a new account instance.
func newAppAccount(pub crypto.PubKey, cash sdk.Coins, creator crypto.Address, typ string,
	isActive bool, isAdmin bool, entityName, entityType string) *AppAccount {
	acct := new(AppAccount)
	acct.SetAddress(pub.Address())
	acct.SetPubKey(pub)
	acct.SetCoins(cash)
	acct.SetCreator(creator)
	acct.EntityName = entityName
	acct.EntityType = entityType
	acct.AccountType = typ
	acct.Active = isActive
	acct.Admin = isAdmin
	return acct
}

// NewOpUser constructs a new account instance, setting cash to nil.
func NewOpUser(pub crypto.PubKey, creator crypto.Address, entityName, entityType string) *AppAccount {
	return newAppAccount(pub, nil, creator, AccountUser, true, false, entityName, entityType)
}

// NewAdminUser constructs a new account instance, setting cash to nil.
func NewAdminUser(pub crypto.PubKey, creator crypto.Address, entityName, entityType string) *AppAccount {
	return newAppAccount(pub, nil, creator, AccountUser, true, true, entityName, entityType)
}

// NewAssetAccount constructs a new account instance.
func NewAssetAccount(pub crypto.PubKey, cash sdk.Coins, creator crypto.Address, entityName, entityType string) *AppAccount {
	return newAppAccount(pub, cash, creator, AccountAsset, true, false, entityName, entityType)
}

// GetCreator returns account's creator.
func GetCreator(a *AppAccount) crypto.Address {
	return a.Creator
}

// SetCreator modifies account's creator.
func (a *AppAccount) SetCreator(creator crypto.Address) {
	a.Creator = creator
}

// GetAccountType returns the account type.
func (a *AppAccount) GetAccountType() string {
	return a.AccountType
}

// IsActive returns true if the account is active; false otherwise.
func (a *AppAccount) IsActive() bool {
	return a.Active
}

// IsAdmin returns true if the account is admin; false otherwise.
func (a *AppAccount) IsAdmin() bool {
	return a.Admin
}

// IsUser returns true if the account holds user data; false otherwise.
func IsUser(a UserAccount) bool {
	return a.GetAccountType() == AccountUser
}

// IsAsset returns true if the account holds assets; false otherwise.
func IsAsset(a UserAccount) bool {
	return a.GetAccountType() == AccountAsset
}

// IsAdminUser returns true if the account is an
// admin user account of its legal entity;
// false otherwise.
func IsAdminUser(a UserAccount) bool {
	return IsUser(a) && a.IsAdmin()
}

// CanCreateUserAccount returns nil if the user can create a new user account.
func CanCreateUserAccount(creator, newAcct *AppAccount) error {
	if !IsAdminUser(creator) {
		return fmt.Errorf("only admins can create user accounts")
	}
	if !creator.IsActive() {
		return fmt.Errorf("the account is disabled")
	}
	isCustodianOrMemberAdmin := IsAdminUser(newAcct) && !IsClearingHouse(newAcct)
	if IsClearingHouse(creator) {
		if !isCustodianOrMemberAdmin && !BelongToSameEntity(creator, newAcct) {
			return fmt.Errorf(
				"can only create admin accounts for its own clearing house or admin accounts for other entities")
		}
		return nil
	}
	// members and custodian can create their own users only
	if !BelongToSameEntity(creator, newAcct) {
		return fmt.Errorf("members and custodian can create their own users only")
	}
	// Only Clearing House's admins can create other admin accounts
	if BelongToSameEntity(creator, newAcct) && IsAdminUser(newAcct) {
		return fmt.Errorf("only admins of the clearing house can create admin accounts")
	}
	return nil
}

// CreateAssetAccount is the function that, given an admin user and the
// new account's public key, instantiate a new asset account owned by
// the admin itsel.
func CreateAssetAccount(creator *AppAccount, pub crypto.PubKey, cash sdk.Coins) (*AppAccount, error) {
	if !IsAdminUser(creator) { // Only admins can create asset accounts.
		return nil, fmt.Errorf("only admins can create asset accounts")
	}
	if !creator.IsActive() {
		return nil, fmt.Errorf("the account is disabled")
	}
	return NewAssetAccount(pub, cash, creator.Address,
		creator.GetLegalEntityName(), creator.GetLegalEntityType()), nil
}

// AccountMapper creates an account mapper given a storekey
func AccountMapper(capKey sdk.StoreKey) sdk.AccountMapper {
	var accountMapper = auth.NewAccountMapper(
		capKey,        // target store
		&AppAccount{}, // prototype
	)

	// Register all interfaces and concrete types that
	// implement those interfaces, here.
	cdc := accountMapper.WireCodec()
	crypto.RegisterWire(cdc)
	// auth.RegisterWireBaseAccount(cdc)

	// Make WireCodec inaccessible before sealing
	res := accountMapper.Seal()
	return res
}
