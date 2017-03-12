package types

import (
	"bytes"
	"errors"

	common "github.com/tendermint/go-common"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
)

// Tx (Transaction) is an atomic operation on the ledger state.
type Tx interface {
	TxType() byte
	SignBytes(chainID string) []byte
}

// SignedTx extends Tx with a method to generate signatures.
type SignedTx interface {
	TxType() byte
	SignBytes(chainID string) []byte
	SignTx(privateKey crypto.PrivKey, chainID string) error
}

// TxExecutor validates Tx execution permission
type TxExecutor interface {
	CanExecTx(byte) bool
}

// CanExecTx is a convenience function that validates
// caller's execution permission on a Tx.
func CanExecTx(executor TxExecutor, tx Tx) bool {
	return executor.CanExecTx(tx.TxType())
}

var _ = wire.RegisterInterface(
	struct{ Tx }{},
	wire.ConcreteType{O: &TransferTx{}, Byte: TxTypeTransfer},
	wire.ConcreteType{O: &BaseQueryTx{}, Byte: TxTypeQueryBase},
	wire.ConcreteType{O: &ObjectsQueryTx{}, Byte: TxTypeQueryObjects},
	wire.ConcreteType{O: &CreateAccountTx{}, Byte: TxTypeCreateAccount},
	wire.ConcreteType{O: &CreateLegalEntityTx{}, Byte: TxTypeCreateLegalEntity},
	wire.ConcreteType{O: &CreateUserTx{}, Byte: TxTypeCreateUser},
)

// SignTx signs the transaction if its address and the privateKey's one match.
func SignTx(signedBytes []byte, addr []byte, privKey crypto.PrivKey) (crypto.Signature, error) {
	if !bytes.Equal(privKey.PubKey().Address(), addr) {
		return nil, errors.New(common.Fmt("SignTx: addresses mismatch: %x != %x",
			privKey.PubKey().Address(), addr))
	}
	return privKey.Sign(signedBytes), nil
}
