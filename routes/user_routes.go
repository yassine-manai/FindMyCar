package routes

import (
	"github.com/gin-gonic/gin"

	"fmc/pkg"

)

func UserRoutes(r *gin.Engine) {
	r.GET("/fyc/users", pkg.GetAllUserApi)
	r.GET("/fyc/userEnabled", pkg.GetUserEnabledAPI)
	r.GET("/fyc/userDeleted", pkg.GetUserDeletedAPI)

	r.POST("/fyc/user", pkg.AddUserAPI)
	r.PUT("/fyc/user/:id", pkg.UpdateUserAPI)
	r.PUT("/fyc/userState", pkg.ChangeUserStateAPI)
	r.DELETE("/fyc/user", pkg.DeleteUserCredAPI)
}
