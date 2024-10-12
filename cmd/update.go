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
	PorkbunAPIDNSUpdateEndpoint   string = "https://api.porkbun.com/api/json/v3/dns/editByNameType"
	PorkbunAPIDNSRetrieveEndpoint string = "https://api.porkbun.com/api/json/v3/dns/retrieveByNameType"
	IpifyAPIEndpoint              string = "https://api.ipify.org?format=json"
)

type Config struct {
	Key       string `json:"key"`
	SecretKey string `json:"secretKey"`
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
	Type      string `json:"type"`
}

// Request body for the DNS Edit Porkbun API endpoint,
// as defined in the documentation: https://porkbun.com/api/json/v3/documentation
type RequestPorkbunDNSUpdate struct {
	SecretKey string `json:"secretapikey"`
	Key       string `json:"apikey"`
	Ip        string `json:"content"`
}

// Response body for the DNS Edit Porkbun API endpoint,
// as defined in the documentation: https://porkbun.com/api/json/v3/documentation
type ResponsePorkbunDNSUpdate struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Request body for the DNS Retrieve Records Porkbun API endpoint,
// as defined in the documentation: https://porkbun.com/api/json/v3/documentation
type RequestPorkbunDNSURetrieve struct {
	SecretKey string `json:"secretapikey"`
	Key       string `json:"apikey"`
}

// Response body for the DNS Retrieve Records Porkbun API endpoint,
// as defined in the documentation: https://porkbun.com/api/json/v3/documentation
type ResponsePorkbunDNSURetrieve struct {
	Status  string `json:"status"`
	Records []struct {
		Content  string `json:"content"`
		Id       string `json:"id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		Ttl      string `json:"ttl"`
		Priority string `json:"prio"`
		Notes    string `json:"notes"`
	} `json:"records"`
}

// Response body for the Ipify API endpoint.
type ResponseIpify struct {
	Ip string `json:"ip"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Gets public IP and update for the domain configured in \"$HOME/.config/bacon/config.json\"",
	RunE:  runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

// Retrieves the current IP address using the Ipify API.
// If an error occurs during the HTTP request or response parsing, it is returned as the second value.
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

// Retrieves the current DNS record from PorkbunAPI for the domain configured in the config file.
// If an error occurs during the HTTP request or response parsing, it is returned as the second value.
func retrieveDNS(c Config) (string, error) {
	var result ResponsePorkbunDNSURetrieve
	endpoint := fmt.Sprintf("%s/%s/%s/%s", PorkbunAPIDNSRetrieveEndpoint, c.Domain, c.Type, c.Subdomain)

	payload := RequestPorkbunDNSURetrieve{
		SecretKey: c.SecretKey,
		Key:       c.Key,
	}

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req := bytes.NewBuffer(encodedPayload)
	res, err := http.Post(endpoint, "application/json", req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if result.Status == "ERROR" {
		return "", fmt.Errorf("error retrieving DNS: %s", result.Status)
	}

	if len(result.Records) == 0 {
		return "", fmt.Errorf("no record for subdomain %s.%s with type %s found", c.Subdomain, c.Domain, c.Type)
	}

	return result.Records[0].Content, nil
}

// Updates the DNS record from PorkbunAPI for the domain configured in the config file.
// If an error occurs during the HTTP request or response parsing, it is returned as the second value.
func updateDNS(c Config, currentIP string) error {
	var result ResponsePorkbunDNSUpdate
	endpoint := fmt.Sprintf("%s/%s/%s/%s", PorkbunAPIDNSUpdateEndpoint, c.Domain, c.Type, c.Subdomain)

	payload := RequestPorkbunDNSUpdate{
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

// Run the update command.
func runUpdate(cmd *cobra.Command, args []string) error {
	var config Config
	homeDir, _ := os.UserHomeDir()
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

	currentDNSIP, err := retrieveDNS(config)
	if err != nil {
		return err
	}

	currentIP, err := getCurrentIP()
	if err != nil {
		return err
	}

	if currentDNSIP == currentIP {
		fmt.Println("Current IP is already up to date!")
		return nil
	}

	err = updateDNS(config, currentIP)
	if err != nil {
		return err
	}

	fmt.Println("DNS updated successfully!")

	return nil
}
