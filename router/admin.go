package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/app/http/controllers/admin"
	"github.com/zhimma/goin-web/app/http/controllers/admin/auth"
	"github.com/zhimma/goin-web/app/http/controllers/admin/category"
	"github.com/zhimma/goin-web/app/http/controllers/admin/passport"
	"github.com/zhimma/goin-web/app/middleware"
)

func InitAdminRouter(Router *gin.RouterGroup) {
	adminRouter := Router.Group("/admin")
	{
		adminRouter.POST("/login", auth.Login)
		adminRouter.POST("/register", auth.Register)
	}

	passportRouter := Router.Group("/passport")
	{
		passportRouter.POST("/apply", passport.Apply)
		passportRouter.POST("/auth", passport.Auth)
	}
	// 使用中间件
	Router.POST("/test", admin.TestList)
	categoryRouter := adminRouter.Group("/categories").Use(middleware.AdminAuth())
	{
		categoryRouter.GET("", category.Index)
		categoryRouter.POST("", category.Store)
		categoryRouter.GET(":id", category.Show)
		categoryRouter.PUT(":id", category.Update)
		categoryRouter.DELETE(":id", category.Destroy)
	}

}
