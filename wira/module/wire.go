package module

import (
	"github.com/google/wire"
)

// CasbinSet wire注入服务
var WireSet = wire.NewSet(
	// vpp.MiddlewareI18n, // 国际化
	// mgo.NewDefaultDatabase, // 数据库

	wire.Struct(new(Injector), "*"), // 注册器
)

type Injector struct {
}

// Init 初始化
func (o *Injector) PostInit() (func(), error) {
	return func() {}, nil
}
