package main

import (
	"fmt"
	"mailer-go/internal/environment"
	route "mailer-go/internal/routes"
)

func main() {
	environment.Environment = environment.NewEnv()
	fmt.Println(environment.Environment.ClientsPath)
	route.SetupRouter()
}
