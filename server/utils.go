package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func loadConfigFromFile(ntwConf *NetworkConfig) error {
	// Get config from json config file
	absPath, err := filepath.Abs("server\\config.json")
	if err != nil {
		return err
	}
	b, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	// Getting the network config variable
	err = json.Unmarshal(b, ntwConf)
	if err != nil {
		return err
	}
	return nil
}

func getServerConfig(id int, config *NetworkConfig) *ServerConfig {
	for _, srv := range config.Servers {
		if srv.Id == id {
			return &srv
		}
	}
	return nil
}

func cntChar(char uint8, text string) int {
	cnt := 0
	for i, _ := range text {
		if text[i] == char {
			cnt++
		}
	}
	return cnt
}
