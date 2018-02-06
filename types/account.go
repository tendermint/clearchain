package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	crypto "github.com/tendermint/go-crypto"
)

// EntityType string identifiers
const (
	EntityClearingHouse            = "ch"
	EntityGeneralClearingMember    = "gcm"
	EntityIndividualClearingMember = "icm"
	EntityCustodian                = "custodian"
)

var (
	// ensure AppAccount implements the sdk.Account interface
	_ sdk.Account = (*AppAccount)(nil)
)

// AppAccount defines the properties of an AppAccount.
type AppAccount struct {
	auth.BaseAccount
	Type            string
	Creator         crypto.Address
	EntityAdmin     bool
	LegalEntityName string
}

// NewAppAccount constructs a new account instance.
func NewAppAccount(pub crypto.PubKey, cash sdk.Coins, typ string, creator crypto.Address, isAdmin bool, entity string) *AppAccount {
	acct := new(AppAccount)
	acct.SetAddress(pub.Address())
	acct.SetPubKey(pub)
	acct.SetCoins(cash)
	acct.SetCreator(creator)
	acct.Type = typ
	acct.EntityAdmin = isAdmin
	acct.LegalEntityName = entity
	return acct
}

// GetCreator returns account's creator.
func GetCreator(a *AppAccount) crypto.Address {
	return a.Creator
}

// SetCreator modifies account's creator.
func (a *AppAccount) SetCreator(creator crypto.Address) {
	a.Creator = creator
}

// IsEntityAdmin returns true if the account is admin
// of its legal entity; false otherwise.
func IsEntityAdmin(a *AppAccount) bool {
	return a.EntityAdmin
}

// IsCustodian returns true if the account's owner entity
// is a custodian; false otherwise.
func IsCustodian(a *AppAccount) bool {
	return a.Type == EntityCustodian
}

// IsClearingHouse returns true if the account's owner entity
// is the clearing house; false otherwise.
func IsClearingHouse(a *AppAccount) bool {
	return a.Type == EntityClearingHouse
}

// IsGeneralClearingMember returns true if the account's owner entity
// is a general clearing member; false otherwise.
func IsGeneralClearingMember(a *AppAccount) bool {
	return a.Type == EntityGeneralClearingMember
}

// IsIndividualClearingMember returns true if the account's owner entity
// is an individual clearing member; false otherwise.
func IsIndividualClearingMember(a *AppAccount) bool {
	return a.Type == EntityIndividualClearingMember
}

// IsMember returns true if the account's owner entity is either
// a general or an individual clearing member; false otherwise.
func IsMember(a *AppAccount) bool {
	return IsIndividualClearingMember(a) ||
		IsGeneralClearingMember(a)
}

// BelongToSameEntity returns true if two accounts
// belong to the same legal entity.
func BelongToSameEntity(acct1, acct2 *AppAccount) bool {
	return (acct1.Type == acct2.Type) && (acct1.LegalEntityName == acct2.LegalEntityName)
}

// CanCreate returns nil if creator can create acct.
func CanCreate(creator, acct *AppAccount) error {
	if !IsEntityAdmin(creator) {
		// Only admins can create accounts.
		return fmt.Errorf("only admins can create accounts")
	}
	if !IsClearingHouse(creator) { // Members and Custodian can only create their own accounts.
		if !BelongToSameEntity(creator, acct) {
			return fmt.Errorf("members and custodian can create their own accounts only")
		}
		// Only Clearing House's admins can create other admin accounts
		if IsEntityAdmin(acct) {
			return fmt.Errorf("only admins of the clearing house can create admin accounts")
		}
	} else {
		// Clearing house's admins can create admin accounts for
		// other entities and any other clearing house accounts
		isCustodianOrMemberAdmin := func(a *AppAccount) bool { return IsEntityAdmin(a) && !IsClearingHouse(a) }
		if !(BelongToSameEntity(creator, acct) || isCustodianOrMemberAdmin(acct)) {
			return fmt.Errorf(
				"can only create other clearing house's accounts or admin accounts for other entities")
		}
	}
	return nil
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

func sliceContainsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
