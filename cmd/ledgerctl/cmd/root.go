package cmd

import (
	"bufio"
	"io"
	"os"

	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "ledgerctl",
	Short: "Query or send commands to the ledger",
	Long:  `Manage, query, and send transactions to the clearchain ledger`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Fprintln(os.Stderr, "Run 'ledgerctl --help' for usage.")
	},
}

// ReadLine reads a single line from a io.Reader.
func ReadLine(rd io.Reader) (string, error) {
	scanner := scanLine(rd)
	return scanner.Text(), scanner.Err()
}

// ReadLineBytes reads a single line from a io.Reader.
func ReadLineBytes(rd io.Reader) ([]byte, error) {
	scanner := scanLine(rd)
	return scanner.Bytes(), scanner.Err()
}

func scanLine(rd io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(rd)
	scanner.Scan()
	return scanner
}
