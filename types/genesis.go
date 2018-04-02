package types

import (
	"encoding/hex"

	crypto "github.com/tendermint/go-crypto"
)

// GenesisState defines the app's initial state to unmarshal.
type GenesisState struct {
	ClearingHouseAdmin GenesisAccount `json:"ch_admin"`
}

// GenesisAccount is an abstraction of the accounts specified in a genesis file
type GenesisAccount struct {
	PubKeyHexa string `json:"public_key"`
	EntityName string `json:"entity_name"`
}

// ToClearingHouseAdmin converts  a GenesisAccount into an AppAccount (a Clearing House admin user)
func (ga *GenesisAccount) ToClearingHouseAdmin() (acc *AppAccount, err error) {
	// Done manually since JSON Unmarshalling does not create a PubKey from a hexa value
	pubBytes, err := hex.DecodeString(ga.PubKeyHexa)
	if err != nil {
		return nil, err
	}
	publicKey, err := crypto.PubKeyFromBytes(pubBytes)
	if err != nil {
		return nil, err
	}

	adminUser := NewAdminUser(publicKey, nil, ga.EntityName, EntityClearingHouse)
	return adminUser, nil
}

// PubKeyFromHexString converts a hexadecimal string representation of
// a public key into a crypto.PubKey instance.
func PubKeyFromHexString(s string) (crypto.PubKey, error) {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return crypto.PubKey{}, err
	}
	return crypto.PubKeyFromBytes(bytes)
}
