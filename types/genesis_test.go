package types

import (
	"encoding/hex"
	"testing"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crypto "github.com/tendermint/go-crypto"
)

// ToClearingHouseAdmin verifies that a GenesisAccount is converted correctly into an AppAccount (a Clearing House admin user)
func Test_ToClearingHouseAdmin(t *testing.T) {
	var publicKey1, publicKey2, publicKey3 crypto.PubKey
	cdc := MakeCodec()
	wire.RegisterCrypto(cdc)

	require.Nil(t, publicKey1.UnmarshalJSON([]byte(`{"type":"ed25519","data":"328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f874"}`)))
	require.Nil(t, publicKey2.UnmarshalJSON([]byte(`{"type":"ed25519","data":"328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f875"}`)))
	require.Nil(t, publicKey3.UnmarshalJSON([]byte(`{"type":"ed25519","data":"328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f876"}`)))

	adminCreated1 := NewAdminUser(publicKey1, nil, "ClearChain", "ch")
	adminCreated2 := NewAdminUser(publicKey2, nil, "ClearingHouse", "ch")
	adminCreated3 := NewAdminUser(publicKey3, nil, "Admin", "ch")

	type args struct {
		jsonValue string
	}
	tests := []struct {
		name            string
		args            args
		expectedAccount *AppAccount
	}{
		{
			"instantiate admin 1 ok", args{jsonValue: "{\"public_key\":\"01328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f874\", \"entity_name\":\"ClearChain\"}"}, adminCreated1,
		},
		{
			"instantiate admin 2 ok", args{jsonValue: "{\"public_key\":\"01328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f875\", \"entity_name\":\"ClearingHouse\"}"}, adminCreated2,
		},
		{
			"instantiate admin 3 ok", args{jsonValue: "{\"public_key\":\"01328eaf59335aa6724f253ca8f1620b249bb83e665d7e5134e9bf92079b2549df3572f876\", \"entity_name\":\"Admin\"}"}, adminCreated3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var genesisAcc GenesisAccount
			err := cdc.UnmarshalJSON([]byte(tt.args.jsonValue), &genesisAcc)
			assert.Nil(t, err)
			adminUser, err := genesisAcc.ToClearingHouseAdmin()
			assert.Nil(t, err)
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
