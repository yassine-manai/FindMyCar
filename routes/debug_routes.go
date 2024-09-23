package routes

import (
	"fmc/debug"

	"github.com/gin-gonic/gin"
)

func DebugRoutes(r *gin.Engine) {

	r.GET("/fyc/debug", debug.Debuger_api)
}
