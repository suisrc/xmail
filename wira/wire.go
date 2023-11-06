//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package wira

// 注意,该文件不会参与编译, 需要通过[make wire]命令生成wire_gen.go文件编译
// 文件第一行是通知golang编译器忽略该文件
/*
	系统使用google/wire作为框架的依赖注入组件,需要注意一下细节.

	1.系统没有自发现依赖组件,需要人工配置
	2.当依赖发生变更后,需要执行[make wire]命令更新
	3.当前可能会注入无用属性,比如[Swagger, Healthz], 该内容主要通知wire,执行构造方法
	4.Injector中的顺序,决定了整体的执行顺序
	5.wire_gen.go文件不要编辑,为自动生成
*/
import (
	"github.com/google/wire"
	"vkc/vpp"
)

// BuildInjector 生成注入器
func NewApp() (vpp.ServeRunner, func(), error) {
	wire.Build(AppSet)
	return new(App), nil, nil
}
