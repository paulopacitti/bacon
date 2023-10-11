package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const PorkbunAPIEndpoint string = "https://porkbun.com/api/json/v3/dns/editByNameType"
const IpifyAPIEndpoint string = "https://api.ipify.org?format=json"

type DNSRecord struct {
	Domain string
	Subdomain string
	RecordType string
}

type RequestPorkbun struct {
	SecretKey string `json:"secretapikey"`
	Key string `json:"apikey"`
	Ip string `json:"content"`
}

type ResponseIpify struct {
	Ip string `json:"ip"`
}

type ResponsePorkbun struct {
	Status string `json:"status"`
	Message string `json:"message"`
}

func getCurrentIP() string {
	var ip ResponseIpify

	res, err := http.Get(IpifyAPIEndpoint)
	if err != nil {
        return err.Error()
    }
    defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
    if err != nil {
        return err.Error()
    }
    json.Unmarshal(body, &ip)

    return ip.Ip
}

func updateDNS(r DNSRecord, sk string, k string) (bool, string) {
	var result ResponsePorkbun;
	endpoint := fmt.Sprintf("%s/%s/%s/%s", PorkbunAPIEndpoint, r.Domain, r.RecordType,  r.Subdomain)

	currentIP := getCurrentIP()
	payload  := RequestPorkbun{
		SecretKey: sk,
		Key: k,
		Ip: currentIP,
	}

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
        return false, err.Error()
    }
	req := bytes.NewBuffer(encodedPayload)

	res, err := http.Post(endpoint, "application/json", req)
	if err != nil {
        return false, err.Error()
    }
    defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
    if err != nil {
        return false, err.Error()
    }
    json.Unmarshal(body, &result)

	return result.Status == "SUCCESS", result.Message
}


func main() {
	var r DNSRecord;
	var sk, k string;

	fmt.Println("Input your sk:")
	fmt.Scan(&sk)
	fmt.Println("Input your k:")
	fmt.Scan(&k)
	fmt.Println("Input your domain:")
	fmt.Scan(&r.Domain)
	fmt.Println("Input your subdomain:")
	fmt.Scan(&r.Subdomain)
	fmt.Println("Input your record type:")
	fmt.Scan(&r.RecordType)
	
	ok, err := updateDNS(r, sk, k)
	fmt.Println(ok, err)
}