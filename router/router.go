package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "gopa/docs"
	m "gopa/gorm"
	"gopa/handler/api"
	"gopa/handler/sd"
	"gopa/router/middleware"
	"net/http"
)

func InitEngine() *gin.Engine {
	g := gin.New()
	m.DB.Init()
	return g
}

// Load loads the middlewares, routes, handler.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// The health check handler
	svcdRouter := g.Group("/sd")
	{
		svcdRouter.GET("/health", sd.HealthCheck)
		svcdRouter.GET("/version", sd.VersionCheck)
		svcdRouter.GET("/disk", sd.DiskCheck)
		svcdRouter.GET("/cpu", sd.CPUCheck)
		svcdRouter.GET("/ram", sd.RAMCheck)
	}
	gapi := g.Group("/api")
	// 登录接口
	gapi.POST("sso-login", api.JwtAuth().LoginHandler)
	gapi.GET("sso-logout", api.JwtAuth().LogoutHandler)
	// GOPA 第一版API
	v1 := gapi.Group("/v1")
	v1.Use(api.JwtAuth().MiddlewareFunc())
	// 管理组的API
	v1.GET("auth", api.Auth)
	groupAPIs := v1.Group("/project")
	{
		groupAPIs.GET("/list", api.ProjectsList)
		groupAPIs.POST("/add", api.ProjectAdd)
		groupAPIs.POST("/delete", api.ProjectDelete)
	}
	// 管理角色的API
	roleAPIs := v1.Group("/role")
	{
		roleAPIs.GET("/list", api.RolesList)
		roleAPIs.POST("/add", api.RoleAdd)
		roleAPIs.POST("/delete", api.RoleDelete)
	}
	// 管理组内角色的API
	projectRolesAPIs := v1.Group("/projectRole")
	{
		projectRolesAPIs.GET("/list", api.ProjectRoleList)
		projectRolesAPIs.POST("/add", api.ProjectRoleAdd)
		projectRolesAPIs.POST("/delete", api.ProjectRoleDelete)
	}
	// 管理用户所属组以及角色的API
	userAPIs := v1.Group("/user")
	{
		userAPIs.GET("/list", api.UserList)
		userAPIs.POST("/add", api.UserAdd)
		userAPIs.POST("/delete", api.UserDelete)
		userAPIs.POST("/update", api.UserUpdate)
	}
	// 管理策略文件的API
	regoAPIs := v1.Group("/rego")
	{
		regoAPIs.GET("/list", api.RegoList)
		regoAPIs.POST("/add", api.RegoAdd)
		regoAPIs.POST("/update", api.RegoUpdate)
		regoAPIs.POST("/delete", api.RegoDelete)
	}
	// 管理应用的API，如perf-server, crawling-server
	applicationAPIs := v1.Group("/application")
	{
		applicationAPIs.GET("/list", api.ApplicationList)
		applicationAPIs.POST("/add", api.ApplicationAdd)
		applicationAPIs.POST("delete", api.ApplicationDelete)
	}
	// 管理应用下的组和角色
	projectResourcesAPIs := v1.Group("/projectResource")
	{
		projectResourcesAPIs.GET("/list", api.ProjectResourceList)
		projectResourcesAPIs.POST("/add", api.ProjectResourceAdd)
		projectResourcesAPIs.POST("/update", api.ProjectResourceUpdate)
		projectResourcesAPIs.POST("/delete", api.ProjectResourceDelete)
	}

	v2 := gapi.Group("/v2")
	{
		v2.GET("/test", api.Test)
	}
	return g
}
