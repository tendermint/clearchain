package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	crypto "github.com/tendermint/go-crypto"
)

func TestDepositMsg_ValidateBasic(t *testing.T) {
	coin := sdk.Coin{Amount: 100, Denom: "ATM"}
	coinNegative := sdk.Coin{Amount: -100, Denom: "ATM"}
	short := crypto.Address("foo")
	long := crypto.Address("hefkuhwqekufghwqekufgwqekufgkwuqgfkugfkuwgek")
	addr := crypto.GenPrivKeyEd25519().PubKey().Address()
	addr2 := crypto.GenPrivKeyEd25519().PubKey().Address()
	addr3 := crypto.GenPrivKeyEd25519().PubKey().Address()

	type fields struct {
		Operator  crypto.Address
		Sender    crypto.Address
		Recipient crypto.Address
		Amount    sdk.Coin
	}
	tests := []struct {
		name      string
		fields    fields
		errorCode sdk.CodeType
	}{
		{
			"empty msg",
			fields{},
			CodeInvalidAmount,
		},
		{
			"no denom",
			fields{Amount: sdk.Coin{Amount: 100}},
			CodeInvalidAmount,
		},
		{
			"no amount",
			fields{Amount: sdk.Coin{Denom: "Foo"}},
			CodeInvalidAmount,
		},
		{
			"missing operator address",
			fields{Amount: coin},
			CodeInvalidAddress,
		},
		{
			"missing sender address",
			fields{Amount: coin, Operator: addr},
			CodeInvalidAddress,
		},
		{
			"short address",
			fields{Amount: coin, Operator: short},
			CodeInvalidAddress,
		},
		{
			"long sender address",
			fields{Amount: coin, Operator: addr, Sender: long},
			CodeInvalidAddress,
		},
		{
			"long recipient address",
			fields{Amount: coin, Operator: addr, Sender: addr, Recipient: long},
			CodeInvalidAddress,
		},
		{
			"same address",
			fields{Amount: coin, Operator: addr2, Sender: addr, Recipient: addr},
			CodeInvalidAddress,
		},
		{
			"proper addresses",
			fields{Amount: coin, Operator: addr, Sender: addr2, Recipient: addr3},
			sdk.CodeOK,
		},
		{
			"negative amount",
			fields{Amount: coinNegative, Operator: addr, Sender: addr2, Recipient: addr3},
			CodeInvalidAmount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DepositMsg{
				Operator:  tt.fields.Operator,
				Sender:    tt.fields.Sender,
				Recipient: tt.fields.Recipient,
				Amount:    tt.fields.Amount,
			}
			got := d.ValidateBasic()
			if got == nil {
				assert.True(t, tt.errorCode.IsOK())
			} else {
				assert.Equal(t, tt.errorCode, got.ABCICode())
			}
		})
	}
}
func TestSettleMsg_ValidateBasic(t *testing.T) {
	coin := sdk.Coin{Amount: 100, Denom: "ATM"}
	coinNegative := sdk.Coin{Amount: -100, Denom: "ATM"}
	short := crypto.Address("foo")
	long := crypto.Address("hefkuhwqekufghwqekufgwqekufgkwuqgfkugfkuwgek")
	addr := crypto.GenPrivKeyEd25519().PubKey().Address()
	addr2 := crypto.GenPrivKeyEd25519().PubKey().Address()
	addr3 := crypto.GenPrivKeyEd25519().PubKey().Address()

	type fields struct {
		Operator  crypto.Address
		Sender    crypto.Address
		Recipient crypto.Address
		Amount    sdk.Coin
	}
	tests := []struct {
		name      string
		fields    fields
		errorCode sdk.CodeType
	}{
		{
			"empty msg",
			fields{},
			CodeInvalidAmount,
		},
		{
			"no denom",
			fields{Amount: sdk.Coin{Amount: 100}},
			CodeInvalidAmount,
		},
		{
			"no amount",
			fields{Amount: sdk.Coin{Denom: "Foo"}},
			CodeInvalidAmount,
		},
		{
			"missing operator address",
			fields{Amount: coin},
			CodeInvalidAddress,
		},
		{
			"short address",
			fields{Amount: coin, Operator: short, Sender: short, Recipient: short},
			CodeInvalidAddress,
		},
		{
			"long address",
			fields{Amount: coin, Operator: long, Sender: short, Recipient: long},
			CodeInvalidAddress,
		},
		{
			"long address2",
			fields{Amount: coin, Operator: addr, Sender: addr2, Recipient: long},
			CodeInvalidAddress,
		},
		{
			"sender and recipient got same address",
			fields{Amount: coin, Operator: addr, Sender: addr2, Recipient: addr2},
			CodeInvalidAddress,
		},
		{
			"proper address",
			fields{Amount: coin, Operator: addr, Sender: addr2, Recipient: addr3},
			sdk.CodeOK,
		},
		{
			"proper negative amount",
			fields{Amount: coinNegative, Operator: addr3, Sender: addr, Recipient: addr2},
			sdk.CodeOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := SettleMsg{
				Operator:  tt.fields.Operator,
				Sender:    tt.fields.Sender,
				Recipient: tt.fields.Recipient,
				Amount:    tt.fields.Amount,
			}
			got := d.ValidateBasic()
			if got == nil {
				assert.True(t, tt.errorCode.IsOK())
			} else {
				assert.Equal(t, tt.errorCode, got.ABCICode())
			}
		})
	}
}
func TestWithdrawMsg_ValidateBasic(t *testing.T) {
	coin := sdk.Coin{Amount: 100, Denom: "ATM"}
	coinNegative := sdk.Coin{Amount: -100, Denom: "ATM"}
	short := crypto.Address("foo")
	long := crypto.Address("hefkuhwqekufghwqekufgwqekufgkwuqgfkugfkuwgek")
	addr := crypto.GenPrivKeyEd25519().PubKey().Address()
	addr2 := crypto.GenPrivKeyEd25519().PubKey().Address()
	addr3 := crypto.GenPrivKeyEd25519().PubKey().Address()

	type fields struct {
		Sender    crypto.Address
		Recipient crypto.Address
		Operator  crypto.Address
		Amount    sdk.Coin
	}
	tests := []struct {
		name      string
		fields    fields
		errorCode sdk.CodeType
	}{
		{"empty msg", fields{}, CodeInvalidAmount},
		{"no denom", fields{Amount: sdk.Coin{Amount: 100}}, CodeInvalidAmount},
		{"no amount", fields{Amount: sdk.Coin{Denom: "Foo"}}, CodeInvalidAmount},
		{"missing address", fields{Amount: coin}, CodeInvalidAddress},
		{"short address", fields{Amount: coin, Sender: short, Recipient: short}, CodeInvalidAddress},
		{"long address", fields{Amount: coin, Sender: long, Recipient: long}, CodeInvalidAddress},
		{"long address2", fields{Amount: coin, Sender: addr, Recipient: long}, CodeInvalidAddress},
		{"same address", fields{Amount: coin, Sender: addr, Recipient: addr, Operator: addr3}, CodeInvalidAddress},
		{"missing proper address", fields{Amount: coin, Sender: addr, Recipient: addr2}, CodeInvalidAddress},
		{"negative amount", fields{Amount: coinNegative, Sender: addr, Recipient: addr2}, CodeInvalidAmount},
		{"proper address", fields{Amount: coin, Sender: addr, Recipient: addr2, Operator: addr3}, sdk.CodeOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WithdrawMsg{
				Operator:  tt.fields.Operator,
				Sender:    tt.fields.Sender,
				Recipient: tt.fields.Recipient,
				Amount:    tt.fields.Amount,
			}
			got := w.ValidateBasic()
			if got == nil {
				assert.True(t, tt.errorCode.IsOK())
			} else {
				assert.Equal(t, tt.errorCode, got.ABCICode(), got.Error())
			}
		})
	}
}

