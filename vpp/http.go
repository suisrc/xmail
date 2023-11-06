package vpp

import (
	"vkc/conf"
	"vkc/core"

	"github.com/gin-gonic/gin"
)

// ===================================================================
// HttpSrvOption 修正http内容
type HttpSrvOption func(*gin.Engine)

// NewHttpSrv engine
func NewHttpSrv(opt HttpSrvOption) *gin.Engine {
	gin.SetMode(conf.C.RunMode)
	//gin.SetMode(gin.DebugMode)

	app := gin.New()
	app.NoMethod(NoMethodHandler)
	app.NoRoute(NoRouteHandler)

	opt(app)
	//app.Use(gin.Logger())
	//app.Use(middleware.LoggerMiddleware())
	//app.Use(gin.Recovery())
	//app.Use(middleware.RecoveryMiddleware())

	return app
}

// NewRouter 初始化根路由
func NewRouter(app *gin.Engine) gin.IRouter {
	var router gin.IRouter
	if v := C.Serve.ContextPath; v != "" {
		router = app.Group(v)
	} else {
		router = app
	}

	return router
}

// NoMethodHandler 未找到请求方法的处理函数
func NoMethodHandler(c *gin.Context) {
	core.ResError(c, core.Err405MethodNotAllowed)
	// Abort, 终止
}

// NoRouteHandler 未找到请求路由的处理函数
func NoRouteHandler(c *gin.Context) {
	core.ResError(c, core.Err404NotFound)
	// Abort, 终止
}
