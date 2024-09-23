package routes

import (
	"fmc/pkg"

	"github.com/gin-gonic/gin"
)

func ErrorRoutes(r *gin.Engine) {
	r.GET("/fyc/errors", pkg.GetAllErrorCode)
	r.POST("/fyc/errors", pkg.CreateErrorMessageAPI)
	r.PUT("/fyc/errors", pkg.UpdateErrorMessageAPI)
	r.DELETE("/fyc/errors", pkg.DeleteErrorMessageAPI)
}
