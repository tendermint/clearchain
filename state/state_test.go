package state

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	bscoin "github.com/tendermint/basecoin/types"
	"github.com/tendermint/clearchain/types"
)

func TestNewState(t *testing.T) {
	s := NewState(bscoin.NewMemKVStore())
	if s == nil {
		t.Error("NewState() return nil")
	}
}

func TestGetChainID(t *testing.T) {
	s := NewState(bscoin.NewMemKVStore())
	defer func() {
		if r := recover(); r == nil {
			t.Error("GetChainID() was expected to panic")
		}
	}()
	s.GetChainID()
}

func TestSetChainID(t *testing.T) {
	s := NewState(bscoin.NewMemKVStore())
	chainID := "test"
	s.SetChainID(chainID)
	if ret := s.GetChainID(); ret != chainID {
		t.Errorf("GetChainID() return %v, expected: %v", ret, chainID)
	}
}

func TestGet(t *testing.T) {
	s := NewState(bscoin.NewMemKVStore())
	if ret := s.Get([]byte("key")); ret != nil {
		t.Errorf("Get() return %v, expected nil", ret)
	}
}

func TestSet(t *testing.T) {
	s := NewState(bscoin.NewMemKVStore())
	key := []byte("key")
	value := []byte("value")
	s.Set(key, value)
	if ret := s.Get(key); string(ret) != string(value) {
		t.Errorf("Get() return %v, expected %v", ret, value)
	}
}

func TestAccountKey(t *testing.T) {
	expected := "base/a/account"
	if ret := AccountKey("account"); string(ret) != expected {
		t.Errorf("AccountKey() return %v, expected %v", ret, expected)
	}
}

func TestUserKey(t *testing.T) {
	addr := "address"
	expected := "base/u/address"
	if ret := UserKey([]byte(addr)); string(ret) != string(expected) {
		t.Errorf("UserKey() return %v, expected %v", ret, expected)
	}
}

func TestLegalEntityKey(t *testing.T) {
	expected := "base/e/entity"
	if ret := LegalEntityKey("entity"); string(ret) != expected {
		t.Errorf("LegalEntityKey() return %v, expected %v", ret, expected)
	}
}

func TestGetAccount(t *testing.T) {
	s := NewState(bscoin.NewMemKVStore())
	acc := &types.Account{ID: uuid.NewV4().String()}
	s.SetAccount(acc.ID, acc)
	if ret := s.GetAccount("nonexisting"); ret != nil {
		t.Errorf("GetAccount() return %v, expected nil", ret)
	}
	if ret := s.GetAccount(acc.ID); ret == nil {
		t.Errorf("GetAccount() return %v, expected: %v", ret, acc)
	}
}

func TestGetUser(t *testing.T) {
	s := NewState(bscoin.NewMemKVStore())
	acc := &types.User{}
	addr := []byte("address")
	s.SetUser(addr, acc)
	if ret := s.GetUser([]byte("nonexisting")); ret != nil {
		t.Errorf("GetUser() return %v, expected nil", ret)
	}
	if ret := s.GetUser(addr); ret == nil {
		t.Errorf("GetUser() return %v, expected: %v", ret, acc)
	}
}

func TestGetLegalEntity(t *testing.T) {
	s := NewState(bscoin.NewMemKVStore())
	e := &types.LegalEntity{ID: uuid.NewV4().String()}
	s.SetLegalEntity(e.ID, e)
	if ret := s.GetLegalEntity("nonexisting"); ret != nil {
		t.Errorf("GetUser() return %v, expected nil", ret)
	}
	if ret := s.GetLegalEntity(e.ID); ret == nil {
		t.Errorf("GetUser() return %v, expected: %v", ret, e)
	}
}
