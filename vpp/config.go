package vpp

import (
	"vkc/conf"
)

var (
	C = new(Config)
)

func init() {
	conf.CS = append(conf.CS, C)
}

type Config struct {
	Serve     ServeConfig
	Grpc      GrpcConfig
	Swww      SwwwConfig
	Cors      CorsConfig
	Gzip      GzipConfig
	Gout      GoutConfig
	ServeName string `default:"vkc"` // 服务
}

// MiddleConfig 中间件启动和关闭
type ServeConfig struct {
	Host             string `default:"0.0.0.0"`
	Port             int    `default:"80"`
	MaxContentLength int64
	ContextPath      string `default:"api"`
	Prefixes         []string
	AuthzServer      string // 认证服务器， 子服务器需要认证权限时候调用
	AuthxServer      string // 认证服务器， 子服务器需要认证权限时候调用
	AuthRelogin      string // `default:"/api/iam/v1/a/odic/login"` // 重新登录地址
	DebugUserForce   string `default:"961212"` // 调试用户KEY
	DebugUserSkyKey  string `default:"X-Request-Sky-Authorize"`
	ReadTimeout      int    `default:"5"`
	WriteTimeout     int    `default:"10"`
	IdleTimeout      int    `default:"15"`
	ShutdownTimeout  int    `default:"10"`
}

type GrpcConfig struct {
	Enable          bool   `default:"false"`
	Network         string `default:"tcp"`
	Address         string `default:"0.0.0.0:9090"`
	ReadTimeout     int    `default:"5"`
	WriteTimeout    int    `default:"10"`
	IdleTimeout     int    `default:"15"`
	ShutdownTimeout int    `default:"10"`
}

// CorsConfig 跨域请求配置参数
type CorsConfig struct {
	Enable                 bool
	AllowAllOrigins        bool     // 允许所有的跨域请求
	AllowOrigins           []string // "*"， 允许访问所有域
	AllowMethods           []string // "POST,GET,OPTIONS,PUT,DELETE,UPDATE"
	AllowHeaders           []string // "Authorization,Content-Length,Accept,Origin,Host,Connection,Accept-Encoding,Accept-Language,Keep-Alive,User-Agent,Cache-Control,Content-Type,Pragma"
	ExposeHeaders          []string // "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma"
	AllowCredentials       bool     // "false"
	MaxAge                 int      `default:"86400"` // "86400"，一天
	AllowWildcard          bool     // origins like http://some-domain/*, https://api.* or http://some.*.subdomain.com
	AllowBrowserExtensions bool     // Allows usage of popular browser extensions schemas
	AllowWebSockets        bool     // Allows usage of WebSocket protocol
	AllowFiles             bool     // file:// schema (dangerous!) use it only when you 100% sure it's needed
}

// GzipConfig gzip压缩
type GzipConfig struct {
	Enable             bool
	ExcludedExtentions []string
	ExcludedPaths      []string
}

// SwwwConfig 静态资源
type SwwwConfig struct {
	Enable  bool
	Index   string `default:"index.html"`
	RootDir string `default:"www"`
}

// GoutConfig gzip压缩
type GoutConfig struct {
	Debug bool
}
