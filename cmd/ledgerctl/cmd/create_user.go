package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tendermint/clearchain/client"
	"github.com/tendermint/clearchain/types"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
	"os"
	"strconv"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "user_to_json",
		Short: "creates JSON representation for user that could be used in genesis.json",
		Run: func(cmd *cobra.Command, args []string) {
			userName := readParameter("username")
			entityID := readParameter("entityID")
			permissions, err := strconv.Atoi(readParameter("permissions"))
			if err != nil {
				panic(err)
			}

			pubKey, err := crypto.PubKeyFromBytes(client.Decode(readParameter("public Key")))
			if err != nil {
				panic(err)
			}

			user := types.NewUser(pubKey, userName, entityID, types.Perm(permissions))

			fmt.Println(string(wire.JSONBytes(user)))
		},
	})
}

func readParameter(name string) string {
	fmt.Print(name + ": ")

	parameter, err := ReadLine(os.Stdin)
	if err != nil {
		panic(err)
	}
	return parameter
}
