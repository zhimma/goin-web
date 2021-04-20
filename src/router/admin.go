package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/app/http/controllers/admin"
	"github.com/zhimma/goin-web/app/http/controllers/admin/api_group"
	apis "github.com/zhimma/goin-web/app/http/controllers/admin/apis"
	"github.com/zhimma/goin-web/app/http/controllers/admin/casbin_auth"
	"github.com/zhimma/goin-web/app/http/controllers/admin/casbin_auth/role"
	"github.com/zhimma/goin-web/app/http/controllers/admin/category"
	"github.com/zhimma/goin-web/app/http/controllers/admin/client_passport"
	"github.com/zhimma/goin-web/app/http/controllers/admin/managers"
	"github.com/zhimma/goin-web/app/http/controllers/admin/passport"
	"github.com/zhimma/goin-web/app/middleware"
)

func InitAdminRouter(Router *gin.RouterGroup) {
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
	categoryRouter := Router.Group("/categories").Use(middleware.AdminAuth(), middleware.AdminCasbin())
	{
		categoryRouter.GET("", category.Index)
		categoryRouter.POST("", category.Store)
		categoryRouter.GET(":id", category.Show)
		categoryRouter.PUT(":id", category.Update)
		categoryRouter.DELETE(":id", category.Destroy)
	}

	// 接口管理
	// apisRouter := Router.Group("/apis").Use(middleware.AdminAuth())
	apisRouter := Router.Group("/apis").Use(middleware.AdminAuth(), middleware.AdminCasbin())
	{
		apisRouter.GET("/", apis.Index)
		apisRouter.GET("/:id", apis.Show)
		apisRouter.POST("/", apis.Store)
		apisRouter.PUT("/:id", apis.Update)
		apisRouter.DELETE("/:id", apis.Destroy)
	}
	// 接口组管理
	apiGroupRouter := Router.Group("api/groups").Use(middleware.AdminAuth(), middleware.AdminCasbin())
	{
		apiGroupRouter.GET("", api_group.Index)
		apiGroupRouter.POST("", api_group.Store)
		apiGroupRouter.PUT(":id", api_group.Update)
		apiGroupRouter.DELETE(":id", api_group.Destroy)
	}

	// 角色管理
	roleRouter := Router.Group("roles").Use(middleware.AdminAuth(), middleware.AdminCasbin())
	{
		roleRouter.GET("/", role.Index)
		roleRouter.POST("/", role.Store)
		roleRouter.PUT("/:id", role.Update)
		roleRouter.DELETE(":id", role.Destroy)
	}
	// 角色分配权限
	cabinAuth := Router.Group("role/policy").Use(middleware.AdminAuth())
	{
		cabinAuth.POST("/apis", casbin_auth.StoreApiPolicies)
		cabinAuth.POST("/menus", casbin_auth.StoreMenuPolicies)
	}

	// 管理员管理
	managersRouter := Router.Group("managers").Use(middleware.AdminAuth())
	{
		managersRouter.GET("/", managers.Index)
		managersRouter.POST("/", managers.Store)
		managersRouter.POST("/bind/role", managers.BindRole)
		managersRouter.PUT("/:id", managers.Update)
		managersRouter.DELETE(":id", managers.Destroy)
		managersRouter.POST("/change/password", managers.ChangePassword)

	}
}
