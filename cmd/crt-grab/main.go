package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	domain := "hackerone.com"

	crtsh(domain)
}

type ObjCrt struct {
	Id             int64  `json:"id"`
	IssuerCaId     int64  `json:"issuer_ca_id"`
	IssuerName     string `json:"issuer_name"`
	CommonName     string `json:"common_name"`
	NameValue      string `json:"name_value"`
	EntryTimestamp string `json:"entry_timestamp"`
	NotBefore      string `json:"not_before"`
	NotAfter       string `json:"not_after"`
	SerialNumber   string `json:"serial_number"`
}

func crtsh(domain string) []string {
	url := "https://crt.sh/?q=." + domain + "&output=json"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	json_str, _ := ioutil.ReadAll(resp.Body)

	var obj_crt []ObjCrt
	if err = json.Unmarshal([]byte(json_str), &obj_crt); err != nil {
		log.Fatalln(err)
	}

	var subdomains []string

	for _, val := range obj_crt {
		for _, s := range strings.Split(val.NameValue, "\n") {
			if !(contains(subdomains, s)) {
				subdomains = append(subdomains, s)
			}
		}
	}
	return subdomains
}

func contains(array []string, str string) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}

	return false
}
