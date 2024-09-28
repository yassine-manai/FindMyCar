package routes

import (
	"github.com/gin-gonic/gin"

	"fmc/pkg"
)

func CameraRoutes(r *gin.Engine) {
	r.GET("/fyc/cameras", pkg.GetCameraAPI)
	r.GET("/fyc/camerasEnabled", pkg.GetCameraEnabledAPI)
	r.GET("/fyc/camerasDeleted", pkg.GetCameraDeletedAPI)
	r.POST("/fyc/cameras", pkg.CreateCameraAPI)
	r.PUT("/fyc/cameras", pkg.UpdateCameraAPI)
	r.PUT("/fyc/cameraState", pkg.ChangeCameraStateAPI)
	r.DELETE("/fyc/cameras", pkg.DeleteCameraAPI)
}
