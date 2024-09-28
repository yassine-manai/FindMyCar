package debug

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fmc/pkg"
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
		"CameraList":  pkg.CameraList,
	})

}
