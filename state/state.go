package state

import (
	"github.com/tendermint/clearchain/types"
	basecoin "github.com/tendermint/basecoin/types"
	common "github.com/tendermint/go-common"
	"github.com/tendermint/go-wire"
)

// State defines the attributes of the system's state
type State struct {
	chainID string
	store   basecoin.KVStore
	cache   *basecoin.KVCache // optional
}

// NewState creates a new State
func NewState(store basecoin.KVStore) *State {
	return &State{
		chainID: "",
		store:   store,
	}
}

// SetChainID sets the State's chain ID
func (s *State) SetChainID(chainID string) {
	s.chainID = chainID
}

// GetChainID retrieves the State's chain ID
func (s *State) GetChainID() string {
	if s.chainID == "" {
		common.PanicSanity("Expected to have set SetChainID")
	}
	return s.chainID
}

// Get retrieves the value for the respective key from the State's store
func (s *State) Get(key []byte) (value []byte) {
	return s.store.Get(key)
}

// Set updates the State's store's key with value
func (s *State) Set(key []byte, value []byte) {
	s.store.Set(key, value)
}

// GetAccount retrieves the Account by address
func (s *State) GetAccount(id string) *types.Account {
	return GetAccount(s.store, id)
}

// SetAccount sets an Account
func (s *State) SetAccount(id string, acc *types.Account) {
	SetAccount(s.store, id, acc)
}

// GetUser retrieves the User by address
func (s *State) GetUser(addr []byte) *types.User {
	return GetUser(s.store, addr)
}

// SetUser sets a User
func (s *State) SetUser(addr []byte, acc *types.User) {
	SetUser(s.store, addr, acc)
}

// GetLegalEntity retrieves the LegalEntity by address
func (s *State) GetLegalEntity(id string) *types.LegalEntity {
	return GetLegalEntity(s.store, id)
}

// SetLegalEntity sets a LegalEntity
func (s *State) SetLegalEntity(id string, l *types.LegalEntity) {
	SetLegalEntity(s.store, id, l)
}

// GetAccountIndex gets the accounts index
func (s *State) GetAccountIndex() *types.AccountIndex {
	return GetAccountIndex(s.store)
}

// SetAccountIndex sets the accounts index
func (s *State) SetAccountIndex(i *types.AccountIndex) {
	SetAccountIndex(s.store, i)
}

func (s *State) GetLegalEntityIndex() *types.LegalEntityIndex {
	data := s.store.Get(LegalEntityIndexKey())
	if len(data) == 0 {
		return nil
	}
	var LegalEntityIndex *types.LegalEntityIndex
	err := wire.ReadBinaryBytes(data, &LegalEntityIndex)
	if err != nil {
		panic(common.Fmt("Error reading LegalEntityIndex %X error: %v",
			data, err.Error()))
	}
	return LegalEntityIndex
}

func (s *State) SetLegalEntityIndex(LegalEntityIndex *types.LegalEntityIndex) {
	LegalEntityIndexBytes := wire.BinaryBytes(LegalEntityIndex)
	s.store.Set(LegalEntityIndexKey(), LegalEntityIndexBytes)
}


//----------------------------------------

func (s *State) CacheWrap() *State {
	cache := basecoin.NewKVCache(s.store)
	return &State{
		chainID: s.chainID,
		store:   cache,
		cache:   cache,
	}
}

// NOTE: errors if s is not from CacheWrap()
func (s *State) CacheSync() {
	s.cache.Sync()
}

//----------------------------------------

// AccountKey generates a data store's unique key for an Account
func AccountKey(id string) []byte {
	return append([]byte("base/a/"), id...)
}

// GetAccount retrieves an Account from the given store
func GetAccount(store basecoin.KVStore, id string) *types.Account {
	data := store.Get(AccountKey(id))
	if len(data) == 0 {
		return nil
	}
	var acc *types.Account
	err := wire.ReadBinaryBytes(data, &acc)
	if err != nil {
		panic(common.Fmt("Error reading account %X error: %v",
			data, err.Error()))
	}
	return acc
}

// SetAccount stores an Account to the given store
func SetAccount(store basecoin.KVStore, id string, acc *types.Account) {
	accBytes := wire.BinaryBytes(acc)
	store.Set(AccountKey(id), accBytes)
}

//----------------------------------------

// UserKey generates a data store's unique key for a User
func UserKey(addr []byte) []byte {
	return append([]byte("base/u/"), addr...)
}

// GetUser retrieves a User from the given store
func GetUser(store basecoin.KVStore, addr []byte) *types.User {
	data := store.Get(UserKey(addr))
	if len(data) == 0 {
		return nil
	}
	var usr *types.User
	err := wire.ReadBinaryBytes(data, &usr)
	if err != nil {
		panic(common.Fmt("Error reading user %X error: %v",
			data, err.Error()))
	}
	return usr
}

// SetUser stores a User to the given store
func SetUser(store basecoin.KVStore, addr []byte, usr *types.User) {
	usrBytes := wire.BinaryBytes(usr)
	store.Set(UserKey(addr), usrBytes)
}

//----------------------------------------

// LegalEntityKey generates a data store's unique key for a LegalEntity
func LegalEntityKey(id string) []byte {
	return append([]byte("base/e/"), id...)
}

// GetLegalEntity retrieves a LegalEntity from the given store
func GetLegalEntity(store basecoin.KVStore, id string) *types.LegalEntity {
	data := store.Get(LegalEntityKey(id))
	if len(data) == 0 {
		return nil
	}
	var ent *types.LegalEntity
	err := wire.ReadBinaryBytes(data, &ent)
	if err != nil {
		panic(common.Fmt("Error reading legal entity %X error: %v",
			data, err.Error()))
	}
	return ent
}

// SetLegalEntity stores a LegalEntity to the given store
func SetLegalEntity(store basecoin.KVStore, id string, ent *types.LegalEntity) {
	entBytes := wire.BinaryBytes(ent)
	store.Set(LegalEntityKey(id), entBytes)
}

//----------------------------------------

// AccountIndexKey generates a data store's unique key for an AccountIndex
func AccountIndexKey() []byte {
	return []byte("base/i/a")
}

func LegalEntityIndexKey() []byte {
	return []byte("base/i/l")
}


// GetAccountIndex retrieves a AccountIndex from the given store
func GetAccountIndex(store basecoin.KVStore) *types.AccountIndex {
	data := store.Get(AccountIndexKey())
	if len(data) == 0 {
		return nil
	}
	var acc *types.AccountIndex
	err := wire.ReadBinaryBytes(data, &acc)
	if err != nil {
		panic(common.Fmt("Error reading account index %X error: %v",
			data, err.Error()))
	}
	return acc
}

// SetAccountIndex stores a AccountIndex to the given store
func SetAccountIndex(store basecoin.KVStore, acc *types.AccountIndex) {
	accBytes := wire.BinaryBytes(acc)
	store.Set(AccountIndexKey(), accBytes)
}

func SetLegalEntityIndex(store basecoin.KVStore, legalEntityIndex *types.LegalEntityIndex) {
	bytes := wire.BinaryBytes(legalEntityIndex)
	store.Set(LegalEntityIndexKey(), bytes)
}

