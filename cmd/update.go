package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	PorkbunAPIEndpoint string = "https://porkbun.com/api/json/v3/dns/editByNameType"
	IpifyAPIEndpoint   string = "https://api.ipify.org?format=json"
)

type Config struct {
	Key       string `json:"key"`
	SecretKey string `json:"secretKey"`
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
	Type      string `json:"type"`
}

type RequestPorkbun struct {
	SecretKey string `json:"secretapikey"`
	Key       string `json:"apikey"`
	Ip        string `json:"content"`
}

type ResponseIpify struct {
	Ip string `json:"ip"`
}

type ResponsePorkbun struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Gets public IP and update for the domain configured in \"$HOME/.config/bacon/config.json\"",
	RunE:  runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func getCurrentIP() (string, error) {
	var ip ResponseIpify

	res, err := http.Get(IpifyAPIEndpoint)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	json.Unmarshal(body, &ip)

	return ip.Ip, nil
}

func updateDNS(c Config) error {
	var result ResponsePorkbun
	endpoint := fmt.Sprintf("%s/%s/%s/%s", PorkbunAPIEndpoint, c.Domain, c.Type, c.Subdomain)

	currentIP, err := getCurrentIP()
	if err != nil {
		return err
	}

	payload := RequestPorkbun{
		SecretKey: c.SecretKey,
		Key:       c.Key,
		Ip:        currentIP,
	}

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req := bytes.NewBuffer(encodedPayload)

	res, err := http.Post(endpoint, "application/json", req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if result.Status == "ERROR" || result.Message == "Edit error: We were unable to edit the DNS record." {
		return fmt.Errorf("error updating DNS: %s", result.Message)
	}

	return nil
}

func runUpdate(cmd *cobra.Command, args []string) error {
	var config Config
	homeDir := os.Getenv("HOME")
	configPath := fmt.Sprintf("%s/%s", homeDir, ".config/bacon")

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}

	err = updateDNS(config)
	if err != nil {
		return err
	}

	fmt.Println("DNS updated successfully!")

	return nil
}
