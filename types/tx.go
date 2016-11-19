package types

import (
	wire "github.com/tendermint/go-wire"
)

// Tx (Transaction) is an atomic operation on the ledger state.
type Tx interface {
	TxType() byte
	SignBytes(chainID string) []byte
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
	wire.ConcreteType{O: &AccountQueryTx{}, Byte: TxTypeQueryAccount},
	wire.ConcreteType{O: &AccountIndexQueryTx{}, Byte: TxTypeQueryAccountIndex},
	wire.ConcreteType{O: &CreateAccountTx{}, Byte: TxTypeCreateAccount},
	wire.ConcreteType{O: &CreateLegalEntityTx{}, Byte: TxTypeCreateLegalEntity},
	wire.ConcreteType{O: &CreateUserTx{}, Byte: TxTypeCreateUser},
)
