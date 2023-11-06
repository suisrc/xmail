.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

APP = vkc
SERVER_BIN = ./docker/app

dev: start

# 初始化mod
init:
	go mod init ${APP}
#go mod init github.com/suisrc/${APP}

# 修正依赖
tidy:
	go mod tidy

# 打包应用 go build -gcflags=-G=3 -o $(SERVER_BIN) ./cmd
build:
	go build -ldflags "-w -s" -o $(SERVER_BIN) ./

# go get -u github.com/go-delve/delve/cmd/dlv
# 使用“--”将标志传递给正在调试的程
# 调试过程中,会产生一些僵尸进程,这个时候,可以通过杀死父进程解决
# ps -ef | grep defunct | more (kill -9 pid 是无法删除进程)
debug:
	dlv debug --headless --api-version=2 --listen=127.0.0.1:2345 main.go -- web -c ./shconf/config.toml
	#echo c | dlv debug --accept-multiclient --api-version=2 --listen=127.0.0.1:2345 main.go -- web -c ./shconf/config.toml

# 运行 -gcflags=-G=3
start:
	go run main.go web -c ./shconf/__config.toml

# go get -u github.com/google/wire/cmd/wire  |  go install github.com/google/wire/cmd/wire@latest
wire:
	wire gen ./wira

# go get -u github.com/nicksnyder/go-i18n/v2/goi18n  |  goi18n -help
i18n:
	goi18n extract -outdir locales -sourceLanguage zh-CN
i18n-update:
	cd locales && goi18n merge -sourceLanguage zh-CN active.zh-CN.toml translate.en-US.toml
i18n-merge:
	cd locales && goi18n merge -sourceLanguage en-US active.en-US.toml translate.en-US.toml

# go get -u github.com/mdempsky/gocode | grpc-proto3
code:
	gocode -s -debug

# =====================================================================================

test:
	@go test -v $(shell go list ./...)

pack: build
	rm -rf $(RELEASE_ROOT) && mkdir -p $(RELEASE_SERVER)
	cp -r $(SERVER_BIN) shconf/config.toml $(RELEASE_SERVER)
	cd $(RELEASE_ROOT) && tar -cvf $(APP).tar ${APP} && rm -rf ${APP}

deploy:
	git checkout deploy-auto && git merge master && git push && git checkout master