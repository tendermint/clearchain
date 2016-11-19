package types

import (
	"bytes"
	"testing"

	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
	tmsp "github.com/tendermint/tmsp/types"
)

func TestCreateUserTx_TxType(t *testing.T) {
	type fields struct {
		Address   []byte
		Name      string
		PubKey    crypto.PubKey
		CanCreate bool
		Signature crypto.Signature
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		{"default", fields{}, TxTypeCreateUser},
	}
	for _, tt := range tests {
		tx := &CreateUserTx{
			Address:   tt.fields.Address,
			Name:      tt.fields.Name,
			PubKey:    tt.fields.PubKey,
			CanCreate: tt.fields.CanCreate,
			Signature: tt.fields.Signature,
		}
		if got := tx.TxType(); got != tt.want {
			t.Errorf("%q. CreateUserTx.TxType() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCreateUserTx_SignBytes(t *testing.T) {
	chainID := "chainID"
	privKey := crypto.GenPrivKeyEd25519()
	tx := &CreateUserTx{
		Address:   privKey.PubKey().Address(),
		Name:      "new_user_name",
		Signature: nil,
	}
	signedBytes := tx.SignBytes(chainID)
	expected := append(wire.BinaryBytes(chainID), wire.BinaryBytes(tx)...)
	if !bytes.Equal(signedBytes, expected) {
		t.Errorf("CreateLegalEntityTx.SignBytes() = %v, want: %v", signedBytes, expected)
	}
}

func TestCreateUserTx_ValidateBasic(t *testing.T) {
	randPubKey := func() crypto.PubKey { return crypto.GenPrivKeyEd25519().PubKey() }
	randBytes := func() []byte { return crypto.CRandBytes(20) }
	genSig := func() crypto.Signature { return crypto.GenPrivKeyEd25519().Sign(randBytes()) }
	type fields struct {
		Address   []byte
		Name      string
		PubKey    crypto.PubKey
		CanCreate bool
		Signature crypto.Signature
	}
	tests := []struct {
		name   string
		fields fields
		want   tmsp.Result
	}{
		{"emptyTx", fields{}, tmsp.ErrBaseInvalidInput},
		{"invalidAddress", fields{Address: []byte("")}, tmsp.ErrBaseInvalidInput},
		{"invalidUserAddr", fields{Address: randBytes()}, tmsp.ErrBaseInvalidPubKey},
		{"invalidName", fields{Address: randBytes(), PubKey: randPubKey()}, tmsp.ErrBaseInvalidInput},
		{"invalidSignature", fields{Address: randBytes(), PubKey: randPubKey(), Name: "name"}, tmsp.ErrBaseInvalidSignature},
		{"valid", fields{Address: randBytes(), PubKey: randPubKey(), Name: "name", Signature: genSig()}, tmsp.OK},
	}
	for _, tt := range tests {
		tx := &CreateUserTx{
			Address:   tt.fields.Address,
			Name:      tt.fields.Name,
			PubKey:    tt.fields.PubKey,
			CanCreate: tt.fields.CanCreate,
			Signature: tt.fields.Signature,
		}
		if got := tx.ValidateBasic(); got.Code != tt.want.Code {
			t.Errorf("%q. CreateUserTx.ValidateBasic() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCreateUserTx_String(t *testing.T) {
	type fields struct {
		Address   []byte
		Name      string
		PubKey    crypto.PubKey
		CanCreate bool
		Signature crypto.Signature
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty", fields{}, "CreateUserTx{,,false,<nil>}"},
		{"stringRepr", fields{[]byte{0}, "", nil, true, nil}, "CreateUserTx{00,,true,<nil>}"},
	}
	for _, tt := range tests {
		tx := &CreateUserTx{
			Address:   tt.fields.Address,
			Name:      tt.fields.Name,
			PubKey:    tt.fields.PubKey,
			CanCreate: tt.fields.CanCreate,
			Signature: tt.fields.Signature,
		}
		if got := tx.String(); got != tt.want {
			t.Errorf("%q. CreateUserTx.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
