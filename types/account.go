package types

import (
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
	creatableEntities = []string{EntityGeneralClearingMember, EntityIndividualClearingMember, EntityCustodian}

	// ensure AppAccount implements the sdk.Account interface
	_ sdk.Account = (*AppAccount)(nil)
)

// AppAccount defines the properties of an AppAccount
type AppAccount struct {
	auth.BaseAccount
	Type            string
	Creator         crypto.Address
	EntityAdmin     bool
	LegalEntityName string
}

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

// GetCreator returns account's creator.
func GetCreator(a *AppAccount) crypto.Address {
	return a.Creator
}

// SetCreator modifies account's creator.
func (a *AppAccount) SetCreator(creator crypto.Address) {
	a.Creator = creator
}

// IsCreatableEntity checks whether it is allowed to create an entity of the given type..
func IsCreatableEntity(entityType string) bool {
	return sliceContainsString(creatableEntities, entityType)
}

func sliceContainsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

// // GetName returns account's name.
// func GetName() string {
// 	return a.Name
// }

// // SetName modifies account's name
// func (a *AppAccount) SetName(name string) {
// 	a.Name = name
// }
