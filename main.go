package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

var apiEndpoint = "https://voip.ms/api/v1/rest.php"

func readCredentials() credentials {
	f, err := os.Open("credentials.json")
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(f)
	var c credentials
	for {
		if err = dec.Decode(&c); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
	// fmt.Printf("%+v\n", c)
	return c
}

func printBalance() {
	c := readCredentials()
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	v := url.Values{}
	v.Add("api_username", c.Email)
	v.Add("api_password", c.Password)
	v.Add("method", "getBalance")
	// v.Add("advanced", "true")
	u.RawQuery = v.Encode()
	// fmt.Println(u.String())

	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body := resp.Body
	dec := json.NewDecoder(body)
	b := BalanceResp{}
	if err := dec.Decode(&b); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", b.Balance.CurrentBalance)
}

func blockNumber(phone *string, note *string) {
	c := readCredentials()

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	bodyWriter.WriteField("api_username", c.Email)
	bodyWriter.WriteField("api_password", c.Password)
	bodyWriter.WriteField("method", "setCallerIDFiltering")
	bodyWriter.WriteField("callerid", *phone)
	bodyWriter.WriteField("did", "all")
	bodyWriter.WriteField("routing", "sys:hangup")
	var filterNote = time.Now().Format("Jan.02/06")
	if *note != "" {
		filterNote += " - " + *note
	}
	bodyWriter.WriteField("note", filterNote)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	/*
	 *fmt.Println(bodyBuf)
	 *return
	 */

	req, err := http.NewRequest("POST", apiEndpoint, bodyBuf)
	req.Header.Set("Content-Type", contentType)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body := resp.Body
	dec := json.NewDecoder(body)
	s := SetCallerIDFilterResp{}
	if err := dec.Decode(&s); err != nil {
		log.Fatal(err)
	}
	switch s.Status {
	case "used_filter":
		fmt.Println("Voip.ms reports a filter for this number already exists:", *phone)
		break
	case "success":
		fmt.Printf("Yay, %s blocked successfully with filter id %s\n", *phone, s.Filtering.String())
		break
	default:
		fmt.Println("Unknown result:")
		fmt.Printf("%+v\n", s)
	}
}

func main() {
	var bNum = flag.String("block-number", "", "Number to block")
	var bNote = flag.String("block-note", "", "Add a note to filter")
	var pb = flag.Bool("print-balance", false, "Show the account balance")
	flag.Parse()

	if *pb {
		printBalance()
		return
	} else if *bNum != "" {
		blockNumber(bNum, bNote)
		return
	}
	flag.Usage()
}
