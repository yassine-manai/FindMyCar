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

	fmt.Println("----------------------------  -+-+-+-+-  -----------------------------------")

	// Carpark routes
	r.GET("/fyc/carparks", func(c *gin.Context) { GetAllCarparksAPI(c, db) })
	r.GET("/fyc/carparks/:id", func(c *gin.Context) { GetCarparkByIDAPI(c, db) })
	r.POST("/fyc/carparks", func(c *gin.Context) { AddCarparkAPI(c, db) })
	r.PUT("/fyc/carparks/:id", func(c *gin.Context) { UpdateCarparkAPI(c, db) })
	r.DELETE("/fyc/carparks/:id", func(c *gin.Context) { DeleteCarparkAPI(c, db) })

	// PresentCar routes
	r.GET("/fyc/presentcars", func(c *gin.Context) { GetPresentCarsAPI(c, db) })
	r.GET("/fyc/presentcars/:lpn", func(c *gin.Context) { GetPresentCarByLPNAPI(c, db) })
	r.POST("/fyc/presentcars", func(c *gin.Context) { CreatePresentCarAPI(c, db) })
	r.PUT("/fyc/presentcars/:id", func(c *gin.Context) { UpdatePresentCarByIdAPI(c, db) })
	r.PUT("/fyc/presentcars", func(c *gin.Context) { UpdatePresentCarBylpnAPI(c, db) })
	r.DELETE("/fyc/presentcars/:id", func(c *gin.Context) { DeletePresentCarAPI(c, db) })

	// Zone routes
	r.GET("/fyc/zones", func(c *gin.Context) { GetZonesAPI(c, db) })
	r.GET("/fyc/zones/:id", func(c *gin.Context) { GetZoneByIDAPI(c, db) })
	r.POST("/fyc/zones", func(c *gin.Context) { CreateZoneAPI(c, db) })
	r.PUT("/fyc/zones/:id", func(c *gin.Context) { UpdateZoneIdAPI(c, db) })
	r.DELETE("/fyc/zones/:id", func(c *gin.Context) { DeleteZoneAPI(c, db) })

	// ZoneImage routes
	r.GET("/fyc/zonesImage", func(c *gin.Context) { GetImageZonesAPI(c, db) })
	r.GET("/fyc/zonesImage/:id", func(c *gin.Context) { GetZoneImageByIDAPI(c, db) })
	r.POST("/fyc/zonesImage", func(c *gin.Context) { CreateZoneImageAPI(c, db) })
	r.PUT("/fyc/zonesImage/:id", func(c *gin.Context) { UpdateZoneImageByIdAPI(c, db) })
	r.DELETE("/fyc/zonesImage/:id", func(c *gin.Context) { DeleteZoneImageAPI(c, db) })

	// Camera routes
	r.GET("/fyc/cameras", func(c *gin.Context) { GetCameraAPI(c, db) })
	r.GET("/fyc/cameras/:id", func(c *gin.Context) { GetCameraByIDAPI(c, db) })
	r.POST("/fyc/cameras", func(c *gin.Context) { CreateCameraAPI(c, db) })
	r.PUT("/fyc/cameras/:id", func(c *gin.Context) { UpdateCameraAPI(c, db) })
	r.DELETE("/fyc/cameras/:id", func(c *gin.Context) { DeleteCameraAPI(c, db) })

	// Car detail routes
	r.GET("/fyc/carDetails", func(c *gin.Context) { GetCarDetailsAPI(c, db) })
	r.GET("/fyc/carDetails/:id", func(c *gin.Context) { GetCarDetailsByIdAPI(c, db) })
	r.POST("/fyc/carDetails", func(c *gin.Context) { CreateCarDetailAPI(c, db) })
	r.PUT("/fyc/carDetails/:id", func(c *gin.Context) { UpdateCarDetailByIdAPI(c, db) })
	r.DELETE("/fyc/carDetails/:id", func(c *gin.Context) { DeleteCarDetailAPI(c, db) })

	// Client Creds routes
	r.GET("/fyc/clientCreds", func(c *gin.Context) { GetAllClientCredsApi(c, db) })
	r.GET("/fyc/clientCreds/:id", func(c *gin.Context) { GetClientCredByIDAPI(c, db) })
	r.POST("/fyc/clientCreds", func(c *gin.Context) { AddClientCredAPI(c, db) })
	r.PUT("/fyc/clientCreds/:id", func(c *gin.Context) { UpdateClientCredAPI(c, db) })
	r.DELETE("/fyc/clientCreds/:id", func(c *gin.Context) { DeleteClientCredAPI(c, db) })

	r.POST("/token", getToken)
	r.GET("/findmycar", findMyCar)
	r.GET("/getpicture", getPicture)

	r.POST("/fyc/v1/Auth/token", TokenHandler)

	// Swagger documentation endpoint
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
