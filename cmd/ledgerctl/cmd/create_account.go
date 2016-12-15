package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tendermint/clearchain/types"
	"github.com/tendermint/go-wire"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "account_to_json",
		Short: "creates JSON representation for account that could be used in genesis.json",
		Run: func(cmd *cobra.Command, args []string) {
			accountID := readParameter("accountID")
			entityID := readParameter("entityID")

			account:= types.NewAccount(accountID, entityID)

			fmt.Println(string(wire.JSONBytes(account)))
		},
	})
}
