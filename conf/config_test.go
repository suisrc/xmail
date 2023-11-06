package conf_test

import (
	"log"
	"testing"

	"vkc/conf"
)

func TestConfig(t *testing.T) {
	conf.CS = append(conf.CS, &TestLogger{})
	log.Println("===========================loading")
	conf.MustLoad("config_test.toml")
	conf.Print()
	// assert.NotNil(t, nil) // 异常才能显示日志
}

// MiddleConfig 中间件启动和关闭
type TestLogger struct {
	Logger  bool `default:"true"`
	Recover bool
}
