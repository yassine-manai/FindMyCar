package routes

import (
	"fmc/pkg"

	"github.com/gin-gonic/gin"
)

func carParkRouter(r *gin.Engine) {

	// Carpark routes
	r.GET("/fyc/carparks", pkg.GetAllCarparksAPI)
	r.GET("/fyc/carparks/:id", pkg.GetCarparkByIDAPI)
	r.POST("/fyc/carparks", pkg.AddCarparkAPI)
	r.PUT("/fyc/carparks/:id", pkg.UpdateCarparkAPI)
	r.DELETE("/fyc/carparks/:id", pkg.DeleteCarparkAPI)

}
