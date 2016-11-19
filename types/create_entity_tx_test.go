package types

import (
	"bytes"
	"testing"

	uuid "github.com/satori/go.uuid"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
	tmsp "github.com/tendermint/tmsp/types"
)

func TestCreateLegalEntityTx_TxType(t *testing.T) {
	type fields struct {
		Address   []byte
		EntityID  string
		Type      byte
		Name      string
		Signature crypto.Signature
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		{"default", fields{}, TxTypeCreateLegalEntity},
	}
	for _, tt := range tests {
		tx := &CreateLegalEntityTx{
			Address:   tt.fields.Address,
			EntityID:  tt.fields.EntityID,
			Type:      tt.fields.Type,
			Name:      tt.fields.Name,
			Signature: tt.fields.Signature,
		}
		if got := tx.TxType(); got != tt.want {
			t.Errorf("%q. CreateLegalEntityTx.TxType() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCreateLegalEntityTx_SignBytes(t *testing.T) {
	chainID := "chainID"
	privKey := crypto.GenPrivKeyEd25519()
	tx := &CreateLegalEntityTx{
		Address:   privKey.PubKey().Address(),
		EntityID:  "entity_id",
		Signature: nil,
	}
	signedBytes := tx.SignBytes(chainID)
	expected := append(wire.BinaryBytes(chainID), wire.BinaryBytes(tx)...)
	if !bytes.Equal(signedBytes, expected) {
		t.Errorf("CreateLegalEntityTx.SignBytes() = %v, want: %v", signedBytes, expected)
	}
}

func TestCreateLegalEntityTx_ValidateBasic(t *testing.T) {
	randBytes := func() []byte { return crypto.CRandBytes(20) }
	genID := func() string { return uuid.NewV4().String() }
	genSig := func() crypto.Signature { return crypto.GenPrivKeyEd25519().Sign(randBytes()) }
	type fields struct {
		Address   []byte
		EntityID  string
		Type      byte
		Name      string
		Signature crypto.Signature
	}
	tests := []struct {
		name   string
		fields fields
		want   tmsp.Result
	}{
		{"emptyTx", fields{}, tmsp.ErrBaseInvalidInput},
		{"invalidAddress", fields{Address: []byte("")}, tmsp.ErrBaseInvalidInput},
		{"invalidSignature", fields{Address: randBytes(), EntityID: genID()}, tmsp.ErrBaseInvalidSignature},
		{"invalidEntityID", fields{Address: randBytes(), EntityID: "", Signature: genSig()}, tmsp.ErrBaseInvalidInput},
		{"invalidEntityType", fields{Address: randBytes(), EntityID: genID(), Signature: genSig(), Type: byte(0xFF)}, tmsp.ErrBaseInvalidInput},
		{"valid", fields{randBytes(), genID(), byte(0xFF), "", genSig()}, tmsp.ErrBaseInvalidInput},
		//		{"valid", fields{randBytes(), uuid.NewV4().String(), crypto.GenPrivKeyEd25519().Sign(crypto.CRandBytes(20))}, tmsp.OK},
	}
	for _, tt := range tests {
		tx := &CreateLegalEntityTx{
			Address:   tt.fields.Address,
			EntityID:  tt.fields.EntityID,
			Type:      tt.fields.Type,
			Name:      tt.fields.Name,
			Signature: tt.fields.Signature,
		}
		if got := tx.ValidateBasic(); got.Code != tt.want.Code {
			t.Errorf("%q. CreateLegalEntityTx.ValidateBasic() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCreateLegalEntityTx_String(t *testing.T) {
	type fields struct {
		Address   []byte
		EntityID  string
		Type      byte
		Name      string
		Signature crypto.Signature
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty", fields{}, "CreateLegalEntityTx{,\"\",0,}"},
		{"stringRepr", fields{Address: []byte{0}, EntityID: "entity_id"}, "CreateLegalEntityTx{00,\"entity_id\",0,}"},
	}
	for _, tt := range tests {
		tx := &CreateLegalEntityTx{
			Address:   tt.fields.Address,
			EntityID:  tt.fields.EntityID,
			Type:      tt.fields.Type,
			Name:      tt.fields.Name,
			Signature: tt.fields.Signature,
		}
		if got := tx.String(); got != tt.want {
			t.Errorf("%q. CreateLegalEntityTx.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
