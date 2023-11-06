package vpp

import (
	"context"
	"net/http"
	"os"
	"vkc/core"
	"vkc/shell"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var (
	DEBUG = false
)

// RunApp(app, CreateRunServe(inj))

// ServeCommand ...
func ServeCommand(ctx context.Context, action func(c *cli.Context) error) *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "运行web服务",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:        "conf",
				Aliases:     []string{"c"},
				Usage:       "配置文件(.json,.yaml,.toml)",
				DefaultText: "config.toml",
				//Required:   true,
			},
		},
		Action: action,
	}
}

//======================================================================================

func RunApp(app *cli.App, run func(ctx context.Context, opts ...OptionSet) error) {
	ctx := context.Background()
	atn := func(c *cli.Context) error {
		sc := SetServeConfig(c.StringSlice("conf"))
		sv := SetServeVersion(app.Version)
		return run(ctx, sc, sv)
	}
	cmd := ServeCommand(ctx, atn)
	cmds := shell.GetCommands()

	app.Commands = append([]*cli.Command{cmd}, cmds...)
	err := app.Run(os.Args)
	if err != nil {
		core.ErrorWithStack1(err)
		//panic(err)
	}
}

//========================================================================================

type ServeRunner interface {
	PostInit() (func(), error)
	GetHttpSrv() http.Handler
	GetGrpcSrv() *grpc.Server
}

type ServeRunnerBuilder func() (ServeRunner, func(), error)

// CreateRunServe 运行服务
func CreateRunServe(builder ServeRunnerBuilder) func(ctx context.Context, opts ...OptionSet) error {
	return func(ctx context.Context, opts ...OptionSet) error {
		iopt := SetBuildServer(func() (ServerInfo, func(), error) {
			runner, cleaner, err := builder()
			if err != nil {
				return nil, nil, err
			}
			cl2, err := runner.PostInit()
			if err != nil {
				cleaner() // 清理
				return nil, nil, err
			}
			cl3 := func() { cl2(); cleaner() }
			return runner, cl3, err
		})
		return RunServe(ctx, append(opts, iopt)...)
	}
}
