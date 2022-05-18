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

	var subdomains []string

	crtsh_subs := crtsh(domain, &subdomains)
	log.Println("News subs from https://crt.sh:", crtsh_subs)

	bufferover_subs := bufferover(domain, &subdomains)
	log.Println("News subs from https://dns.bufferover.run:", bufferover_subs)

	log.Println(len(subdomains), "subs foundeed ", subdomains)
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

func crtsh(domain string, subdomains_list *[]string) []string {
	var new_subs []string

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

	for _, val := range obj_crt {
		for _, s := range strings.Split(val.NameValue, "\n") {
			if !(contains((*subdomains_list), s)) {
				(*subdomains_list) = append((*subdomains_list), s)
				new_subs = append(new_subs, s)
			}
		}
	}

	return new_subs
}

type ObjMeta struct {
	Runtime   string   `json:"Runtime"`
	Errors    *string  `json:"Errors"`
	Message   string   `json:"Message"`
	FileNames []string `json:"FileNames"`
	TOS       string   `json:"TOS"`
}

type ObjBufferover struct {
	Meta ObjMeta  `json:"Meta"`
	FDNS []string `json:"FDNS_A"`
	RDNS []string `json:"RDNS"`
}

func bufferover(domain string, subdomains_list *[]string) []string {
	url := "https://dns.bufferover.run/dns?q=." + domain + "&output=json"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	var obj_bufferover ObjBufferover

	if err = json.NewDecoder(resp.Body).Decode(&obj_bufferover); err != nil {
		log.Fatalln(err)
	}

	var new_subs []string

	for _, val := range obj_bufferover.FDNS {
		sub := strings.Split(val, ",")[1]
		if !(contains((*subdomains_list), sub)) {
			(*subdomains_list) = append((*subdomains_list), sub)
			new_subs = append(new_subs, sub)
		}
	}

	for _, val := range obj_bufferover.RDNS {
		sub := strings.Split(val, ",")[1]
		if !(contains((*subdomains_list), sub)) {
			(*subdomains_list) = append((*subdomains_list), sub)
			new_subs = append(new_subs, sub)
		}
	}

	return new_subs
}

func contains(array []string, str string) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}

	return false
}
