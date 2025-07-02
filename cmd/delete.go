package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tryoasnafi/gate/internal/session"
	"github.com/tryoasnafi/gate/internal/store.go"
)

func DeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete an SSH label",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pwd := session.Require()
			if err := store.DeleteEntry(pwd, args[0]); err != nil {
				fmt.Println("Error deleting:", err)
				os.Exit(1)
			}
			fmt.Println("Label deleted.")
		},
	}
}
