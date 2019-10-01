package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

type BungoRequester struct {
	APIKey string
}

// Get a response from Bungo
func (r BungoRequester) Get(url string) []byte {
	client := &http.Client{}
	req, clientErr := http.NewRequest("GET", url, nil)
	if clientErr != nil {
		log.Fatal(clientErr)
	}
	req.Header.Add("X-API-Key", r.APIKey)
	resp, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	return body
}
