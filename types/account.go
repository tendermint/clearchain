package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	crypto "github.com/tendermint/go-crypto"
)

// ensure AppAccount implements the sdk.Account interface
var _ sdk.Account = (*AppAccount)(nil)

// AppAccount defines the properties of an AppAccount
type AppAccount struct {
	auth.BaseAccount
	Type string

	// TODO: fields that may potentially be introduced in future
	// Name               string
	// LegalEntityAddress crypto.Address
	// Creator            crypto.Address
}

func IsCustodian(a *AppAccount) bool {
	return a.Type == EntityCustodian
}

func IsClearingHouse(a *AppAccount) bool {
	return a.Type == EntityClearingHouse
}

func IsGeneralClearingMember(a *AppAccount) bool {
	return a.Type == EntityGeneralClearingMember
}

func IsIndividualClearingMember(a *AppAccount) bool {
	return a.Type == EntityIndividualClearingMember
}

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

// // GetName returns account's name.
// func  GetName() string {
// 	return a.Name
// }

// // SetName modifies account's name
// func (a *AppAccount) SetName(name string) {
// 	a.Name = name
// }

// // GetCreator returns account's name.
// func  GetCreator() crypto.Address {
// 	return a.Creator
// }

// // SetCreator modifies account's name
// func (a *AppAccount) SetCreator(creator crypto.Address) {
// 	a.Creator = creator
// }
