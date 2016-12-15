package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output version information and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ledgerctl -- HEAD")
	},
}
