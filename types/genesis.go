package types

import (
	"encoding/hex"
	crypto "github.com/tendermint/go-crypto"
)

// State to Unmarshal
type GenesisState struct {
	AdminUsers []GenesisAccount `json:"admin_users"`
}

// GenesisAccount is an abstraction of the accounts specified in a genesis file
type GenesisAccount struct {
	PubKeyHexa string `json:"public_key"`
	EntityName string `json:"entity_name"`
	EntityType string `json:"entity_type"`

	/* future support
	Creator 	crypto.Address 	`json:"creator_address"`
	Coins 		sdk.Coins    	`json:"coins"`
	IsActive    bool 			`json:"is_active"`
	IsAdmin     bool 			`json:"is_admin"`
	AccountType string 			`json:"account_type"`
	*/
}

// ToAdminUser converts  GenesisAccount into an AppAccount (Admin User)
func (ga *GenesisAccount) ToAdminUser() (acc *AppAccount, err error) {

	// Done manually since JSON Unmarshalling does not create a PubKey from a hexa value
	pubBytes, _ := hex.DecodeString(ga.PubKeyHexa)
	publicKey, _ := crypto.PubKeyFromBytes(pubBytes)

	adminUser := NewAdminUser(publicKey, nil, ga.EntityName, ga.EntityType)
	return adminUser, nil
}
