package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/app/http/controllers/admin"
	"github.com/zhimma/goin-web/app/middleware"
)

func InitAdminRouter(Router *gin.RouterGroup) {
	AdminRouter := Router.Group("/admin")
	{
		AdminRouter.POST("/login", admin.Login)
		AdminRouter.POST("/register", admin.Register)
		AdminRouter.POST("/logout", admin.Logout)
	}
	// 使用中间件
	AdminRouter.Use(middleware.AdminAuth())
	{
		AdminRouter.POST("/test", admin.TestList)
	}
}
