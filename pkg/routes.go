package pkg

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

	fmt.Println("----------------------------  -+-+-+-+-  -----------------------------------")

	// Carpark routes
	r.GET("/fyc/carparks", GetAllCarparksAPI)
	r.GET("/fyc/carparks/:id", GetCarparkByIDAPI)
	r.POST("/fyc/carparks", AddCarparkAPI)
	r.PUT("/fyc/carparks/:id", UpdateCarparkAPI)
	r.DELETE("/fyc/carparks/:id", DeleteCarparkAPI)

	// PresentCar routes
	r.GET("/fyc/presentcars", GetPresentCarsAPI)
	r.GET("/fyc/presentcars/:lpn", GetPresentCarByLPNAPI)
	r.POST("/fyc/presentcars", CreatePresentCarAPI)
	r.PUT("/fyc/presentcars/:id", UpdatePresentCarByIdAPI)
	r.PUT("/fyc/presentcars", UpdatePresentCarBylpnAPI)
	r.DELETE("/fyc/presentcars/:id", DeletePresentCarAPI)

	// Zone routes
	r.GET("/fyc/zones", GetZonesAPI)
	r.GET("/fyc/zones/:id", GetZoneByIDAPI)
	r.POST("/fyc/zones", CreateZoneAPI)
	r.PUT("/fyc/zones/:id", UpdateZoneIdAPI)
	r.DELETE("/fyc/zones/:id", DeleteZoneAPI)

	// ZoneImage routes
	r.GET("/fyc/zonesImages", GetAllImageZonesAPI)
	r.GET("/fyc/zonesImage", GetZoneImageByIDAPI)
	r.POST("/fyc/zonesImage", CreateZoneImageAPI)
	r.PUT("/fyc/zonesImage/:id", UpdateZoneImageByIdAPI)
	r.DELETE("/fyc/zonesImage/:id", DeleteZoneImageAPI)

	// Camera routes
	r.GET("/fyc/cameras", GetCameraAPI)
	r.GET("/fyc/cameras/:id", GetCameraByIDAPI)
	r.POST("/fyc/cameras", CreateCameraAPI)
	r.PUT("/fyc/cameras/:id", UpdateCameraAPI)
	r.DELETE("/fyc/cameras/:id", DeleteCameraAPI)

	// Car detail routes
	r.GET("/fyc/carDetails", GetCarDetailsAPI)
	r.GET("/fyc/carDetails/:id", GetCarDetailsByIdAPI)
	r.POST("/fyc/carDetails", CreateCarDetailAPI)
	r.PUT("/fyc/carDetails/:id", UpdateCarDetailByIdAPI)
	r.DELETE("/fyc/carDetails/:id", DeleteCarDetailAPI)

	// Client Creds routes
	r.GET("/fyc/clientCreds", GetAllClientCredsApi)
	r.GET("/fyc/clientCreds/:id", GetClientCredByIDAPI)
	r.POST("/fyc/clientCreds", AddClientCredAPI)
	r.PUT("/fyc/clientCreds/:id", UpdateClientCredAPI)
	r.DELETE("/fyc/clientCreds/:id", DeleteClientCredAPI)

	// PresentCarHistory routes
	r.GET("/fyc/history", GetHistoryAPI)
	r.GET("/fyc/history/:lpn", GetHistoryByLPNAPI)
	r.POST("/fyc/history", CreateHistoryAPI)
	r.PUT("/fyc/history/:id", UpdateHistoryAPI)
	r.DELETE("/fyc/history/:id", DeleteHistoryAPI)

	r.POST("/token", getToken)
	r.GET("/findmycar", findMyCar)
	r.GET("/getpicture", getPicture)
	r.GET("/getSettings", getsettings)

	r.POST("/fyc/v1/Auth/token", TokenHandler)
	r.GET("/fyc/debug", Debuger_api)

	// Swagger  endpoint
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
