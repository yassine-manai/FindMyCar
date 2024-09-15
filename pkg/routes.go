package pkg

import (
	"fmc/config"
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/uptrace/bun"
)

func SetupRouter(db *bun.DB) *gin.Engine {
	r := gin.Default()
	r.Use(config.CustomErrorHandler())
	
	carParkAPI := NewCarparkAPI(db)
	presentCarAPI := NewPresentCarAPI(db)
	zoneAPI := NewZoneAPI(db)
	zoneImageAPI := NewZoneImageAPI(db)
	CameraAPI := NewCameraAPI(db)
	CarDetailAPI := NewCarDetailAPI(db)

	fmt.Println("----------------------------  -+-+-+-+-  -----------------------------------")

	// Carpark routes
	r.GET("fyc/carparks", carParkAPI.GetAllCarparks)
	r.GET("fyc/carparks/:id", carParkAPI.GetCarparkByID)
	r.POST("fyc/carparks", carParkAPI.AddCarpark)
	r.PUT("fyc/carparks/:id", carParkAPI.UpdateCarpark)
	r.DELETE("fyc/carparks/:id", carParkAPI.DeleteCarpark)

	// PresentCar routes
	r.GET("fyc/presentcars", presentCarAPI.GetPresentCars)
	r.GET("fyc/presentcars/:lpn", presentCarAPI.GetPresentCarByLPN)
	r.POST("fyc/presentcars", presentCarAPI.CreatePresentCar)
	r.PUT("fyc/presentcars/:id", presentCarAPI.UpdatePresentCarById)
	r.PUT("fyc/presentcars", presentCarAPI.UpdatePresentCarBylpn)
	r.DELETE("fyc/presentcars/:id", presentCarAPI.DeletePresentCar)

	// Zone routes
	r.GET("/fyc/zones", zoneAPI.GetZones)
	r.POST("/fyc/zones", zoneAPI.CreateZone)
	r.PUT("/fyc/zones/:id", zoneAPI.UpdateZoneId)
	r.DELETE("/fyc/zones/:id", zoneAPI.DeleteZone)

	// ZoneImage routes
	r.GET("/fyc/zonesImage", zoneImageAPI.GetImageZones)
	r.POST("/fyc/zonesImage", zoneImageAPI.CreateZoneImage)
	r.PUT("/fyc/zonesImage/:id", zoneImageAPI.UpdateZoneImageById)
	r.DELETE("/fyc/zonesImage/:id", zoneImageAPI.DeleteZoneImage)

	// Camera routes
	r.GET("/fyc/cameras", CameraAPI.GetCamera)
	r.POST("/fyc/cameras", CameraAPI.CreateCamera)
	r.PUT("/fyc/cameras/:id", CameraAPI.UpdateCamera)
	r.DELETE("/fyc/cameras/:id", CameraAPI.DeleteCamera)

	// Car detail routes
	r.GET("/fyc/carDetails", CarDetailAPI.GetCarDetails)
	r.POST("/fyc/carDetails", CarDetailAPI.CreateCarDetail)
	r.PUT("/fyc/carDetails/:id", CarDetailAPI.UpdateCarDetailById)
	r.DELETE("/fyc/carDetails/:id", CarDetailAPI.DeleteCarDetail)

	r.GET("/fyc/v1", presentCarAPI.FYCHandler)
	r.POST("/fyc/v1/Auth/token", TokenHandler)

	// Swagger documentation endpoint
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
