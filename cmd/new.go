package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/tryoasnafi/gate/internal/crypto"
	"github.com/tryoasnafi/gate/internal/session"
	"github.com/tryoasnafi/gate/internal/store.go"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "Add new SSH label",
		Run: func(cmd *cobra.Command, args []string) {
			pwd := session.Require()

			var user, host, label string
			var port int
			fmt.Print("Label: ")
			fmt.Scanln(&label)
			fmt.Print("User: ")
			fmt.Scanln(&user)
			fmt.Print("Host: ")
			fmt.Scanln(&host)
			fmt.Print("Port: ")
			fmt.Scanln(&port)
			fmt.Print("Password: ")
			password, _ := crypto.ReadPassword()

			entry := store.Entry{
				Label:     label,
				User:      user,
				Host:      host,
				Port:      port,
				Password:  string(password),
				CreatedAt: time.Now(),
			}
			if err := store.AddEntry(pwd, entry); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			fmt.Println("Label added.")
		},
	}
}
