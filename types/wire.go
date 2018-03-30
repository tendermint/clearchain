package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	oldwire "github.com/tendermint/go-wire"
)

const (
	typeDepositMsg            = 0x1
	typeSettleMsg             = 0x2
	typeWithdrawMsg           = 0x3
	typeCreateAdminMsg        = 0x4
	typeCreateOperatorMsg     = 0x5
	typeCreateAssetAccountMsg = 0x6
	typeFreezeAdminMsg        = 0x7
	typeFreezeOperatorMsg     = 0x8

	typeAppAccount = 0x1
)

// MakeCodec instantiate a wire.Codec and register
// all application's types; it returns the new codec.
func MakeCodec() *wire.Codec {
	var _ = oldwire.RegisterInterface(
		struct{ sdk.Msg }{},
		oldwire.ConcreteType{DepositMsg{}, typeDepositMsg},
		oldwire.ConcreteType{SettleMsg{}, typeSettleMsg},
		oldwire.ConcreteType{WithdrawMsg{}, typeWithdrawMsg},
		oldwire.ConcreteType{CreateAdminMsg{}, typeCreateAdminMsg},
		oldwire.ConcreteType{CreateOperatorMsg{}, typeCreateOperatorMsg},
		oldwire.ConcreteType{CreateAssetAccountMsg{}, typeCreateAssetAccountMsg},
		oldwire.ConcreteType{FreezeAdminMsg{}, typeFreezeAdminMsg},
		oldwire.ConcreteType{FreezeOperatorMsg{}, typeFreezeOperatorMsg},
	)
	var _ = oldwire.RegisterInterface(
		struct{ sdk.Account }{},
		oldwire.ConcreteType{&AppAccount{}, typeAppAccount},
	)
	return wire.NewCodec()
}
