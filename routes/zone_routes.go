package routes

import (
	"github.com/gin-gonic/gin"

	"fmc/pkg"
)

func ZoneRoutes(r *gin.Engine) {
	r.GET("/fyc/zones", pkg.GetZonesAPI)
	r.GET("/fyc/zonesEnabled", pkg.GetZoneEnabledAPI)
	r.GET("/fyc/zonesDeleted", pkg.GetZoneDeletedAPI)
	r.POST("/fyc/zones", pkg.CreateZoneAPI)
	r.PUT("/fyc/zones/:id", pkg.UpdateZoneIdAPI)
	r.PUT("/fyc/zoneState", pkg.ChangeZoneStateAPI)
	r.DELETE("/fyc/zones/:id", pkg.DeleteZoneAPI)
}
