package routes

import (
	"fmc/pkg"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/token", pkg.GetToken)
	r.GET("/findmycar", pkg.FindMyCar)
	r.GET("/getpicture", pkg.GetPicture)
	r.GET("/getSettings", pkg.Getsettings)
	r.POST("/fyc/v1/Auth/token", pkg.TokenHandler)
}
