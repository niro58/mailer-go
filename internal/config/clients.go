package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type SenderAuth struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var Senders map[string]SenderAuth

func GetSenders() map[string]SenderAuth {
	clientsPath := os.Getenv("CLIENTS_PATH")

	if clientsPath == "" {
		panic("CLIENTS_PATH not set")
	}
	jsonFile, err := os.Open(clientsPath)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var result map[string]SenderAuth
	json.Unmarshal([]byte(byteValue), &result)

	return result
}

var ErrSenderNotFound = errors.New("sender not found")
