package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/app/http/controllers/admin"
	apis "github.com/zhimma/goin-web/app/http/controllers/admin/apis"
	"github.com/zhimma/goin-web/app/http/controllers/admin/category"
	"github.com/zhimma/goin-web/app/http/controllers/admin/client_passport"
	"github.com/zhimma/goin-web/app/http/controllers/admin/passport"
	"github.com/zhimma/goin-web/app/middleware"
)

func InitAdminRouter(Router *gin.RouterGroup) {
	adminRouter := Router.Group("/passport")
	{
		adminRouter.POST("/login", passport.Login)
		adminRouter.POST("/register", passport.Register)
	}

	passportRouter := Router.Group("/client/passport")
	{
		passportRouter.POST("/apply", client_passport.Apply)
		passportRouter.POST("/auth", client_passport.Auth)
	}
	// 使用中间件
	Router.POST("/test", admin.TestList)

	// 分类
	categoryRouter := adminRouter.Group("/categories").Use(middleware.AdminAuth())
	{
		categoryRouter.GET("", category.Index)
		categoryRouter.POST("", category.Store)
		categoryRouter.GET(":id", category.Show)
		categoryRouter.PUT(":id", category.Update)
		categoryRouter.DELETE(":id", category.Destroy)
	}

	// 接口管理
	apisRouter := adminRouter.Group("/apis").Use(middleware.AdminAuth())
	{
		apisRouter.GET("", apis.Index)
		apisRouter.POST("", apis.Store)
		apisRouter.PUT(":id", apis.Update)
		apisRouter.DELETE(":id", apis.Destroy)
	}
	// 接口组管理
	apiGroupRouter := adminRouter.Group("api/groups").Use(middleware.AdminAuth())
	{
		apiGroupRouter.GET("", apis.Index)
		apiGroupRouter.POST("", apis.Store)
		apiGroupRouter.PUT(":id", apis.Update)
		apiGroupRouter.DELETE(":id", apis.Destroy)
	}
}
