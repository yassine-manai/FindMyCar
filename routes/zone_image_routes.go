package routes

import (
	"fmc/pkg"

	"github.com/gin-gonic/gin"
)

func ZoneImageRoutes(r *gin.Engine) {
	r.GET("/fyc/zonesImages", pkg.GetAllImageZonesAPI)
	r.GET("/fyc/zonesImage", pkg.GetZoneImageByIDAPI)
	r.POST("/fyc/zonesImage", pkg.CreateZoneImageAPI)
	r.PUT("/fyc/zonesImage/:id", pkg.UpdateZoneImageByIdAPI)
	r.DELETE("/fyc/zonesImage/:id", pkg.DeleteZoneImageAPI)
}
