package conf

var (
	// C 全局配置(需要先执行MustLoad，否则拿不到配置)
	C  = new(Config)
	CS = []interface{}{}
	FS = []func(){}
)

func init() {
	CS = append(CS, C)
}

// Config 配置参数
type Config struct {
	RunMode     string `default:"release"`
	PrintConfig bool   `default:"false"`
	Middleware  MiddleConfig
}

// MiddleConfig 中间件启动和关闭
type MiddleConfig struct {
	LoggerFile string `default:""`
	Recover    bool
}

// IsDebugMode 是否是debug模式
func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}
