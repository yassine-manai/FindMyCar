package routes

import (
	"fmc/pkg"

	"github.com/gin-gonic/gin"
)

func HistoryRoutes(r *gin.Engine) {
	r.GET("/fyc/history", pkg.GetHistoryAPI)
	r.GET("/fyc/history/:lpn", pkg.GetHistoryByLPNAPI)
	r.POST("/fyc/history", pkg.CreateHistoryAPI)
	r.PUT("/fyc/history/:id", pkg.UpdateHistoryAPI)
	r.DELETE("/fyc/history/:id", pkg.DeleteHistoryAPI)
}
