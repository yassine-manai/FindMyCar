package routes

import (
	"fmc/pkg"

	"github.com/gin-gonic/gin"
)

func ClientCredsRoutes(r *gin.Engine) {
	r.GET("/fyc/clientCreds", pkg.GetAllClientCredsApi)
	r.GET("/fyc/clientCreds/:id", pkg.GetClientCredByIDAPI)
	r.POST("/fyc/clientCreds", pkg.AddClientCredAPI)
	r.PUT("/fyc/clientCreds/:id", pkg.UpdateClientCredAPI)
	r.DELETE("/fyc/clientCreds/:id", pkg.DeleteClientCredAPI)
}
