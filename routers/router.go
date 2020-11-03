package routers

import (
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	var Router = gin.New()
	AdminGroup := Router.Group("")
	InitAdminRouter(AdminGroup)
	return Router
}
