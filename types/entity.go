package types

import (
	"bytes"
	"fmt"
	"github.com/tendermint/go-common"
)

// EntityType byte identifiers
const (
	EntityTypeCHByte        = byte(0x01)
	EntityTypeGCMByte       = byte(0x02)
	EntityTypeICMByte       = byte(0x03)
	EntityTypeCustodianByte = byte(0x04)
)

// IsValidEntityType checks whether a byte is a valid type for an entity.
func IsValidEntityType(b byte) bool {
	return bytes.Contains([]byte{EntityTypeCHByte, EntityTypeGCMByte, EntityTypeICMByte, EntityTypeCustodianByte}, []byte{b})
}

// LegalEntity defines the attributes of a legal entity
type LegalEntity struct {
	ID          string `json:"id"`           // LegalEntity's ID
	EntityID    string `json:"entity_id"`    // LegalEntity's owner
	Type        byte   `json:"type"`         // Mandatory
	Name        string `json:"name"`         // This could be empty
	Permissions Perm   `json:"permissions"`  // Set of allowed Txs
	CreatorAddr []byte `json:"creator_addr"` // ID of the creator of the Clearing House that created the legal entity
}

// NewLegalEntityByType is a convenience function to create a legal entity according to the type given.
func NewLegalEntityByType(t byte, id string, name string, creatorAddr []byte, EntityID string) *LegalEntity {
	switch t {
	case EntityTypeCHByte:
		return NewCH(id, name, creatorAddr, EntityID)
	case EntityTypeCustodianByte:
		return NewCustodian(id, name, creatorAddr, EntityID)
	case EntityTypeGCMByte:
		return NewGCM(id, name, creatorAddr, EntityID)
	case EntityTypeICMByte:
		return NewICM(id, name, creatorAddr, EntityID)
	default:
		common.PanicSanity(common.Fmt("Unexpected TxType: %x", t))
	}
	return nil
}

// NewCH is a convenience function to create a new CH
func NewCH(id string, name string, creatorAddr []byte, EntityID string) *LegalEntity {
	return NewLegalEntity(id, EntityTypeCHByte, name, NewPermByTxType(
		TxTypeTransfer, TxTypeQueryAccount, TxTypeCreateAccount, TxTypeCreateLegalEntity, TxTypeCreateUser,
	), creatorAddr, EntityID)
}

// NewGCM is a convenience function to create a new GCM
func NewGCM(id string, name string, creatorAddr []byte, EntityID string) *LegalEntity {
	return NewLegalEntity(id, EntityTypeGCMByte, name, NewPermByTxType(
		TxTypeTransfer, TxTypeQueryAccount, TxTypeCreateAccount, TxTypeCreateUser), creatorAddr, EntityID)
}

// NewICM is a convenience function to create a new ICM
func NewICM(id string, name string, creatorAddr []byte, EntityID string) *LegalEntity {
	return NewLegalEntity(id, EntityTypeICMByte, name, NewPermByTxType(
		TxTypeTransfer, TxTypeQueryAccount, TxTypeCreateAccount, TxTypeCreateUser), creatorAddr, EntityID)
}

// NewCustodian is a convenience function to create a new Custodian
func NewCustodian(id string, name string, creatorAddr []byte, EntityID string) *LegalEntity {
	return NewLegalEntity(id, EntityTypeGCMByte, name, NewPermByTxType(
		TxTypeTransfer, TxTypeQueryAccount, TxTypeCreateAccount, TxTypeCreateUser), creatorAddr, EntityID)
}

// NewLegalEntity initializes a new LegalEntity
func NewLegalEntity(id string, t byte, name string, permissions Perm, creatorAddr []byte, EntityID string) *LegalEntity {
	return &LegalEntity{ID: id, Type: t, Name: name, Permissions: permissions, CreatorAddr: creatorAddr, EntityID: EntityID}
}

// Equal provides an equality operator
func (l *LegalEntity) Equal(e *LegalEntity) bool {
	if l != nil && e != nil {
		return l.ID == e.ID && l.Type == e.Type && bytes.Equal(l.CreatorAddr, e.CreatorAddr) &&
			l.Name == e.Name && l.Permissions == e.Permissions && l.EntityID == e.EntityID
	}
	return l == e
}

// CanExecTx determines whether a LegalEntity can execute a Tx
func (l *LegalEntity) CanExecTx(txType byte) bool {
	return l.Permissions.Has(permissionsMapByTxType[txType])
}

func (l *LegalEntity) String() string {
	if l == nil {
		return "nil-LegalEntity"
	}
	return fmt.Sprintf("LegalEntity{%x %s %q %v %x %v}", l.Type, l.ID, l.Name, l.Permissions, l.CreatorAddr, l.EntityID)
}

//--------------------------------------------

// LegalEntityGetter is implemented by any value that has a GetLegalEntity
type LegalEntityGetter interface {
	GetLegalEntity(id string) *LegalEntity
}

// LegalEntitySetter is implemented by any value that has a SetLegalEntity
type LegalEntitySetter interface {
	SetLegalEntity(id string, acc *LegalEntity)
}

// LegalEntityGetterSetter is implemented by any value that has both
// GetLegalEntity and SetLegalEntity
type LegalEntityGetterSetter interface {
	GetLegalEntity(id string) *LegalEntity
	SetLegalEntity(id string, acc *LegalEntity)
}
