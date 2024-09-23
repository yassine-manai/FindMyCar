package routes

import (
	"fmc/pkg"

	"github.com/gin-gonic/gin"
)

func CameraRoutes(r *gin.Engine) {
	r.GET("/fyc/cameras", pkg.GetCameraAPI)
	r.POST("/fyc/cameras", pkg.CreateCameraAPI)
	r.PUT("/fyc/cameras", pkg.UpdateCameraAPI)
	r.DELETE("/fyc/cameras", pkg.DeleteCameraAPI)
}
