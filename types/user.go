package types

import (
	"fmt"

	"github.com/tendermint/go-crypto"
)

// User defines the attribute of a ledger's user
type User struct {
	PubKey      crypto.PubKey `json:"pub_key"`     // May be nil, if not known.
	Name        string        `json:"name"`        // Human-readable identifier, mandatory
	EntityID    string        `json:"entity_id"`   // LegalEntity's ID
	Permissions Perm          `json:"permissions"` // User is disabled if empty
}

// NewUser initializes a new user
func NewUser(pubKey crypto.PubKey, name string, entityID string, permissions Perm) *User {
	if len(name) == 0 || pubKey == nil {
		return nil
	}
	return &User{
		PubKey:      pubKey,
		Name:        name,
		EntityID:    entityID,
		Permissions: permissions}
}

// Equal provides an equality operator
func (u *User) Equal(v *User) bool {
	if u != nil && v != nil {
		return u.PubKey.Equals(v.PubKey) && u.Name == v.Name && u.EntityID == v.EntityID && u.Permissions == v.Permissions
	}
	return u == v
}

// CanExecTx determines whether a LegalEntity can execute a Tx
func (u *User) CanExecTx(txType byte) bool {
	return u.Permissions.Has(permissionsMapByTxType[txType])
}

// VerifySignature verifies a signed message against the User's PubKey.
func (u *User) VerifySignature(signBytes []byte, signature crypto.Signature) bool {
	return u.PubKey.VerifyBytes(signBytes, signature)
}

func (u *User) String() string {
	if u == nil {
		return "nil-User"
	}
	return fmt.Sprintf("User{%s %q %v}", u.EntityID, u.Name, u.Permissions)
}

//--------------------------------------------

// UserGetter is implemented by any value that has a GetUser
type UserGetter interface {
	GetUser(addr []byte) *User
}

// UserSetter is implemented by any value that has a SetUser
type UserSetter interface {
	SetUser(addr []byte, acc *User)
}

// UserGetterSetter is implemented by any value that has both
// GetUser and SetUser
type UserGetterSetter interface {
	GetUser(addr []byte) *User
	SetUser(addr []byte, acc *User)
}

//-----------------------------------------

// PrivUser defines the attributes of a private user
type PrivUser struct {
	crypto.PrivKey
	User
}
