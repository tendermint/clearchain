package types

// Index defines the operations that can be performed on
// an objects index.
type Index interface {
	Has(s string) bool
	ToStringSlice() []string
	Add(s string)
}

// AccountIndex stores the list of accounts managed on the ledger.
type AccountIndex struct {
	Accounts []string `json:"accounts"`
}

// NewAccountIndex creates a new accounts index
func NewAccountIndex() *AccountIndex {
	return &AccountIndex{Accounts: []string{}}
}

// Has returns whether s is listed in the accounts index.
func (i *AccountIndex) Has(s string) bool {
	for _, t := range i.Accounts {
		if t == s {
			return true
		}
	}
	return false
}

// ToStringSlice returns a string slice representation of the index.
func (i *AccountIndex) ToStringSlice() []string {
	return i.Accounts
}

// Add adds an account to the index, if it's not yet there.
func (i *AccountIndex) Add(s string) {
	if !i.Has(s) {
		i.Accounts = append(i.Accounts, s)
	}
}

//-----------------------------------------

// AccountIndexGetter is implemented by any value that has a GetAccountIndex
type AccountIndexGetter interface {
	GetAccountIndex() *AccountIndex
}

// AccountIndexSetter is implemented by any value that has a SetAccountIndex
type AccountIndexSetter interface {
	SetAccountIndex(i *AccountIndex)
}

// AccountIndexGetterSetter is implemented by any value that has both
// GetAccountIndex and SetAccountIndex
type AccountIndexGetterSetter interface {
	GetAccountIndex() *AccountIndex
	SetAccountIndex(i *AccountIndex)
}
