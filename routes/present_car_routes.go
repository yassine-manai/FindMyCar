package routes

import (
	"fmc/pkg"

	"github.com/gin-gonic/gin"
)

func PresentCarRoutes(r *gin.Engine) {
	r.GET("/fyc/presentcars", pkg.GetPresentCarsAPI)
	r.GET("/fyc/presentcars/:lpn", pkg.GetPresentCarByLPNAPI)
	r.POST("/fyc/presentcars", pkg.CreatePresentCarAPI)
	r.PUT("/fyc/presentcars/:id", pkg.UpdatePresentCarByIdAPI)
	r.PUT("/fyc/presentcars", pkg.UpdatePresentCarBylpnAPI)
	r.DELETE("/fyc/presentcars/:id", pkg.DeletePresentCarAPI)
}

