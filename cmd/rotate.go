package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tryoasnafi/gate/internal/crypto"
	"github.com/tryoasnafi/gate/internal/store.go"
)

func RotateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rotate",
		Short: "Rotate the master password",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Enter old master password: ")
			oldPwd, err := crypto.ReadPassword()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			fmt.Print("Enter new master password: ")
			newPwd, err := crypto.ReadPassword()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}

			if err := store.RotateMasterPassword(oldPwd, newPwd); err != nil {
				fmt.Println("Rotation failed:", err)
				os.Exit(1)
			}

			fmt.Println("Master password rotated.")
		},
	}
}
