package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func get(url, token string) (body []byte, err error) {
	log.Println("[get] Entered")
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
