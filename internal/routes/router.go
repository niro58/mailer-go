package route

import (
	handler "mailer-go/internal/handlers"
	"mailer-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	app := handler.CreateApplication()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	authRequired := r.Group("/")
	authRequired.Use(middleware.AuthRequired())
	{
		r.GET("/health", app.Health)
		r.POST("/send", app.Send)
	}

	r.Run(":8085")
	return r
}
