package types

import (		
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"	

	crypto "github.com/tendermint/go-crypto"	
)

// Test_ToAdminUser verifies that a GenesisAccount is converted correctly into AppAccount (of type Admin User)
func Test_ToAdminUser(t *testing.T) {	
	
	cdc := MakeTxCodec()

	pubBytes, _ := hex.DecodeString("328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f874")
	publicKey1, _ := crypto.PubKeyFromBytes(pubBytes)

	pubBytes, _ = hex.DecodeString("328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f875")
	publicKey2, _ := crypto.PubKeyFromBytes(pubBytes)

	pubBytes, _ = hex.DecodeString("328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f876")
	publicKey3, _ := crypto.PubKeyFromBytes(pubBytes)

	adminCreated1 := NewAdminUser(publicKey1, nil, "FXCH", "ch")
	adminCreated2 := NewAdminUser(publicKey2, nil, "ClearingHouse", "ch")
	adminCreated3 := NewAdminUser(publicKey3, nil, "Admin", "gcm")
	
	type args struct {	
		jsonValue string
	}
	tests := []struct {
		name       string
		args       args
		expectedAccount   *AppAccount	
	}{
		{
			"instantiate admin 1 ok", args{jsonValue: "{\"public_key\":\"328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f874\", \"entity_name\":\"FXCH\", \"entity_type\":\"ch\"}"}, adminCreated1,
		},	
		{
			"instantiate admin 2 ok", args{jsonValue: "{\"public_key\":\"328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f875\", \"entity_name\":\"ClearingHouse\", \"entity_type\":\"ch\"}"}, adminCreated2,
		},
		{
			"instantiate admin 3 ok", args{jsonValue: "{\"public_key\":\"328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f876\", \"entity_name\":\"Admin\", \"entity_type\":\"gcm\"}"}, adminCreated3,
		},	
	}

	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {	
			var genesisAcc GenesisAccount 		
			err := cdc.UnmarshalJSON([]byte(tt.args.jsonValue), &genesisAcc)		
			assert.Nil(t, err)			
			adminUser, _ := genesisAcc.ToAdminUser()
			assert.Equal(t, hex.EncodeToString(tt.expectedAccount.Address), hex.EncodeToString(adminUser.Address))
			assert.True(t, tt.expectedAccount.PubKey.Equals(adminUser.PubKey))	
			assert.True(t, adminUser.Coins.IsZero())		
			assert.True(t, adminUser.Creator == nil)	
			assert.Equal(t, tt.expectedAccount.EntityName, adminUser.EntityName)			
			assert.Equal(t, tt.expectedAccount.EntityType, adminUser.EntityType)	
			assert.Equal(t, tt.expectedAccount.AccountType, adminUser.AccountType)
			assert.Equal(t, AccountUser, adminUser.AccountType)
			assert.Equal(t, tt.expectedAccount.Active, adminUser.Active)
			assert.True(t, adminUser.Active)
			assert.Equal(t, tt.expectedAccount.Admin, adminUser.Admin)
			assert.True(t, adminUser.Admin)
		})
	}
}