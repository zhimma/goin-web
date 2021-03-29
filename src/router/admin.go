package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/app/http/controllers/admin"
	"github.com/zhimma/goin-web/app/http/controllers/admin/api_group"
	apis "github.com/zhimma/goin-web/app/http/controllers/admin/apis"
	"github.com/zhimma/goin-web/app/http/controllers/admin/casbin/role"
	"github.com/zhimma/goin-web/app/http/controllers/admin/category"
	"github.com/zhimma/goin-web/app/http/controllers/admin/client_passport"
	"github.com/zhimma/goin-web/app/http/controllers/admin/passport"
)

func InitAdminRouter(Router *gin.RouterGroup) {
	Router = Router.Group("admin")
	// 管理后台登陆
	adminRouter := Router.Group("/passport")
	{
		adminRouter.POST("/login", passport.Login)
		adminRouter.POST("/register", passport.Register)
	}
	// service client 登陆
	passportRouter := Router.Group("/client/passport")
	{
		passportRouter.POST("/apply", client_passport.Apply)
		passportRouter.POST("/auth", client_passport.Auth)
	}
	// 使用中间件
	Router.POST("/test", admin.TestList)

	// 分类
	categoryRouter := Router.Group("/categories")
	{
		categoryRouter.GET("", category.Index)
		categoryRouter.POST("", category.Store)
		categoryRouter.GET(":id", category.Show)
		categoryRouter.PUT(":id", category.Update)
		categoryRouter.DELETE(":id", category.Destroy)
	}

	// 接口管理
	// apisRouter := Router.Group("/apis").Use(middleware.AdminAuth())
	apisRouter := Router.Group("/apis")
	{
		apisRouter.GET("", apis.Index)
		apisRouter.POST("", apis.Store)
		apisRouter.PUT(":id", apis.Update)
		apisRouter.DELETE(":id", apis.Destroy)
	}
	// 接口组管理
	apiGroupRouter := Router.Group("api/groups")
	{
		apiGroupRouter.GET("", api_group.Index)
		apiGroupRouter.POST("", api_group.Store)
		apiGroupRouter.PUT(":id", api_group.Update)
		apiGroupRouter.DELETE(":id", api_group.Destroy)
	}

	// 接口组管理
	casbinRole := Router.Group("roles")
	{
		casbinRole.GET("", role.Index)
		casbinRole.POST("", role.Store)
		casbinRole.PUT(":id", role.Update)
		casbinRole.DELETE(":id", role.Destroy)
	}
}
