package routes

import (
	"github.com/gin-gonic/gin"

	"fmc/pkg"
)

func UserAuditRoutes(r *gin.Engine) {
	r.GET("/fyc/UserAudit", pkg.GetUserAuditAPI)
	r.POST("/fyc/UserAudit", pkg.CreateUserAuditAPI)
	r.PUT("/fyc/UserAudit", pkg.UpdateUserAuditAPI)
	r.DELETE("/fyc/UserAudit", pkg.DeleteUserAuditAPI)
}
