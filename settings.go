package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type settings struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func ReadConsumerKey() string {
	return readSettings().ConsumerKey
}

func ReadConsumerSecret() string {
	return readSettings().ConsumerSecret
}

func ReadAccessToken() string {
	return readSettings().AccessToken
}

func ReadAccessTokenSecret() string {
	return readSettings().AccessTokenSecret
}

func readSettings() settings {
	file, err := ioutil.ReadFile("./settings.json")
	if err != nil {
		fmt.Errorf("File error: %v\n", err)
	}
	var s settings
	err = json.Unmarshal(file, &s)
	if err != nil {
		fmt.Errorf("File failed to decode: %v\n", err)
	}
	return s
}
