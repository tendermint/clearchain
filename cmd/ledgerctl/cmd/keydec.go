package cmd

import (
	"os"

	"fmt"

	"log"

	"github.com/spf13/cobra"
	"github.com/tendermint/clearchain/client"
)

func init() {
	RootCmd.AddCommand(keydecCmd)
}

var keydecCmd = &cobra.Command{
	Use:   "keydec",
	Short: "Read a key from stdin and print its binary representation to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		key, err := ReadLine(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(client.Decode(key))
	},
}
