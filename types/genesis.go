package types

import (
	"encoding/hex"
	crypto "github.com/tendermint/go-crypto"
)

// State to Unmarshal
type GenesisState struct {
	ClearingHouseAdmin GenesisAccount `json:"ch_admin"`
}

// GenesisAccount is an abstraction of the accounts specified in a genesis file
type GenesisAccount struct {
	PubKeyHexa string `json:"public_key"`
	EntityName string `json:"entity_name"`	

	/* future support
	EntityType string `json:"entity_type"`
	Creator 	crypto.Address 	`json:"creator_address"`
	Coins 		sdk.Coins    	`json:"coins"`
	IsActive    bool 			`json:"is_active"`
	IsAdmin     bool 			`json:"is_admin"`
	AccountType string 			`json:"account_type"`
	*/
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
