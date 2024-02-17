package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bacon",
	Short: "A DNS updater for Porkbun",
	Long:  "A cli to update the public IP of your domain registered in Porkbun.",
	Run:   nil,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
