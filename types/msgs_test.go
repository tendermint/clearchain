package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

func TestDepositMsg_ValidateBasic(t *testing.T) {
	type fields struct {
		Sender    crypto.Address
		Recipient crypto.Address
		Amount    sdk.Coin
	}

	coin := sdk.Coin{Amount: 100, Denom: "ATM"}
	coinNegative := sdk.Coin{Amount: -100, Denom: "ATM"}

	short := crypto.Address("foo")
	long := crypto.Address("hefkuhwqekufghwqekufgwqekufgkwuqgfkugfkuwgek")
	addr := crypto.GenPrivKeyEd25519().PubKey().Address()
	addr2 := crypto.GenPrivKeyEd25519().PubKey().Address()

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
			"missing address",
			fields{Amount: coin},
			CodeInvalidAddress,
		},
		{
			"short address",
			fields{Amount: coin, Sender: short, Recipient: short},
			CodeInvalidAddress,
		},
		{
			"long address",
			fields{Amount: coin, Sender: long, Recipient: long},
			CodeInvalidAddress,
		},
		{
			"long address2",
			fields{Amount: coin, Sender: addr, Recipient: long},
			CodeInvalidAddress,
		},
		{
			"same address",
			fields{Amount: coin, Sender: addr, Recipient: addr},
			CodeInvalidAddress,
		},
		{
			"proper address",
			fields{Amount: coin, Sender: addr, Recipient: addr2},
			sdk.CodeOK,
		},
		{
			"negative amount",
			fields{Amount: coinNegative, Sender: addr, Recipient: addr2},
			CodeInvalidAmount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DepositMsg{
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
	type fields struct {
		Sender    crypto.Address
		Recipient crypto.Address
		Amount    sdk.Coin
	}

	coin := sdk.Coin{Amount: 100, Denom: "ATM"}
	coinNegative := sdk.Coin{Amount: -100, Denom: "ATM"}

	short := crypto.Address("foo")
	long := crypto.Address("hefkuhwqekufghwqekufgwqekufgkwuqgfkugfkuwgek")
	addr := crypto.GenPrivKeyEd25519().PubKey().Address()
	addr2 := crypto.GenPrivKeyEd25519().PubKey().Address()

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
			"missing address",
			fields{Amount: coin},
			CodeInvalidAddress,
		},
		{
			"short address",
			fields{Amount: coin, Sender: short, Recipient: short},
			CodeInvalidAddress,
		},
		{
			"long address",
			fields{Amount: coin, Sender: long, Recipient: long},
			CodeInvalidAddress,
		},
		{
			"long address2",
			fields{Amount: coin, Sender: addr, Recipient: long},
			CodeInvalidAddress,
		},
		{
			"same address",
			fields{Amount: coin, Sender: addr, Recipient: addr},
			CodeInvalidAddress,
		},
		{
			"proper address",
			fields{Amount: coin, Sender: addr, Recipient: addr2},
			sdk.CodeOK,
		},
		{
			"proper negative amount",
			fields{Amount: coinNegative, Sender: addr, Recipient: addr2},
			sdk.CodeOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := SettleMsg{
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
	type fields struct {
		Sender    crypto.Address
		Recipient crypto.Address
		Operator  crypto.Address
		Amount    sdk.Coin
	}

	coin := sdk.Coin{Amount: 100, Denom: "ATM"}
	coinNegative := sdk.Coin{Amount: -100, Denom: "ATM"}

	short := crypto.Address("foo")
	long := crypto.Address("hefkuhwqekufghwqekufgwqekufgkwuqgfkugfkuwgek")
	addr := crypto.GenPrivKeyEd25519().PubKey().Address()
	addr2 := crypto.GenPrivKeyEd25519().PubKey().Address()
	addr3 := crypto.GenPrivKeyEd25519().PubKey().Address()

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
			"missing address",
			fields{Amount: coin},
			CodeInvalidAddress,
		},
		{
			"short address",
			fields{Amount: coin, Sender: short, Recipient: short},
			CodeInvalidAddress,
		},
		{
			"long address",
			fields{Amount: coin, Sender: long, Recipient: long},
			CodeInvalidAddress,
		},
		{
			"long address2",
			fields{Amount: coin, Sender: addr, Recipient: long},
			CodeInvalidAddress,
		},
		{
			"same address",
			fields{Amount: coin, Sender: addr, Recipient: addr, Operator: addr3},
			CodeInvalidAddress,
		},
		{
			"missing proper address",
			fields{Amount: coin, Sender: addr, Recipient: addr2},
			CodeInvalidAddress,
		},
		{
			"negative amount",
			fields{Amount: coinNegative, Sender: addr, Recipient: addr2},
			CodeInvalidAmount,
		},
		{
			"proper address",
			fields{Amount: coin, Sender: addr, Recipient: addr2, Operator: addr3},
			sdk.CodeOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WithdrawMsg{
				Sender:    tt.fields.Sender,
				Recipient: tt.fields.Recipient,
				Operator:  tt.fields.Operator,
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

func TestCreateAccountMsg_ValidateBasic(t *testing.T) {
	type fields struct {
		Creator         crypto.Address
		PubKey          crypto.PubKey
		AccountType     string
		LegalEntityName string
	}

	creatorAddress := crypto.GenPrivKeyEd25519().PubKey().Address()
	newPubKey := crypto.GenPrivKeyEd25519().PubKey()
	entity := "entity"

	tests := []struct {
		name      string
		fields    fields
		errorCode sdk.CodeType
	}{
		{
			"new CH acc ok",
			fields{Creator: creatorAddress, PubKey: newPubKey, AccountType: EntityClearingHouse, LegalEntityName: entity},
			sdk.CodeOK,
		},
		{
			"new CUS acc ok",
			fields{Creator: creatorAddress, PubKey: newPubKey, AccountType: EntityCustodian, LegalEntityName: entity},
			sdk.CodeOK,
		},
		{
			"new GCM acc ok",
			fields{Creator: creatorAddress, PubKey: newPubKey, AccountType: EntityGeneralClearingMember, LegalEntityName: entity},
			sdk.CodeOK,
		},
		{
			"new ICM acc ok",
			fields{Creator: creatorAddress, PubKey: newPubKey, AccountType: EntityIndividualClearingMember, LegalEntityName: entity},
			sdk.CodeOK,
		},
		{
			"legal entity name is empty",
			fields{Creator: creatorAddress, PubKey: newPubKey, AccountType: EntityIndividualClearingMember, LegalEntityName: "    "},
			CodeInvalidAccount,
		},
		{
			"creator addr is null",
			fields{Creator: nil, PubKey: newPubKey, AccountType: EntityIndividualClearingMember},
			CodeInvalidAddress,
		},
		{
			"invalid len creator address",
			fields{Creator: crypto.Address("short"), PubKey: newPubKey, AccountType: EntityIndividualClearingMember},
			CodeInvalidAddress,
		},
		{
			"new acc pubkey is null",
			fields{Creator: creatorAddress, PubKey: nil, AccountType: EntityIndividualClearingMember},
			CodeInvalidPubKey,
		},
		{
			"same creator and new acc",
			fields{Creator: newPubKey.Address(), PubKey: nil, AccountType: EntityIndividualClearingMember},
			CodeInvalidPubKey,
		},
		{
			"creator is not CH",
			fields{Creator: newPubKey.Address(), PubKey: nil, AccountType: EntityIndividualClearingMember},
			CodeInvalidPubKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := CreateAccountMsg{
				Creator:         tt.fields.Creator,
				PubKey:          tt.fields.PubKey,
				AccountType:     tt.fields.AccountType,
				LegalEntityName: tt.fields.LegalEntityName,
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

func TestSignBytes(t *testing.T) {
	type fields struct {
		Sender    crypto.Address
		Recipient crypto.Address
		Operator  crypto.Address
		Amount    sdk.Coin
	}

	coin := sdk.Coin{Amount: 100, Denom: "ATM"}
	coin2 := sdk.Coin{Amount: 300, Denom: "ATM"}
	coin3 := sdk.Coin{Amount: 100, Denom: "RTD"}
	addr := crypto.GenPrivKeyEd25519().PubKey().Address()
	addr2 := crypto.GenPrivKeyEd25519().PubKey().Address()
	addr3 := crypto.GenPrivKeyEd25519().PubKey().Address()

	seen := make(map[string]bool)
	hasSeen := func(key []byte) bool {
		skey := string(key)
		if seen[skey] {
			return true
		}
		seen[skey] = true
		return false
	}

	tests := []struct {
		name     string
		fields   fields
		beenSeen bool
	}{
		{"one", fields{Sender: addr, Recipient: addr2, Operator: addr3, Amount: coin}, false},
		{"one-rep", fields{Sender: addr, Recipient: addr2, Operator: addr3, Amount: coin}, true},
		{"two", fields{Sender: addr2, Recipient: addr, Operator: addr3, Amount: coin}, false},
		{"three", fields{Sender: addr, Recipient: addr2, Operator: addr3, Amount: coin2}, false},
		{"four", fields{Sender: addr, Recipient: addr3, Operator: addr2, Amount: coin}, false},
		{"five", fields{Sender: addr, Recipient: addr2, Operator: addr3, Amount: coin3}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DepositMsg{
				Sender:    tt.fields.Sender,
				Recipient: tt.fields.Recipient,
				Amount:    tt.fields.Amount,
			}
			bz := d.GetSignBytes()
			assert.Equal(t, tt.beenSeen, hasSeen(bz))

			s := SettleMsg{
				Sender:    tt.fields.Sender,
				Recipient: tt.fields.Recipient,
				Amount:    tt.fields.Amount,
			}
			bz = s.GetSignBytes()
			assert.Equal(t, tt.beenSeen, hasSeen(bz))

			w := WithdrawMsg{
				Sender:    tt.fields.Sender,
				Recipient: tt.fields.Recipient,
				Operator:  tt.fields.Operator,
				Amount:    tt.fields.Amount,
			}
			bz = w.GetSignBytes()
			assert.Equal(t, tt.beenSeen, hasSeen(bz))
		})
	}

}
