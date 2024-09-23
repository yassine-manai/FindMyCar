package routes

import (
	"fmc/config"
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(config.CustomErrorHandler())

	fmt.Println("--------------------------  START ROUTING  ----------------------")

	carParkRouter(r)
	PresentCarRoutes(r)
	ZoneRoutes(r)
	ZoneImageRoutes(r)
	CameraRoutes(r)
	CarDetailRoutes(r)
	ClientCredsRoutes(r)
	HistoryRoutes(r)
	AuthRoutes(r)
	DebugRoutes(r)

	fmt.Println("--------------------------  END ROUTING  ----------------------")

	// Swagger endpoint
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
