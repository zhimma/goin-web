package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/app/middleware"
)

func Routers() *gin.Engine {
	var Router = gin.New()
	Router.Use(middleware.Cors())
	AdminGroup := Router.Group("admin")
	// 初始化所有的路由
	InitAdminRouter(AdminGroup)
	// 初始化所有的路由
	return Router
}
