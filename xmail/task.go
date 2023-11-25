package xmail

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"vkc/core"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/sirupsen/logrus"
)

// ==================================================================
// 自动同步刷新邮件
type SyncMailTask struct {
	State SyncMailState // 同步任务信息
	CtlMS chan string   // 控制消息管道
}

type SyncMailState struct {
	Running bool `json:"running"`       // 是否正在运行
	MailCnt int  `json:"mail_count"`    // 已经同步的邮件数量
	ReqCnt  int  `json:"request_count"` // 请求同步的调用次数

	NextSyncAt time.Time `json:"next_sync_at"` //   下一次同步的时间
	LastSyncAt time.Time `json:"last_sync_at"` // 最后一次同步的时间
	LastError  string    `json:"last_error"`   // 最后一次同步的错误

	StartedAt time.Time `json:"started_at"` // 开始时间， 最后一次开始的时间
	StoppedAt time.Time `json:"stopped_at"` // 停止时间， 最后一次停止的时间
	CreatedAt time.Time `json:"created_at"` // 创建时间，   第一次启动的时间

}

func (aa *MailManager) SyncEmailGet(ctx *gin.Context) {
	core.ResSuccess(ctx, aa.ST.State) // 返回状态
}

func (aa *MailManager) SyncEmailCtl(ctx *gin.Context) {
	active := ctx.Query("active") // 是否运行
	if active == "1" && !aa.ST.State.Running {
		go aa.SyncEmailRun()
		core.ResSuccess(ctx, "started") // 已经启动
	} else if active == "1" && aa.ST.State.Running {
		core.ResSuccess(ctx, "running") // 正在运行
	} else if active == "0" && aa.ST.State.Running {
		aa.ST.CtlMS <- "stop"
		core.ResSuccess(ctx, "stoping") // 正在停止
	} else if active == "0" && !aa.ST.State.Running {
		core.ResSuccess(ctx, "stopped") // 已经停止
	} else {
		core.ResError(ctx, core.Err400BadParam)
	}
}

// 推荐邮件同步地址每次只读取一封邮件，以免在同步过程中，前面封邮件没有成功导致后面的邮件无法同步
// 但是如果效率第一的情况下，推荐一次性读取多封邮件，以提高效率
func (aa *MailManager) SyncEmailRun() {
	if aa.ST.State.Running {
		return // 正在运行, 不重复启动
	}
	if C.XMail.SyncUri == "" {
		aa.ST.State.LastError = "no sync address: mail.sync: https://..."
		return // 没有配置同步地址
	}

	if aa.ST.State.CreatedAt.IsZero() {
		aa.ST.State.CreatedAt = time.Now()
	}
	// 启动
	aa.ST.State.StartedAt = time.Now()
	aa.ST.CtlMS = make(chan string, 1)
	aa.ST.State.Running = true
	// 终止
	defer func() {
		aa.ST.State.Running = false
		close(aa.ST.CtlMS) // 关闭管道
		aa.ST.State.StoppedAt = time.Now()

		if err := recover(); err != nil { // 发生意外错误
			aa.ST.State.LastError = fmt.Sprintf("%v", err)
		}
	}()
	// 同步业务
	for aa.ST.State.Running {
		aa.SleepMailSync(1, "", false) // 每次同步间隔
		aa.ST.State.LastSyncAt = time.Now()
		aa.ST.State.ReqCnt++ // 请求次数

		// 获取远程同步数据
		data := []byte{}
		code := 0
		err := gout.GET(C.XMail.SyncUri).BindBody(&data).Code(&code).Do()
		sec := C.XMail.SyncSec
		if err != nil {
			aa.SleepMailSync(sec, "request network err: "+err.Error(), true)
			continue // 网络错误
		}
		if code != 200 {
			aa.SleepMailSync(sec, "request content err: code -> "+strconv.Itoa(code), true)
			continue // 内容错误
		}
		result := &EmailTemp{}
		if err := json.Unmarshal(data, &result); err != nil {
			aa.SleepMailSync(sec, "request unmarshal err: "+err.Error(), true)
			continue // 解析错误
		}
		if len(result.Data) == 0 {
			aa.SleepMailSync(sec, "", true)
			continue // 没有新邮件
		}

		// 将邮件持久化到数据库
		ctx := context.TODO()
		for _, item := range result.Data {
			aa.ST.State.MailCnt++ // 已经同步的邮件数量
			eml, err := aa.InsertEmail2(ctx, strings.NewReader(item.Raw))
			if err != nil {
				logrus.Error("insert email err: ", err.Error(), " raw: \n", item.Raw)
			} else {
				go aa.SyncMailNotice(eml) // 邮件已经保存，异步通知其他模块
			}
		}
	}
}

func (aa *MailManager) SleepMailSync(sec int, err string, nxt bool) {
	if err != "" {
		// 记录错误
		aa.ST.State.LastError = err
	}
	if nxt {
		// 下一次同步时间
		aa.ST.State.NextSyncAt = time.Now().Add(time.Duration(sec) * time.Second)
	}
	// 同步中断
	select {
	case ctl := <-aa.ST.CtlMS:
		if ctl == "stop" {
			logrus.Info("sync mail stop")
			aa.ST.State.Running = false // 停止
		} else {
			logrus.Info("sync mail ctl: ", ctl)
		}
	case <-time.After(time.Duration(sec) * time.Second):
		// logrus.Info("sync mail sleep: ", sec, "s")
	}
}

type EmailTemp struct {
	Success bool `json:"success,omitempty"`
	Data    []struct {
		Id  string `json:"id,omitempty"`
		Raw string `json:"raw,omitempty"`
	} `json:"data,omitempty"`
}
