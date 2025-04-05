package router

import (
	"fmt"
	env "mailer-go/internal/environment"
	handler "mailer-go/internal/handlers"
	"mailer-go/internal/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() *gin.Engine {
	app := handler.CreateApplication()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	//Gin Logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	output := zerolog.ConsoleWriter{Out: os.Stdout}

	log.Logger = zerolog.New(output).
		With().
		Timestamp().
		Str("service", "mailer-go").
		Logger()

	r.Use(middleware.Logger())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := r.Group("/")
	api.Use(middleware.AuthRequired())

	api.GET("/health", app.Health)
	api.POST("/send", app.Send)
	api.POST("/send-template", app.SendTemplate)
	api.GET("/status", app.Status)

	if env.Env.GinMode == "release" {
		r.Run(fmt.Sprintf(":%s", env.Env.Port))
	} else {
		r.Run(fmt.Sprintf("localhost:%s", env.Env.Port))
	}

	return r
}
