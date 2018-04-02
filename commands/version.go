package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/clearchain"
)

var (
	// VersionCmd prints the program's version to stderr and exits.
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print clearchain's version",
		Run:   doVersionCmd,
	}
)

func doVersionCmd(cmd *cobra.Command, args []string) {
	v := clearchain.Version
	if len(v) == 0 {
		fmt.Fprintln(os.Stderr, "unset")
		return
	}
	fmt.Fprintln(os.Stderr, v)
}
