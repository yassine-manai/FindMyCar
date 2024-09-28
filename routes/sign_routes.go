package routes

import (
	"github.com/gin-gonic/gin"

	"fmc/pkg"
)

func SignRoutes(r *gin.Engine) {
	r.GET("/fyc/sign", pkg.GetSignAPI)
	r.GET("/fyc/signEnabled", pkg.GetSignEnabledAPI)
	r.GET("/fyc/signDeleted", pkg.GetSignDeletedAPI)
	r.POST("/fyc/sign", pkg.CreateSignAPI)
	r.PUT("/fyc/sign", pkg.UpdateSignAPI)
	r.PUT("/fyc/signState", pkg.ChangeSigntateAPI)
	r.DELETE("/fyc/sign", pkg.DeleteSignAPI)
}
