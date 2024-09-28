package routes

import (
	"github.com/gin-gonic/gin"

	"fmc/pkg"
)

func ClientCredsRoutes(r *gin.Engine) {
	// Routes for client credentials operations
	r.GET("/fyc/clientCreds", pkg.GetAllClientCredsApi)
	r.GET("/fyc/clientEnabled", pkg.GetClientEnabledAPI)
	r.GET("/fyc/clientsDeleted", pkg.GetClientDeletedAPI)

	r.POST("/fyc/clientCreds", pkg.AddClientCredAPI)
	r.PUT("/fyc/clientCreds", pkg.UpdateClientCredAPI)
	r.PUT("/fyc/clientState", pkg.ChangeClientStateAPI)

	r.DELETE("/fyc/clientCreds", pkg.DeleteClientCredAPI)
}
