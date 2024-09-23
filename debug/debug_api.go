package debug

import (
	"fmc/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DebugAPI godoc
//
//	@Summary		Debug API
//	@Tags			Debug
//	@Produce		json
//	@Router			/fyc/debug [get]
func Debuger_api(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"ZoneList":    pkg.Zonelist,
		"CarParkList": pkg.CarParkList,
		"CameraList":  pkg.CameraList,
	})

}
