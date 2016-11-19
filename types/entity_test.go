package types

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/tendermint/go-common"
)

func TestNewLegalEntity(t *testing.T) {
	uuid := uuid.NewV4().String()
	type args struct {
		id          string
		t           byte
		name        string
		permissions Perm
		creatorAddr []byte
	}
	tests := []struct {
		name string
		args args
		want *LegalEntity
	}{
		{"newEntity", args{uuid, EntityTypeCHByte, "", Perm(0), []byte{}}, &LegalEntity{ID: uuid, Type: EntityTypeCHByte, Name: "", Permissions: Perm(0)}},
	}
	for _, tt := range tests {
		if got := NewLegalEntity(tt.args.id, tt.args.t, tt.args.name, tt.args.permissions, tt.args.creatorAddr); !got.Equal(tt.want) {
			t.Errorf("%q. NewLegalEntity() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestLegalEntity_Equal(t *testing.T) {
	id := uuid.NewV4().String()
	type fields struct {
		ID          string
		Type        byte
		Name        string
		Permissions Perm
		CreatorAddr []byte
	}
	type args struct {
		e *LegalEntity
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"equal",
			fields{id, byte(0xFF), "test_name", PermTransferTx, []byte{}},
			args{&LegalEntity{id, byte(0xFF), "test_name", PermTransferTx, []byte{}}},
			true,
		},
		{"notEqual",
			fields{id, byte(0xFF), "test_name", PermTransferTx, []byte{}},
			args{&LegalEntity{uuid.NewV4().String(), byte(0xFF), "test_name", PermTransferTx, []byte{}}},
			false,
		},
	}

	for _, tt := range tests {
		l := &LegalEntity{
			ID:          tt.fields.ID,
			Type:        tt.fields.Type,
			Name:        tt.fields.Name,
			Permissions: tt.fields.Permissions,
			CreatorAddr: tt.fields.CreatorAddr,
		}
		if got := l.Equal(tt.args.e); got != tt.want {
			t.Errorf("%q. LegalEntity.Equal() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestLegalEntity_CanExecTx(t *testing.T) {
	allowedTxs := []byte{TxTypeTransfer, TxTypeQueryAccount}
	notAllowedTxs := []byte{TxTypeCreateUser, TxTypeCreateLegalEntity}
	type fields struct {
		Permissions Perm
	}
	type args struct {
		txs []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"canExec", fields{NewPermByTxType(allowedTxs...)}, args{allowedTxs}, true},
		{"canExec", fields{NewPermByTxType(allowedTxs...)}, args{notAllowedTxs}, false},
	}
	for _, tt := range tests {
		e := LegalEntity{Permissions: tt.fields.Permissions}
		for _, b := range tt.args.txs {
			got := e.CanExecTx(b)
			if got != tt.want {
				t.Errorf("%q. LegalEntity.CanExecTx() = %v, want %v", tt.name, got, tt.want)
			}
		}
	}
}

func TestLegalEntity_String(t *testing.T) {
	id := uuid.NewV4().String()
	type args struct {
		id string
		t  byte
		n  string
		p  Perm
		c  []byte
	}
	testcases := []struct {
		name string
		args *args
		want string
	}{
		{"nonEmpty", &args{id, 0x01, "CH1", PermTransferTx, []byte{}}, common.Fmt("LegalEntity{%x %s \"CH1\" %v %s}", EntityTypeCHByte, id, PermTransferTx, "")},
		{"nil", nil, "nil-LegalEntity"},
	}

	for _, tc := range testcases {
		var e *LegalEntity
		if tc.args != nil {
			e = NewLegalEntity(tc.args.id, tc.args.t, tc.args.n, tc.args.p, tc.args.c)
		}
		if ret := e.String(); ret != tc.want {
			t.Errorf("%q: String() return %q, expected: %q", tc.name, ret, tc.want)
		}
	}
}

func TestIsValidEntityType(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"valid", args{EntityTypeGCMByte}, true},
		{"invalid", args{byte(0xFF)}, false},
	}
	for _, tt := range tests {
		if got := IsValidEntityType(tt.args.b); got != tt.want {
			t.Errorf("%q. IsValidEntityType() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
