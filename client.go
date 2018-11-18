package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

var apiEndpoint = "https://voip.ms/api/v1/rest.php"

type client struct {
	url string
	credentials
}

func newClient() *client {
	c := readCredentials()
	return &client{apiEndpoint, c}
}

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
	return c
}

func (c *client) doRequest(req *http.Request, respStruct interface{}) {

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body := resp.Body
	dec := json.NewDecoder(body)
	if err := dec.Decode(respStruct); err != nil {
		log.Fatal(err)
	}
	return

}

func (c *client) getRequest(method string, values url.Values, respStruct interface{}) {
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	values.Add("api_username", c.Email)
	values.Add("api_password", c.Password)
	values.Add("method", method)
	u.RawQuery = values.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	c.doRequest(req, respStruct)
	return
}

func (c *client) postRequest(method string, p postParams, respStruct interface{}) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	bodyWriter.WriteField("api_username", c.Email)
	bodyWriter.WriteField("api_password", c.Password)
	bodyWriter.WriteField("method", method)

	for k, v := range p {
		bodyWriter.WriteField(k, v)
	}

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

	c.doRequest(req, respStruct)
	return
}
