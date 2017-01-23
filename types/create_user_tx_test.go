package types

import (
	"bytes"
	"testing"

	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
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
		want   abci.Result
	}{
		{"emptyTx", fields{}, abci.ErrBaseInvalidInput},
		{"invalidAddress", fields{Address: []byte("")}, abci.ErrBaseInvalidInput},
		{"invalidUserAddr", fields{Address: randBytes()}, abci.ErrBaseInvalidPubKey},
		{"invalidName", fields{Address: randBytes(), PubKey: randPubKey()}, abci.ErrBaseInvalidInput},
		{"invalidSignature", fields{Address: randBytes(), PubKey: randPubKey(), Name: "name"}, abci.ErrBaseInvalidSignature},
		{"valid", fields{Address: randBytes(), PubKey: randPubKey(), Name: "name", Signature: genSig()}, abci.OK},
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

func TestCreateUserTx_SignTx(t *testing.T) {
	privKey := crypto.GenPrivKeyEd25519()
	addr := privKey.PubKey().Address()
	type fields struct {
		Address   []byte
		Name      string
		PubKey    crypto.PubKey
		CanCreate bool
		Signature crypto.Signature
	}
	type args struct {
		privateKey crypto.PrivKey
		chainID    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"validSignature", fields{addr, "user", nil, false, nil}, args{privKey, "chainID"}, false},
		{"invalidSignature", fields{addr, "user", nil, false, nil}, args{crypto.GenPrivKeyEd25519(), "chainID"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &CreateUserTx{
				Address:   tt.fields.Address,
				Name:      tt.fields.Name,
				PubKey:    tt.fields.PubKey,
				CanCreate: tt.fields.CanCreate,
				Signature: tt.fields.Signature,
			}
			if err := tx.SignTx(tt.args.privateKey, tt.args.chainID); (err != nil) != tt.wantErr {
				t.Errorf("CreateUserTx.SignTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
