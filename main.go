package main

import (
	"fmt"
	"mailer-go/internal/config"
	route "mailer-go/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	config.Senders = config.GetSenders()
	route.SetupRouter()
}
