package types

import (
	"fmt"
	"strings"
)

// Identifiers for legal entities types.
const (
	EntityClearingHouse            = "ch"
	EntityGeneralClearingMember    = "gcm"
	EntityIndividualClearingMember = "icm"
	EntityCustodian                = "custodian"
)

// LegalEntity is the interface that wraps the basic accessor methods
// to set and get entities attributes.
type LegalEntity interface {
	SetLegalEntityName(string)
	GetLegalEntityName() string
	SetLegalEntityType(string)
	GetLegalEntityType() string
}

// BaseLegalEntity defines the properties of a legal entity.
type BaseLegalEntity struct {
	EntityName string
	EntityType string
}

// SetLegalEntityName modifies an entity's name.
func (e BaseLegalEntity) SetLegalEntityName(name string) {
	e.EntityName = name
}

// GetLegalEntityName returns the entity's name.
func (e BaseLegalEntity) GetLegalEntityName() string {
	return e.EntityName
}

// SetLegalEntityType modifies an entity's type.
func (e BaseLegalEntity) SetLegalEntityType(typ string) {
	e.EntityType = typ
}

// GetLegalEntityType returns the entity's type.
func (e BaseLegalEntity) GetLegalEntityType() string {
	return e.EntityType
}

// IsCustodian returns true if the account's owner entity
// is a custodian; false otherwise.
func IsCustodian(e LegalEntity) bool {
	return e.GetLegalEntityType() == EntityCustodian
}

// IsClearingHouse returns true if the account's owner entity
// is the clearing house; false otherwise.
func IsClearingHouse(e LegalEntity) bool {
	return e.GetLegalEntityType() == EntityClearingHouse
}

// IsGeneralClearingMember returns true if the account's owner entity
// is a general clearing member; false otherwise.
func IsGeneralClearingMember(e LegalEntity) bool {
	return e.GetLegalEntityType() == EntityGeneralClearingMember
}

// IsIndividualClearingMember returns true if the account's owner entity
// is an individual clearing member; false otherwise.
func IsIndividualClearingMember(e LegalEntity) bool {
	return e.GetLegalEntityType() == EntityIndividualClearingMember
}

// IsMember returns true if the account's owner entity is either
// a general or an individual clearing member; false otherwise.
func IsMember(e LegalEntity) bool {
	return IsIndividualClearingMember(e) || IsGeneralClearingMember(e)
}

// BelongToSameEntity returns true if two accounts
// belong to the same legal entity.
func BelongToSameEntity(e1, e2 LegalEntity) bool {
	return (e1.GetLegalEntityName() == e2.GetLegalEntityName()) &&
		(e1.GetLegalEntityType() == e2.GetLegalEntityType())
}

// ValidateLegalEntity performs basic validation
// on types that implement LegalEntity.
func ValidateLegalEntity(e LegalEntity) error {
	if len(strings.TrimSpace(e.GetLegalEntityName())) == 0 {
		return fmt.Errorf("legal entity name must be non-nil")
	}
	if !sliceContainsString([]string{
		EntityClearingHouse,
		EntityGeneralClearingMember,
		EntityIndividualClearingMember,
		EntityCustodian}, e.GetLegalEntityType()) {
		return fmt.Errorf("legal entity type %q is invalid", e.GetLegalEntityType())
	}
	return nil
}

func sliceContainsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
