package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	file, err := ioutil.ReadFile(filepath.Join(dir, "settings.json"))
	if err != nil {
		fmt.Errorf("File error: %v\n", err)
		panic(err)
	}
	var s settings
	err = json.Unmarshal(file, &s)
	if err != nil {
		fmt.Errorf("File failed to decode: %v\n", err)
		panic(err)
	}
	return s
}
