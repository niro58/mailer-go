package main

import (
	"fmt"
	route "mailer-go/internal/routes"
	util "mailer-go/internal/utils"
	"path"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(path.Join(util.Root, "/.env"))
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	route.SetupRouter()
}
