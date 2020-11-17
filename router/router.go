package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/app/middleware"
)

func Routers() *gin.Engine {
	var Router = gin.New()
	Router.Use(middleware.Cors())
	AdminGroup := Router.Group("")
	InitAdminRouter(AdminGroup)
	return Router
}
