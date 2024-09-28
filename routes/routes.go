package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"fmc/config"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(config.CustomErrorHandler())

	fmt.Println("--------------------------  START ROUTING  ----------------------")
	DebugRoutes(r) // Debug routes
	AuthRoutes(r)  // TestRoutes

	ZoneImageRoutes(r)
	CameraRoutes(r)
	ZoneRoutes(r)
	PresentCarRoutes(r)
	CarDetailRoutes(r)
	ClientCredsRoutes(r)
	SignRoutes(r)
	UserAuditRoutes(r)
	UserRoutes(r)
	SettingsRoutes(r)

	HistoryRoutes(r)
	ErrorRoutes(r)

	fmt.Println("--------------------------  END ROUTING  ----------------------")

	// Swagger endpoint
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
