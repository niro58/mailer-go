package route

import (
	handler "email-sender/handlers"
	"email-sender/middleware"

	"github.com/gin-gonic/gin"
)
func SetupRouter() *gin.Engine {
	app := handler.CreateApplication()

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.POST("/send", app.Send)
	r.Run(":8081")
	return r
}
