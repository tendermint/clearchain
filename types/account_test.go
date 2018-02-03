package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanCreate(t *testing.T) {
	chAdmin := makeAccount(EntityClearingHouse, nil, true)
	gcmAdmin := makeAccount(EntityGeneralClearingMember, nil, true)
	icmAdmin := makeAccount(EntityIndividualClearingMember, nil, true)
	custAdmin := makeAccount(EntityCustodian, nil, true)
	chOp := makeAccount(EntityClearingHouse, nil, false)
	gcmOp := makeAccount(EntityGeneralClearingMember, nil, false)
	icmOp := makeAccount(EntityIndividualClearingMember, nil, false)
	custOp := makeAccount(EntityCustodian, nil, false)

	chOp.LegalEntityName = chAdmin.LegalEntityName
	gcmOp.LegalEntityName = gcmAdmin.LegalEntityName
	icmOp.LegalEntityName = icmAdmin.LegalEntityName
	custOp.LegalEntityName = custAdmin.LegalEntityName

	type args struct {
		creator *AppAccount
		acct    *AppAccount
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"CH admin can create CH admin", args{chAdmin, chAdmin}, false},
		{"CH admin can create CUS admin", args{chAdmin, custAdmin}, false},
		{"CH admin can create GCM admin", args{chAdmin, gcmAdmin}, false},
		{"CH admin can create ICM admin", args{chAdmin, icmAdmin}, false},
		{"CH admin can create CH operator", args{chAdmin, chOp}, false},
		{"CH admin cannot create CUS operator", args{chAdmin, custOp}, true},
		{"CH admin cannot create GCM operator", args{chAdmin, gcmOp}, true},
		{"CH admin cannot create ICM operator", args{chAdmin, icmOp}, true},
		{"CH operator cannot create accounts", args{chOp, gcmOp}, true},
		{"GCM admin cannot create ICM accounts", args{gcmAdmin, icmOp}, true},
		{"GCM admin can create GCM accounts", args{gcmAdmin, gcmOp}, false},
		{"ICM admin can create ICM accounts", args{icmAdmin, icmOp}, false},
		{"CUST admin can create CUST accounts", args{custAdmin, custOp}, false},
		{"same entity types, different entity", args{gcmAdmin, makeAccount(EntityGeneralClearingMember, nil, false)}, true},
		{"both CH, different entity", args{chAdmin, makeAccount(EntityClearingHouse, nil, false)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CanCreate(tt.args.creator, tt.args.acct)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}
