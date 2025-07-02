package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tryoasnafi/gate/cmd"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gate",
		Short: "gate - The gatekeeper to your SSH access",
	}

	rootCmd.AddCommand(
		cmd.InitCmd(),
		cmd.RotateCmd(),
		cmd.ListCmd(),
		cmd.NewCmd(),
		cmd.ConnectCmd(),
		cmd.DeleteCmd(),
		cmd.ImportCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
