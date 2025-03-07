package route

import (
	handler "mailer-go/internal/handlers"
	"mailer-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter() *gin.Engine {
	app := handler.CreateApplication()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := r.Group("/")
	api.Use(middleware.AuthRequired())

	api.GET("/health", app.Health)
	api.POST("/send", app.Send)

	r.Run(":8085")
	return r
}
