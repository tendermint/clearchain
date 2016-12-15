package types

import (
	"github.com/satori/go.uuid"
	common "github.com/tendermint/go-common"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
	tmsp "github.com/tendermint/tmsp/types"
)

const (
	// TxTypeCreateAccount defines CreateAccountTx's code
	TxTypeCreateAccount = byte(0x02)
)

// CreateAccountTx defines the attributes of an account create.
type CreateAccountTx struct {
	Address   []byte           `json:"address"`    // Hash of the user's PubKey
	AccountID string           `json:"account_id"` // ID of the new account
	Signature crypto.Signature `json:"signature"`
}

func (tx *CreateAccountTx) SignTx(privateKey crypto.PrivKey, chainID string) {
	tx.Signature = privateKey.Sign(tx.SignBytes(chainID))
}

// TxType returns the byte type of CreateAccountTx
func (tx *CreateAccountTx) TxType() byte {
	return TxTypeCreateAccount
}

// SignBytes generates a byte-to-byte signature
func (tx *CreateAccountTx) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sig := tx.Signature
	tx.Signature = nil
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	tx.Signature = sig
	return signBytes
}

// ValidateBasic performs basic validation on the Tx.
func (tx *CreateAccountTx) ValidateBasic() tmsp.Result {
	if len(tx.Address) != 20 {
		return tmsp.ErrBaseInvalidInput.AppendLog("Invalid address length")
	}
	if tx.Signature == nil {
		return tmsp.ErrBaseInvalidSignature.AppendLog("The transaction must be signed")
	}
	if _, err := uuid.FromString(tx.AccountID); err != nil {
		return tmsp.ErrBaseInvalidInput.AppendLog(common.Fmt("Invalid account_id: %s", err))
	}
	return tmsp.OK
}

func (tx *CreateAccountTx) String() string {
	return common.Fmt("CreateAccountTx{%x,%q}", tx.Address, tx.AccountID)
}
