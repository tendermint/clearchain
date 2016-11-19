package types

import "testing"

func TestAccount_Equal(t *testing.T) {
	type fields struct {
		ID       string
		EntityID string
		Wallets  []Wallet
	}
	type args struct {
		a *Account
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"empty", fields{}, args{&Account{}}, true},
		{"notEqual", fields{"address", "entity", nil}, args{&Account{}}, false},
		{"equal", fields{"address", "entity", []Wallet{}},
			args{&Account{"address", "entity", []Wallet{}}}, true},
	}
	for _, tt := range tests {
		acc := &Account{
			ID:       tt.fields.ID,
			EntityID: tt.fields.EntityID,
		}
		if got := acc.Equal(tt.args.a); got != tt.want {
			t.Errorf("%q. Account.Equal() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestAccount_Copy(t *testing.T) {
	type fields struct {
		ID       string
		EntityID string
		Wallets  []Wallet
	}
	tests := []struct {
		name   string
		fields fields
		want   *Account
	}{
		{
			"copy", fields{"address", "entity", []Wallet{}}, &Account{"address", "entity", []Wallet{}},
		},
	}
	for _, tt := range tests {
		acc := &Account{
			ID:       tt.fields.ID,
			EntityID: tt.fields.EntityID,
		}
		if got := acc.Copy(); !got.Equal(tt.want) {
			t.Errorf("%q. Account.Copy() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestAccount_String(t *testing.T) {
	type fields struct {
		ID       string
		EntityID string
	}
	tests := []struct {
		name   string
		fields *fields
		want   string
	}{
		{"nil", nil, "nil-Account"},
		{"empty", &fields{}, "Account{ }"},
		{"nonEmpty", &fields{"ID", "entityID"}, "Account{ID entityID}"},
	}
	for _, tt := range tests {
		var acc *Account
		if tt.fields != nil {
			acc = &Account{
				ID:       tt.fields.ID,
				EntityID: tt.fields.EntityID,
			}
		}
		if got := acc.String(); got != tt.want {
			t.Errorf("%q. Account.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestAccount_BelongsTo(t *testing.T) {
	type fields struct {
		ID       string
		EntityID string
	}
	type args struct {
		legalEntityID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"emptyFields", fields{}, args{}, true},
		{"emptyEntity", fields{EntityID: "entity"}, args{}, false},
		{"doesNotBelongToEntity", fields{"ID", "otherEntityID"}, args{"entityID"}, false},
		{"belongsToEntity", fields{"ID", "entityID"}, args{"entityID"}, true},
	}
	for _, tt := range tests {
		acc := &Account{
			ID:       tt.fields.ID,
			EntityID: tt.fields.EntityID,
		}
		if got := acc.BelongsTo(tt.args.legalEntityID); got != tt.want {
			t.Errorf("%q. Account.BelongsTo() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestAccount_GetWallet(t *testing.T) {
	type fields struct {
		ID       string
		EntityID string
		Wallets  []Wallet
	}
	type args struct {
		currency string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Wallet
	}{
		{"nil", fields{}, args{"USD"}, nil},
		{"nil", fields{Wallets: []Wallet{Wallet{Currency: "USD", Balance: 10}}}, args{"USD"}, &Wallet{Currency: "USD", Balance: 10}},
		// {"notEqual", fields{"address", "entity", nil}, args{&Account{}}, false},
		// {"equal", fields{"address", "entity", []Wallet{}},
		// 	args{&Account{"address", "entity", []Wallet{}}}, true},
	}
	for _, tt := range tests {
		acc := &Account{
			ID:       tt.fields.ID,
			EntityID: tt.fields.EntityID,
			Wallets:  tt.fields.Wallets,
		}
		if got := acc.GetWallet(tt.args.currency); !got.Equal(tt.want) {
			t.Errorf("%q. Account.GetWallet() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestWallet_Equal(t *testing.T) {
	type fields struct {
		Currency string
		Balance  int64
		Sequence int
	}
	type args struct {
		z *Wallet
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"equal", fields{"USD", 10, 1}, args{&Wallet{Currency: "USD", Balance: 10, Sequence: 1}}, true},
		{"notEqual", fields{"USD", 10, 1}, args{&Wallet{}}, false},
	}
	for _, tt := range tests {
		w := &Wallet{
			Currency: tt.fields.Currency,
			Balance:  tt.fields.Balance,
			Sequence: tt.fields.Sequence,
		}
		if got := w.Equal(tt.args.z); got != tt.want {
			t.Errorf("%q. Wallet.Equal() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestWallet_String(t *testing.T) {
	type fields struct {
		Currency string
		Balance  int64
		Sequence int
	}
	tests := []struct {
		name   string
		fields *fields
		want   string
	}{
		{"nil", nil, "nil-Wallet"},
		{"nonEmpty", &fields{"USD", 10, 1}, "Wallet{USD 1 10}"},
		{"empty", &fields{}, "Wallet{ 0 0}"},
	}
	for _, tt := range tests {
		var w *Wallet
		if tt.fields != nil {
			w = &Wallet{
				Currency: tt.fields.Currency,
				Balance:  tt.fields.Balance,
				Sequence: tt.fields.Sequence,
			}
		}
		if got := w.String(); got != tt.want {
			t.Errorf("%q. Wallet.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
