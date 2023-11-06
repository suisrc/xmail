package vpp

import (
	"log"
	"os"
	"path/filepath"
	"time"
	"vkc/conf"
	"vkc/core"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	g18n "github.com/suisrc/gin-i18n"
	"google.golang.org/grpc"
)

// EmptyMiddleware 不执行业务处理的中间件
func EmptyMiddleware(ctx *gin.Context) {
	ctx.Next() // Pass, 跳过
}

//===================================================================
//===================================================================
//===================================================================

// UseHttpSrv 绑定中间件
func UseHttpSrv(bundle *g18n.Bundle) HttpSrvOption {
	cnf1 := conf.C
	conf := C
	return func(app *gin.Engine) {
		if cnf1.RunMode != "release" {
			app.Use(gin.Logger()) // LOGGER
			//app.Use(LoggerMiddleware())
		}
		// CORS
		if conf.Cors.Enable {
			app.Use(MiddlewareCors())
		}
		// GZIP
		if conf.Gzip.Enable {
			app.Use(MiddlewareGzip())
		}
		// 异常处理
		//app.Use(gin.Recovery())
		app.Use(MiddlewareRecovery)
		// 国际化, 全部国际化
		app.Use(MiddlewareI18n(bundle))
	}
}

// UseHttpSrv 绑定中间件
func UseHttpSrv0() HttpSrvOption {
	cnf1 := conf.C
	conf := C
	return func(app *gin.Engine) {
		if cnf1.RunMode != "release" {
			app.Use(gin.Logger()) // LOGGER
			//app.Use(LoggerMiddleware())
		}
		// CORS
		if conf.Cors.Enable {
			app.Use(MiddlewareCors())
		}
		// GZIP
		if conf.Gzip.Enable {
			app.Use(MiddlewareGzip())
		}
		// 异常处理
		//app.Use(gin.Recovery())
		app.Use(MiddlewareRecovery)
	}
}

// UseGrpcSrv 绑定中间件
func UseGrpcSrv0() GrpcSrvOption {
	return func() []grpc.ServerOption {
		return []grpc.ServerOption{}
	}
}

// ===================================================================
// MiddlewareI18n 国际化
func MiddlewareI18n(bundle *g18n.Bundle) gin.HandlerFunc {
	// bundle := i18n.NewBundle(language.Chinese)
	// bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	// bundle.LoadMessageFile("locales/active.zh-CN.toml")
	// bundle.LoadMessageFile("locales/active.en-US.toml")
	return g18n.Serve(bundle)
}

// ===================================================================
// MiddlewareGzip Giz, 主要部署前端时候(www中间件)对静态资源进行压缩
func MiddlewareGzip() gin.HandlerFunc {
	conf := C.Gzip
	if !conf.Enable {
		return EmptyMiddleware
	}
	return gzip.Gzip(gzip.BestCompression,
		gzip.WithExcludedExtensions(conf.ExcludedExtentions),
		gzip.WithExcludedPaths(conf.ExcludedPaths),
	)
}

// ===================================================================
// MiddlewareCors 跨域
func MiddlewareCors() gin.HandlerFunc {
	conf := C.Cors
	if !conf.Enable {
		return EmptyMiddleware
	}
	return cors.New(cors.Config{
		AllowAllOrigins:        conf.AllowAllOrigins,
		AllowOrigins:           conf.AllowOrigins,
		AllowMethods:           conf.AllowMethods,
		AllowHeaders:           conf.AllowHeaders,
		ExposeHeaders:          conf.ExposeHeaders,
		AllowCredentials:       conf.AllowCredentials,
		MaxAge:                 time.Second * time.Duration(conf.MaxAge),
		AllowWildcard:          conf.AllowWildcard,
		AllowBrowserExtensions: conf.AllowBrowserExtensions,
		AllowWebSockets:        conf.AllowWebSockets,
		AllowFiles:             conf.AllowFiles,
	})
}

// ===================================================================
// MiddlewareSwww 静态站点中间件
func MiddlewareSwww(root string) gin.HandlerFunc {
	conf := C.Swww
	if !conf.Enable {
		return EmptyMiddleware
	}
	return func(c *gin.Context) {
		if root == "" {
			root = conf.RootDir
		}

		p := c.Request.URL.Path
		fpath := filepath.Join(root, filepath.FromSlash(p))
		_, err := os.Stat(fpath)
		if err != nil && os.IsNotExist(err) {
			fpath = filepath.Join(root, conf.Index)
		}

		c.File(fpath)
		c.Abort()
	}
}

// ===================================================================
// MiddlewareRecovery 崩溃恢复中间件
func MiddlewareRecovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			stack := core.Stack(3)
			log.Printf("[panic]: %v\n%v", stack, err)
			core.ResError(c, core.Err500InternalServer)
			// 结束
		}
	}()
	c.Next()
}
