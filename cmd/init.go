package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tryoasnafi/gate/internal/crypto"
	"github.com/tryoasnafi/gate/internal/session"
	"github.com/tryoasnafi/gate/internal/store.go"
)

func InitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize gate with a master password",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Enter master password: ")
			pwd, err := crypto.ReadPassword()
			if err != nil {
				fmt.Println("Error reading password:", err)
				os.Exit(1)
			}

			if err := store.InitStore(pwd); err != nil {
				fmt.Println("Initialization failed:", err)
				os.Exit(1)
			}

			session.Create(pwd)
			fmt.Println("gate initialized.")
		},
	}
}
