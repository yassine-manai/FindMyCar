package routes

import (
	"fmc/pkg"

	"github.com/gin-gonic/gin"
)

func CarDetailRoutes(r *gin.Engine) {
	r.GET("/fyc/carDetails", pkg.GetCarDetailsAPI)
	r.POST("/fyc/carDetails", pkg.CreateCarDetailAPI)
	r.PUT("/fyc/carDetails", pkg.UpdateCarDetailByIdAPI)
	r.DELETE("/fyc/carDetails", pkg.DeleteCarDetailAPI)
}
