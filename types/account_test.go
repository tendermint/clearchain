package types

import (
	"bytes"
	"fmt"
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

func TestAppAccount_CanCreateUserAccount(t *testing.T) {
	chAdmin, _ := makeAdminUser("CH", EntityClearingHouse)
	chOp, _ := makeUser("CH", EntityClearingHouse)
	otherChAdmin, _ := makeAdminUser("CH2", EntityClearingHouse)
	otherChUser, _ := makeUser("CH2", EntityClearingHouse)
	custAdmin, _ := makeAdminUser("CUST", EntityCustodian)
	custOp, _ := makeUser("CUST", EntityCustodian)
	otherCustOp, _ := makeUser("CUST2", EntityCustodian)
	icmAdmin, _ := makeAdminUser("ICM", EntityIndividualClearingMember)
	icmOp, _ := makeUser("ICM", EntityIndividualClearingMember)
	gcmAdmin, _ := makeAdminUser("GCM", EntityGeneralClearingMember)
	gcmOp, _ := makeUser("GCM", EntityGeneralClearingMember)
	chAsset, _ := makeAssetAccount(nil, chAdmin.EntityName, EntityClearingHouse)
	disabledAdmin, _ := makeAdminUser("ICM", EntityIndividualClearingMember)
	disabledAdmin.Active = false

	tests := []struct {
		name    string
		creator *AppAccount
		newAcct *AppAccount
		wantErr bool
	}{
		{"ch admin can create ch admin", chAdmin, chAdmin, false},
		{"ch admin can create ch ops", chAdmin, chOp, false},
		{"ch admin cannot create foreign ch admin", chAdmin, otherChAdmin, true},
		{"ch admin cannot create foreign ch ops", chAdmin, otherChUser, true},
		{"ch admin can create custodian admin", chAdmin, custAdmin, false},
		{"ch admin cannot create custodian ops", chAdmin, custOp, true},
		{"ch admin can create icm admin", chAdmin, icmAdmin, false},
		{"ch admin cannot create icm ops", chAdmin, icmOp, true},
		{"ch admin can create gcm admin", chAdmin, gcmAdmin, false},
		{"ch admin cannot create gcm ops", chAdmin, gcmOp, true},
		{"non-ch op cannot create admin", custOp, custAdmin, true},
		{"non-ch admin can create ops", custAdmin, custOp, false},
		{"non-ch admin and op must be same entity", custAdmin, otherCustOp, true},
		{"can't create asset accounts", custAdmin, chAsset, true},
		{"icm admin can create op", icmAdmin, icmOp, false},
		{"disabled icm admin cannot create op", disabledAdmin, icmOp, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CanCreateUserAccount(tt.creator, tt.newAcct)
			assert.Equal(t, tt.wantErr, (err != nil), fmt.Sprintf("%v", err))
		})
	}
}

func TestCreateAssetAccount(t *testing.T) {
	genKey := func() crypto.PubKey {
		return crypto.GenPrivKeyEd25519().PubKey()
	}
	chAdmin, _ := makeAdminUser("CH", EntityClearingHouse)
	custAdmin, _ := makeAdminUser("CUST", EntityCustodian)
	icmAdmin, _ := makeAdminUser("ICM", EntityIndividualClearingMember)
	gcmAdmin, _ := makeAdminUser("GCM", EntityGeneralClearingMember)
	gcmOp, _ := makeUser("GCM", EntityGeneralClearingMember)
	chInactiveAdm, _ := makeAdminUser("CH", EntityClearingHouse)
	chInactiveAdm.Active = false
	chAsset, _ := makeAssetAccount(nil, "CH", EntityClearingHouse)
	type args struct {
		creator *AppAccount
		pub     crypto.PubKey
		cash    sdk.Coins
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ch admin create asset", args{chAdmin, genKey(), sdk.Coins{{"USD", 1000}}}, false},
		{"cust admin create asset", args{custAdmin, genKey(), nil}, false},
		{"icm admin create asset", args{icmAdmin, genKey(), nil}, false},
		{"gcm admin create asset", args{gcmAdmin, genKey(), nil}, false},
		{"gcm op cannot create asset", args{gcmOp, genKey(), nil}, true},
		{"ch inactive admin cannot create asset", args{chInactiveAdm, genKey(), nil}, true},
		{"asset cannot create asset", args{chAsset, genKey(), nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateAssetAccount(tt.args.creator, tt.args.pub, tt.args.cash)
			assert.Equal(t, tt.wantErr, (err != nil))
			if err == nil {
				assert.EqualValues(t, tt.args.creator.BaseLegalEntity, got.BaseLegalEntity)
				assert.True(t, bytes.Equal(tt.args.creator.Address, got.Creator))
				assert.True(t, sdk.Coins.IsEqual(got.GetCoins(), tt.args.cash))
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
	return NewOpUser(priv.PubKey(), nil, entname, typ), priv
}

func makeAdminUser(entname, typ string) (*AppAccount, crypto.PrivKey) {
	priv := crypto.GenPrivKeyEd25519()
	return NewAdminUser(priv.PubKey(), nil, entname, typ), priv
}

func makeAssetAccount(cash sdk.Coins, entname, typ string) (*AppAccount, crypto.Address) {
	pub := crypto.GenPrivKeyEd25519().PubKey()
	return NewAssetAccount(pub, cash, nil, entname, typ), pub.Address()
}
