package types

import "fmt"

// Account defines the attributes of an account
type Account struct {
	ID       string   `json:"id"`        // Account's address
	EntityID string   `json:"entity_id"` // Account's owner
	Wallets  []Wallet `json:"wallets"`   // Account's wallets
}

// NewAccount creates a new account.
func NewAccount(id, entityID string) *Account {
	return &Account{
		ID:       id,
		EntityID: entityID,
		Wallets:  make([]Wallet, 0),
	}
}

// Equal provides an equality operator
func (acc *Account) Equal(a *Account) bool {
	if acc != nil && a != nil {
		return acc.ID == a.ID && acc.EntityID == a.EntityID && acc.walletsEqual(a)
	}
	return acc == a
}

func (acc *Account) walletsEqual(a *Account) bool {
	if len(acc.Wallets) != len(a.Wallets) {
		return false
	}
	for i, w := range acc.Wallets {
		if !w.Equal(&a.Wallets[i]) {
			return false
		}
	}
	return true
}

// Copy make a copy of an Account
func (acc *Account) Copy() *Account {
	accCopy := *acc
	return &accCopy
}

func (acc *Account) String() string {
	if acc == nil {
		return "nil-Account"
	}
	return fmt.Sprintf("Account{%s %s}", acc.ID, acc.EntityID)
}

// BelongsTo checks whether an Account belongs to a given LegalEntity
func (acc *Account) BelongsTo(legalEntityID string) bool {
	return acc.EntityID == legalEntityID
}

// GetWallet retrieves the Account's wallet for the given currency.
func (acc *Account) GetWallet(currency string) *Wallet {
	for i, wal := range acc.Wallets {
		if wal.Currency == currency {
			return &acc.Wallets[i]
		}
	}
	return nil
}

func (account *Account) SetWallet(wallet Wallet) {
	for i, wal := range account.Wallets {
		if wal.Currency == wallet.Currency {
			account.Wallets[i] = wallet
			return
		}
	}
	
	account.Wallets = append(account.Wallets, wallet)
}

//-----------------------------------------

// Wallet defines the attributes of an account's wallet
type Wallet struct {
	Currency string `json:"currency"`
	Balance  int64  `json:"balance"`
	Sequence int    `json:"sequence"`
}

// Equal provides an equality operator
func (w *Wallet) Equal(z *Wallet) bool {
	if w != nil && z != nil {
		return w.Currency == z.Currency && w.Balance == z.Balance && w.Sequence == z.Sequence
	}
	return w == z
}

func (w *Wallet) String() string {
	if w == nil {
		return "nil-Wallet"
	}
	return fmt.Sprintf("Wallet{%s %v %v}", w.Currency, w.Sequence, w.Balance)
}

// AccountsReturned defines the attributes of response's payload
type AccountsReturned struct {
	Account []*Account `json:"accounts"`
}

//-----------------------------------------

// AccountGetter is implemented by any value that has a GetAccount
type AccountGetter interface {
	GetAccount(id string) *Account
}

// AccountSetter is implemented by any value that has a SetAccount
type AccountSetter interface {
	SetAccount(id string, acc *Account)
}

// AccountGetterSetter is implemented by any value that has both
// GetAccount and SetAccount
type AccountGetterSetter interface {
	GetAccount(id string) *Account
	SetAccount(id string, acc *Account)
}
