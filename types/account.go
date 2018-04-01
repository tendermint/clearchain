package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/go-crypto"
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
	Creator     sdk.Address
	AccountType string
	Active      bool
	Admin       bool
}

// NewAppAccount constructs a new account instance.
func newAppAccount(pub crypto.PubKey, cash sdk.Coins, creator sdk.Address, typ string,
	isActive bool, isAdmin bool, entityName, entityType string) *AppAccount {
	baseaccount := auth.BaseAccount{
		Address: pub.Address(),
		PubKey:  pub,
		Coins:   cash,
	}
	entity := BaseLegalEntity{
		EntityName: entityName,
		EntityType: entityType,
	}
	return &AppAccount{
		BaseAccount:     baseaccount,
		BaseLegalEntity: entity,
		Creator:         creator,
		AccountType:     typ,
		Active:          isActive,
		Admin:           isAdmin,
	}
}

// NewOpUser constructs a new account instance, setting cash to nil.
func NewOpUser(pub crypto.PubKey, creator sdk.Address, entityName, entityType string) *AppAccount {
	return newAppAccount(pub, nil, creator, AccountUser, true, false, entityName, entityType)
}

// NewAdminUser constructs a new account instance, setting cash to nil.
func NewAdminUser(pub crypto.PubKey, creator sdk.Address, entityName, entityType string) *AppAccount {
	return newAppAccount(pub, nil, creator, AccountUser, true, true, entityName, entityType)
}

// NewAssetAccount constructs a new account instance.
func NewAssetAccount(pub crypto.PubKey, cash sdk.Coins, creator sdk.Address, entityName, entityType string) *AppAccount {
	return newAppAccount(pub, cash, creator, AccountAsset, true, false, entityName, entityType)
}

// GetAccountType returns the account type.
func (a AppAccount) GetAccountType() string {
	return a.AccountType
}

// IsActive returns true if the account is active; false otherwise.
func (a AppAccount) IsActive() bool {
	return a.Active
}

// IsAdmin returns true if the account is admin; false otherwise.
func (a AppAccount) IsAdmin() bool {
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

// NewAccountMapper creates an account mapper given a storekey
func NewAccountMapper(capKey sdk.StoreKey) sdk.AccountMapper {
	return auth.NewAccountMapperSealed(capKey, &AppAccount{})
}

// Get the AccountDecoder function for the custom AppAccount
func GetAccountDecoder(cdc *wire.Codec) sdk.AccountDecoder {
	return func(accBytes []byte) (res sdk.Account, err error) {
		if len(accBytes) == 0 {
			return nil, sdk.ErrTxDecode("accBytes are empty")
		}
		acct := new(AppAccount)
		err = cdc.UnmarshalBinary(accBytes, &acct)
		if err != nil {
			panic(err)
		}
		return acct, err
	}
}

/* auxiliary functions */

func accountEqual(a1, a2 *AppAccount) bool {
	return ((a1.AccountType == a2.AccountType) &&
		(a1.Admin == a2.Admin) &&
		(a1.Active == a2.Active) &&
		bytes.Equal(a1.Address, a2.Address) &&
		(bytes.Equal(a1.GetPubKey().Bytes(), a2.GetPubKey().Bytes())) &&
		BelongToSameEntity(a1, a2) &&
		bytes.Equal(a1.Creator, a2.Creator) &&
		a1.GetCoins().IsEqual(a2.GetCoins()))
}
