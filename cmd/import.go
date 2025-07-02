package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tryoasnafi/gate/internal/crypto"
	"github.com/tryoasnafi/gate/internal/session"
	"github.com/tryoasnafi/gate/internal/store.go"
)

func ImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import",
		Short: "Import SSH entries from file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pwd := session.Require()

			fmt.Print("Enter import passphrase: ")
			importPwd, _ := crypto.ReadPassword()

			if err := store.ImportEntries(pwd, importPwd, args[0]); err != nil {
				fmt.Println("Import failed:", err)
				os.Exit(1)
			}
			fmt.Println("Import complete.")
		},
	}
}