func TestCreateUserAccountMsg_ValidateBasic(t *testing.T) {
	creatorAddress := crypto.GenPrivKeyEd25519().PubKey().Address()
	newPubKey := crypto.GenPrivKeyEd25519().PubKey()
	entity := "entity"
	type fields struct {
		Creator         crypto.Address
		PubKey          crypto.PubKey
		LegalEntityType string
		LegalEntityName string
		IsAdmin         bool
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.CodeType
	}{
		{"new CH acc ok", fields{creatorAddress, newPubKey, EntityClearingHouse, entity, false}, sdk.CodeOK},
		{"new CUS acc ok", fields{creatorAddress, newPubKey, EntityCustodian, entity, false}, sdk.CodeOK},
		{"new GCM acc ok", fields{creatorAddress, newPubKey, EntityGeneralClearingMember, entity, true}, sdk.CodeOK},
		{"new ICM acc ok", fields{creatorAddress, newPubKey, EntityIndividualClearingMember, entity, true}, sdk.CodeOK},
		{"legal entity name is empty", fields{creatorAddress, newPubKey, EntityIndividualClearingMember, "", true}, CodeInvalidAccount},
		{"wrong legal entity type", fields{creatorAddress, newPubKey, "invalid", entity, true}, CodeInvalidAccount},
		{"creator is nil", fields{nil, newPubKey, EntityIndividualClearingMember, entity, true}, CodeInvalidAddress},
		{"invalid creator len", fields{crypto.Address("short"), newPubKey, EntityIndividualClearingMember, entity, true}, CodeInvalidAddress},
		{"new pubkey is nil", fields{creatorAddress, nil, EntityIndividualClearingMember, entity, true}, CodeInvalidPubKey},
		{"same creator and acct", fields{newPubKey.Address(), newPubKey, EntityIndividualClearingMember, entity, true}, CodeInvalidPubKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := CreateUserAccountMsg{
				Creator:         tt.fields.Creator,
				PubKey:          tt.fields.PubKey,
				LegalEntityType: tt.fields.LegalEntityType,
				LegalEntityName: tt.fields.LegalEntityName,
				IsAdmin:         tt.fields.IsAdmin,
			}
			got := msg.ValidateBasic()
			if got == nil {
				assert.True(t, tt.want == sdk.CodeOK)
			} else {
				assert.Equal(t, tt.want, got.ABCICode(), got.Error())
			}
		})
	}
}
func TestCreateAssetAccountMsg_ValidateBasic(t *testing.T) {
	creatorAddress := crypto.GenPrivKeyEd25519().PubKey().Address()
	newPubKey := crypto.GenPrivKeyEd25519().PubKey()
	type fields struct {
		Creator crypto.Address
		PubKey  crypto.PubKey
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.CodeType
	}{
		{"new CH acc ok", fields{creatorAddress, newPubKey}, sdk.CodeOK},
		{"new CUS acc ok", fields{creatorAddress, newPubKey}, sdk.CodeOK},
		{"new GCM acc ok", fields{creatorAddress, newPubKey}, sdk.CodeOK},
		{"new ICM acc ok", fields{creatorAddress, newPubKey}, sdk.CodeOK},
		{"creator is nil", fields{nil, newPubKey}, CodeInvalidAddress},
		{"invalid creator len", fields{crypto.Address("short"), newPubKey}, CodeInvalidAddress},
		{"new pubkey is nil", fields{creatorAddress, nil}, CodeInvalidPubKey},
		{"same creator and acct", fields{newPubKey.Address(), newPubKey}, CodeInvalidPubKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := CreateAssetAccountMsg{
				Creator: tt.fields.Creator,
				PubKey:  tt.fields.PubKey,
			}
			got := msg.ValidateBasic()
			if got == nil {
				assert.True(t, tt.want == sdk.CodeOK)
			} else {
				assert.Equal(t, tt.want, got.ABCICode(), got.Error())
			}
		})
	}
}
