package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Creates config file with the domain that needs to be updated.",

	Run: runConfig,
}

func init() {
	configCmd.Flags().StringP("key", "k", "", "Porkbun API key for the account.")
	configCmd.Flags().StringP("secretKey", "p", "", "Porkbun API secret key for the account.")
	configCmd.Flags().StringP("domain", "d", "", "Domain that needs to be updated (example.com).")
	configCmd.Flags().StringP("subdomain", "b", "", "Subdomain that needs to be updated (www, *...).")
	configCmd.Flags().StringP("type", "t", "", "Type of DNS record that needs to be updated (A, CNAME...).")

	configCmd.MarkFlagRequired("key")
	configCmd.MarkFlagRequired("secretKey")
	configCmd.MarkFlagRequired("domain")
	configCmd.MarkFlagRequired("subdomain")
	configCmd.MarkFlagRequired("type")

	viper.BindPFlag("key", configCmd.Flags().Lookup("key"))
	viper.BindPFlag("secretKey", configCmd.Flags().Lookup("secretKey"))
	viper.BindPFlag("domain", configCmd.Flags().Lookup("domain"))
	viper.BindPFlag("subdomain", configCmd.Flags().Lookup("subdomain"))
	viper.BindPFlag("type", configCmd.Flags().Lookup("type"))

	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) {
	homeDir := os.Getenv("HOME")
	configPath := fmt.Sprintf("%s/%s", homeDir, ".config/bacon")
	filePath := fmt.Sprintf("%s/%s", configPath, "config.json")
	err := os.MkdirAll(configPath, 0700)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	viper.SetConfigFile(filePath)
	viper.SetConfigType("json")

	viper.WriteConfigAs(filePath)
}
