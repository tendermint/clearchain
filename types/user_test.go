package types

import (
	"testing"

	"github.com/tendermint/go-crypto"
)

func TestNewUser(t *testing.T) {
	privKey := crypto.GenPrivKeyEd25519()
	type args struct {
		pubKey      crypto.PubKey
		name        string
		entityID    string
		permissions Perm
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		{"nilPubKey", args{nil, "test", "entity", 0}, nil},
		{"emptyName", args{crypto.GenPrivKeyEd25519().PubKey(), "", "entity", 0}, nil},
		{"nonNil", args{privKey.PubKey(), "test", "entity", 0}, &User{privKey.PubKey(), "test", "entity", 0}},
	}
	for _, tt := range tests {
		if got := NewUser(tt.args.pubKey, tt.args.name, tt.args.entityID, tt.args.permissions); !got.Equal(tt.want) {
			t.Errorf("%q. NewUser() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestUser_Equal(t *testing.T) {
	privKey := crypto.GenPrivKeyEd25519()
	type fields struct {
		PubKey      crypto.PubKey
		Name        string
		EntityID    string
		Permissions Perm
	}
	type args struct {
		v *User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"equal", fields{privKey.PubKey(), "test", "entity", 0}, args{&User{privKey.PubKey(), "test", "entity", 0}}, true},
		{"notEqual", fields{privKey.PubKey(), "test", "", 0}, args{&User{privKey.PubKey(), "test", "entity", 0}}, false},
	}
	for _, tt := range tests {
		u := &User{
			PubKey:      tt.fields.PubKey,
			Name:        tt.fields.Name,
			EntityID:    tt.fields.EntityID,
			Permissions: tt.fields.Permissions,
		}
		if got := u.Equal(tt.args.v); got != tt.want {
			t.Errorf("%q. User.Equal() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestUser_CanExecTx(t *testing.T) {
	type fields struct {
		PubKey      crypto.PubKey
		Name        string
		EntityID    string
		Permissions Perm
	}
	type args struct {
		txType byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"canExec", fields{nil, "", "", NewPermByTxType(TxTypeTransfer, TxTypeCreateUser)}, args{TxTypeTransfer}, true},
		{"cantExec", fields{nil, "", "", NewPermByTxType(TxTypeTransfer, TxTypeCreateUser)}, args{}, false},
		{"noPermisssions", fields{nil, "", "", PermNone}, args{TxTypeTransfer}, false},
	}
	for _, tt := range tests {
		u := &User{
			PubKey:      tt.fields.PubKey,
			Name:        tt.fields.Name,
			EntityID:    tt.fields.EntityID,
			Permissions: tt.fields.Permissions,
		}
		if got := u.CanExecTx(tt.args.txType); got != tt.want {
			t.Errorf("%q. User.CanExecTx() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestUser_String(t *testing.T) {
	type fields struct {
		PubKey      crypto.PubKey
		Name        string
		EntityID    string
		Permissions Perm
	}
	tests := []struct {
		name   string
		fields *fields
		want   string
	}{
		{"nil", nil, "nil-User"},
		{"empty", &fields{}, "User{ \"\" 0}"},
		{"nonEmpty", &fields{nil, "test", "entity", PermTransferTx}, "User{entity \"test\" 1}"},
	}
	for _, tt := range tests {
		var u *User
		if tt.fields != nil {
			u = &User{
				PubKey:      tt.fields.PubKey,
				Name:        tt.fields.Name,
				EntityID:    tt.fields.EntityID,
				Permissions: tt.fields.Permissions,
			}
		}
		if got := u.String(); got != tt.want {
			t.Errorf("%q. User.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestUser_VerifySignature(t *testing.T) {
	privKey := crypto.GenPrivKeyEd25519()
	bs := crypto.CRandBytes(100)
	genKey := func() crypto.PrivKey {
		return crypto.GenPrivKeyEd25519()
	}
	type fields struct {
		PubKey      crypto.PubKey
		Name        string
		EntityID    string
		Permissions Perm
	}
	type args struct {
		signBytes []byte
		signature crypto.Signature
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"invalidSignature", fields{privKey.PubKey(), "", "", PermNone}, args{bs, genKey().Sign(bs)}, false},
		{"validSignature", fields{privKey.PubKey(), "", "", PermNone}, args{bs, privKey.Sign(bs)}, true},
	}
	for _, tt := range tests {
		u := &User{
			PubKey:      tt.fields.PubKey,
			Name:        tt.fields.Name,
			EntityID:    tt.fields.EntityID,
			Permissions: tt.fields.Permissions,
		}
		if got := u.VerifySignature(tt.args.signBytes, tt.args.signature); got != tt.want {
			t.Errorf("%q. User.VerifySignature() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
