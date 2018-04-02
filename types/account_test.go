package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	crypto "github.com/tendermint/go-crypto"
)

func TestBelongToSameEntity(t *testing.T) {
	acct1, _ := makeAdminUser("ent1", EntityClearingHouse)
	acct2, _ := makeUser("ent2", EntityClearingHouse)
	acct3, _ := makeAssetAccount(nil, "ent2", EntityClearingHouse)
	acct4, _ := makeAssetAccount(nil, "ent2", EntityCustodian)
	type args struct {
		acct1 *AppAccount
		acct2 *AppAccount
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"identity", args{acct1, acct1}, true},
		{"same type, different name", args{acct1, acct2}, false},
		{"different type, same name", args{acct3, acct4}, false},
		{"different account, same entity", args{acct2, acct3}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BelongToSameEntity(tt.args.acct1, tt.args.acct2); got != tt.want {
				t.Errorf("BelongToSameEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sliceContainsString(t *testing.T) {
	stringSlice := []string{"xxx", "yyy", "zzz"}
	type args struct {
		slice  []string
		target string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"nil string", args{stringSlice, ""}, false},
		{"spotted", args{stringSlice, "xxx"}, true},
		{"empty slice", args{[]string{}, "xxx"}, false},
		{"nil slice", args{[]string{}, "xxx"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sliceContainsString(tt.args.slice, tt.args.target)
			assert.Equal(t, tt.want, got)
		})
	}
}

/* Auxiliary functions.
 */

func makeUser(entname, typ string) (*AppAccount, crypto.PrivKey) {
	priv := crypto.GenPrivKeyEd25519()
	return NewOpUser(priv.PubKey(), nil, entname, typ), priv.Wrap()
}

func makeAdminUser(entname, typ string) (*AppAccount, crypto.PrivKey) {
	priv := crypto.GenPrivKeyEd25519()
	return NewAdminUser(priv.PubKey(), nil, entname, typ), priv.Wrap()
}

func makeAssetAccount(cash sdk.Coins, entname, typ string) (*AppAccount, crypto.Address) {
	pub := crypto.GenPrivKeyEd25519().PubKey()
	return NewAssetAccount(pub, cash, nil, entname, typ), pub.Address()
}

func TestGetAccountDecoder(t *testing.T) {
	cdc := MakeCodec()
	decoder := GetAccountDecoder(cdc)
	admin, _ := makeAdminUser("member", "gcm")
	adminBz, _ := cdc.MarshalBinary(admin)
	tests := []struct {
		name    string
		account *AppAccount
		bz      []byte
		wantErr bool
	}{
		{"admin", admin, adminBz, false},
		{"nil", &AppAccount{}, []byte{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acct, err := decoder(tt.bz)
			assert.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				concreteAcct := acct.(*AppAccount)
				assert.True(t, accountEqual(concreteAcct, tt.account))
			}
		})
	}
}
