package main

import (
	env "mailer-go/internal/environment"
	"mailer-go/internal/router"
)

func main() {
	env.NewEnv()
	router.Init()
}
