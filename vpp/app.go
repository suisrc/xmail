package vpp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vkc/conf"
	"vkc/core"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Option 定义配置项
type OptionSet func(*Options)

type ServerInfo interface {
	GetHttpSrv() http.Handler
	GetGrpcSrv() *grpc.Server
}

// BuildServer 创建服务
type BuildServer func() (serve ServerInfo, clean func(), err error)

// Options options
type Options struct {
	ConfigFile  []string
	Version     string
	BuildServer BuildServer
}

// SetConfigFile 设定配置文件
func SetServeConfig(s []string) OptionSet {
	return func(o *Options) {
		o.ConfigFile = s
	}
}

// SetVersion 设定版本号b
func SetServeVersion(s string) OptionSet {
	return func(o *Options) {
		o.Version = s
	}
}

// SetBuildEngine 设定注入助手
func SetBuildServer(b BuildServer) OptionSet {
	return func(o *Options) {
		o.BuildServer = b
	}
}

// RunServe 运行服务, 注意,必须对BuildInjector进行初始化
func RunServe(ctx context.Context, opts ...OptionSet) error {
	return ServeShutdown(ctx, func() (func(), error) {
		return ServeRun(ctx, opts...)
	})
}

// ServeRunAndShutdown 运行服务
func ServeShutdown(ctx context.Context, runServe func() (func(), error)) error {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	shutdownServe, err := runServe()
	if err != nil {
		return err
	}

	sig := <-sc // 等待服务器中断
	log.Printf("received a signal [%s], serve shutdown ...", sig.String())
	shutdownServe()
	log.Printf("application serve exiting")
	time.Sleep(time.Second) // 延迟1s, 用于日志等信息保存
	return nil
}

// ServeConfig 启动服务
func ServeRun(ctx context.Context, opts ...OptionSet) (func(), error) {
	var o Options
	for _, opt := range opts {
		opt(&o)
	}

	conf.MustLoad(o.ConfigFile...) // 加载配置文件
	conf.Print()                   // 配置打印
	conf.NextInit()                // 加载初始化方法

	// 启动日志
	log.Printf("http serve startup, M[%s]-V[%s]-P[%d]", conf.C.RunMode, o.Version, os.Getpid())

	// 初始化依赖注入器
	srv, cl1, err := o.BuildServer()
	if err != nil {
		return nil, err
	}

	// 初始化HTTP服务
	cl2 := RunHttpServe(ctx, srv.GetHttpSrv())
	cl3 := RunGrpcServe(ctx, srv.GetGrpcSrv())

	shutdown := func() { cl3(); cl3(); cl2(); cl1() }
	return shutdown, nil
}

// RunHttpServe 初始化http服务
func RunHttpServe(ctx context.Context, handler http.Handler) func() {
	if handler == nil {
		return func() {} // 忽略HTTP服务
	}
	conf := C.Serve
	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(conf.IdleTimeout) * time.Second,
	}

	go func() {
		log.Printf("http server is running at %s.", addr)
		// var err error
		// if conf.CertFile != "" && conf.KeyFile != "" {
		// 	srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		// 	err = srv.ListenAndServeTLS(conf.CertFile, conf.KeyFile)
		// } else {
		// 	err = srv.ListenAndServe()
		// }
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(conf.ShutdownTimeout))
		defer cancel()
		log.Printf("http server shutdown ...")
		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			core.ErrorWithStack1(err)
		}
	}
}

func RunGrpcServe(ctx context.Context, server *grpc.Server) func() {
	conf := C.Grpc
	if !conf.Enable || server == nil {
		return func() {} // 忽略GRPC服务
	}
	// http协议，这是实验性的
	if conf.Network == "http" {
		srv := &http.Server{
			Addr:         conf.Address,
			Handler:      server,
			ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(conf.IdleTimeout) * time.Second,
		}
		go func() {
			log.Printf("grpc server (http) is running at %s.", conf.Address)
			err := srv.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				panic(err)
			}
		}()
		return func() {
			log.Printf("grpc server shutdown ...")
			server.GracefulStop()
		}
	}
	// tcp or unix 协议
	srv, err := net.Listen(conf.Network, conf.Address) // tpc, :9090
	if err != nil {
		panic(err)
	}
	go func() {
		log.Printf("grpc server is running at %s.", conf.Address)
		reflection.Register(server)
		err := server.Serve(srv)
		if err != nil && !errors.Is(err, net.ErrClosed) {
			panic(err)
		}
	}()
	return func() {
		log.Printf("grpc server shutdown ...")
		server.GracefulStop()
	}
}
