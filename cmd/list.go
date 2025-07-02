package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/tryoasnafi/gate/internal/session"
	"github.com/tryoasnafi/gate/internal/store.go"
)

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all SSH labels",
		Run: func(cmd *cobra.Command, args []string) {
			pwd := session.Require()
			items, err := store.ListEntries(pwd)
			if err != nil {
				fmt.Println("Error listing entries:", err)
				os.Exit(1)
			}
			fmt.Printf("List gate(%d):\n", len(items))
			for _, item := range items {
				fmt.Printf("[%s] %s@%s:%d (%s)\n", item.Label, item.User, item.Host, item.Port, item.CreatedAt.Format(time.DateTime))
			}
		},
	}
}
