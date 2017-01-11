package testutil

import (
	"math/rand"

	"reflect"

	"github.com/tendermint/clearchain/types"
	"github.com/dustinkirkland/golang-petname"
	"github.com/satori/go.uuid"
	"github.com/tendermint/go-common"
	"github.com/tendermint/go-crypto"
)

// RandCH creates a new CH
func RandCH() *types.LegalEntity {
	return types.NewCH(uuid.NewV4().String(), petname.Generate(2, "-"), nil, "")
}

// RandGCM creates a new GCM
func RandGCM(creatorAddr []byte) *types.LegalEntity {
	return types.NewGCM(uuid.NewV4().String(), petname.Generate(2, "-"), creatorAddr, "")
}

// RandICM creates a new ICM
func RandICM(creatorAddr []byte) *types.LegalEntity {
	return types.NewICM(uuid.NewV4().String(), petname.Generate(2, "-"), creatorAddr, "")
}

// RandCustodian creates a new Custodian
func RandCustodian(creatorAddr []byte) *types.LegalEntity {
	return types.NewCustodian(uuid.NewV4().String(), petname.Generate(2, "-"), creatorAddr, "")
}

// RandEntity generates a LegalEntity
func RandEntity(t byte, permissions types.Perm) *types.LegalEntity {
	return types.NewLegalEntity(uuid.NewV4().String(), t, petname.Generate(2, "-"), permissions, nil, "")
}

// RandAccounts generates num random accounts for the given LegalEntity.
func RandAccounts(num int, e *types.LegalEntity) []*types.Account {
	accounts := make([]*types.Account, num)
	for i := 0; i < num; i++ {
		accounts[i] = RandAccount(e)
	}
	return accounts
}

// RandAccount generate an Account for a given LegalEntity if not nil.
func RandAccount(e *types.LegalEntity) *types.Account {
	if e == nil {
		return types.NewAccount(uuid.NewV4().String(), "")
	}
	return types.NewAccount(uuid.NewV4().String(), e.ID)
}

// PrivUserWithLegalEntityFromSecret generates a PrivUser from
// a secret and populate the inner User with LegalEntity's
// address and the given permissions.
func PrivUserWithLegalEntityFromSecret(secret string, e *types.LegalEntity, perms types.Perm) *types.PrivUser {
	privUser := PrivUserFromSecret(secret)
	privUser.User.Permissions = perms
	privUser.User.EntityID = e.ID
	return privUser
}

// PrivUserFromSecret generates a new PrivUser from a secret
func PrivUserFromSecret(secret string) *types.PrivUser {
	var privKey crypto.PrivKey
	if len(secret) == 0 {
		privKey = crypto.GenPrivKeyEd25519()
	} else {
		privKey = crypto.GenPrivKeyEd25519FromSecret([]byte(secret))
	}

	privUser := types.PrivUser{
		PrivKey: privKey,
		User: types.User{
			PubKey: privKey.PubKey(),
			Name:   petname.Generate(2, "-"),
		},
	}

	return &privUser
}

// RandUsersWithLegalEntity generates num random PrivUsers
// and attach them a LegalEntity
func RandUsersWithLegalEntity(num int, e *types.LegalEntity, perms types.Perm) []*types.PrivUser {
	privUsers := make([]*types.PrivUser, num)
	for i := 0; i < num; i++ {
		privUsers[i] = PrivUserWithLegalEntityFromSecret("", e, perms)
	}
	return privUsers
}

// RandUsers generates num random PrivUsers
func RandUsers(num int) []*types.PrivUser {
	privUsers := make([]*types.PrivUser, num)
	for i := 0; i < num; i++ {
		privUsers[i] = PrivUserFromSecret("")
	}
	return privUsers
}

// RandWallet creates a new Wallet with a random Balance
func RandWallet(currency types.Currency, minBalanceMul, maxBalanceMul int64) types.Wallet {
	balance := common.RandInt64() * currency.MinimumUnit() * (maxBalanceMul - minBalanceMul)
	return types.Wallet{
		Balance: balance, Currency: currency.Symbol()}
}

// RandWallets generates num random Wallets
func RandWallets(num int, minBalanceMul, maxBalanceMul int64) []types.Wallet {
	wallets := make([]types.Wallet, num)
	currenciesKeys := reflect.ValueOf(types.Currencies).MapKeys()

	for i := 0; i < num; i++ {
		randIndex := rand.Intn(len(currenciesKeys))
		ccy := types.Currencies[currenciesKeys[randIndex].String()]
		wallets[i] = RandWallet(ccy, minBalanceMul, maxBalanceMul)
	}
	return wallets
}
