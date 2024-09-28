package routes

import (
	"github.com/gin-gonic/gin"

	"fmc/pkg"
)

func SettingsRoutes(r *gin.Engine) {
	r.GET("/fyc/settings", pkg.GetSettingsAPI)
	r.POST("/fyc/settings", pkg.AddSettingsAPI)
	r.PUT("/fyc/settings", pkg.UpdateSettingsAPI)
}
