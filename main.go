package main

import (
	"fmt"
	route "mailer-go/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	route.SetupRouter()
}
