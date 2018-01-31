package types

import (
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

var entityTypes []string = []string{EntityClearingHouse, EntityGeneralClearingMember, EntityIndividualClearingMember, EntityCustodian}

// LegalEntity defines the properties of a legal entity
type LegalEntity interface {
	GetAddress() crypto.Address
	GetType() string
	GetAdmins() []crypto.Address
	GetOps() []crypto.Address
	GetParent() crypto.Address
}

// BaseLegalEntity defines the properties of a concrete legal entity
type BaseLegalEntity struct {
	auth.BaseAccount
	Type   string
	Admins []crypto.Address
	Ops    []crypto.Address
	Parent crypto.Address
}

func (e BaseLegalEntity) GetAddress() crypto.Address {
	return e.Address
}

func (e BaseLegalEntity) GetType() string {
	return e.Type
}

func (e BaseLegalEntity) GetAdmins() []crypto.Address {
	return e.Admins
}

func (e BaseLegalEntity) GetOps() []crypto.Address {
	return e.Ops
}

func (e BaseLegalEntity) GetParent() crypto.Address {
	return e.Parent
}

func IsValidEntityType(entityType string) bool {
	return sliceContainsString(entityTypes, entityType)
}

func sliceContainsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}