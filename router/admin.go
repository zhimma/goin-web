package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/app/http/controllers/admin"
	"github.com/zhimma/goin-web/app/http/controllers/admin/category"
	"github.com/zhimma/goin-web/app/http/controllers/admin/passport"
	"github.com/zhimma/goin-web/app/middleware"
)

func InitAdminRouter(Router *gin.RouterGroup) {
	adminRouter := Router.Group("/admin")
	{
		adminRouter.POST("/login", admin.Login)
		adminRouter.POST("/register", admin.Register)
		adminRouter.POST("/logout", admin.Logout)
	}

	passportRouter := Router.Group("/passport")
	{
		passportRouter.POST("/apply", passport.Apply)
		passportRouter.POST("/auth", passport.Auth)
	}
	// 使用中间件
	testRouters := adminRouter.Use(middleware.AdminAuth())
	{
		testRouters.POST("/test", admin.TestList)
	}
	categoryRouter := adminRouter.Group("/categories")
	{
		categoryRouter.GET("", category.Index)
		categoryRouter.POST("", category.Store)
		categoryRouter.GET(":id", category.Show)
		categoryRouter.PUT(":id", category.Update)
		categoryRouter.DELETE(":id", category.Destroy)
	}

}
