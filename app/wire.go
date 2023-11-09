package app

import (
	"vkc/vpp"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// app package

// WireSet wire注入声明
var WireSet = wire.NewSet(
	vpp.NewRouter, // router
	// wire.Struct(new(Demo), "*"),
	// wire.Struct(new(Index), "*"),
	wire.Struct(new(Email), "*"),

	// wire.Struct(new(apis.HelloworldServer), "*"),
	wire.Struct(new(Injector), "*"), // 注册器
)

type Injector struct {
	//==========================================
	// Http 引擎
	Engine *gin.Engine // 引擎
	Router gin.IRouter // 根路由

	// Http API
	// Index *Index
	// Demo  *Demo
	Email *Email

	//==========================================
	// Grpc 引擎
	// Server *grpc.Server // 引擎

	// Grpc SVC
	// HelloworldServer *apis.HelloworldServer
}

// Init 初始化
func (aa *Injector) PostInit() (func(), error) {
	// Http API
	// engine := aa.Engine
	// aa.Index.Register(engine) // 首页

	router := aa.Router
	// aa.Demo.Register(router)
	aa.Email.Register(router)

	// Grpc API
	// server := aa.Server
	// aa.HelloworldServer.Register(server)

	return func() {}, nil
}
