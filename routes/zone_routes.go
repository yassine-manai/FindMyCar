package routes

import (
	"fmc/pkg"

	"github.com/gin-gonic/gin"
)

func ZoneRoutes(r *gin.Engine) {
	r.GET("/fyc/zones", pkg.GetZonesAPI)
	r.GET("/fyc/zones/:id", pkg.GetZoneByIDAPI)
	r.POST("/fyc/zones", pkg.CreateZoneAPI)
	r.PUT("/fyc/zones/:id", pkg.UpdateZoneIdAPI)
	r.DELETE("/fyc/zones/:id", pkg.DeleteZoneAPI)
}
