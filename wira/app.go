package wira

import (
	"net/http"
	"vkc/app"
	"vkc/vpp"
	"vkc/wira/module"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"google.golang.org/grpc"
	// 引入swagger, 测试过程中， 使用restclient更方便一些， 先移除对swagger的支持
	// _ "github.com/suisrc/api/swagger"
)

// AppSet
var AppSet = wire.NewSet(
	module.WireSet, // 组件

	app.WireSet, // 业务

	vpp.NewHttpSrv,  //*gin.Engine
	vpp.UseHttpSrv0, // 增加引擎中间件
	vpp.NewHealthz,  // 健康检查

	wire.Struct(new(App), "*"),                 // 注册
	wire.Bind(new(vpp.ServeRunner), new(*App)), // 接口
)

//=====================================================================

// Injector 注入器(用于初始化完成之后的引用)
type App struct {
	Engine *gin.Engine // 引擎

	HLZ vpp.Healthz     // 健康
	MOD module.Injector // 组件
	APP app.Injector    // 业务

}

var _ vpp.ServeRunner = &App{}

func (that *App) GetHttpSrv() http.Handler {
	return that.Engine
}

func (that *App) GetGrpcSrv() *grpc.Server {
	return nil
}

// Construct 初始化
func (that *App) PostInit() (func(), error) {
	pcx := []func() (func(), error){
		that.MOD.PostInit,
		that.APP.PostInit,
	}
	clx := []func(){}
	clf := func() {
		for _, f := range clx {
			f()
		}
	}
	for _, f := range pcx {
		cln, err := f()
		if err != nil {
			defer clf()
			return nil, err
		}
		clx = append(clx, cln)
	}
	// clx = append(clx, grpcx.CloseAllConn)
	return clf, nil
}
