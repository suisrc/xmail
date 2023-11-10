package xmail

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strconv"
	"sync"
	"time"
	"vkc/core"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

// 获取多个邮件
func (aa *MailManager) GetEmails(ctx *gin.Context) {
	co := &GetEmailsCO{}
	if err := ctx.ShouldBindQuery(co); err != nil {
		core.ResError2(ctx, err, core.Err400BadParam)
		return // 参数错误
	}
	eml, err := aa.GetEmails2(ctx, co)
	if err != nil {
		core.Fix500Logger(ctx, err)
		return // 服务器错误
	}
	if !co.Show {
		// 不显示邮件内容
		for _, em := range eml {
			em.Html = ""
			em.Text = ""
		}
	}
	core.ResSuccess(ctx, eml)
}

// 获取多个邮件
func (aa *MailManager) GetEmailsText(ctx *gin.Context) {
	co := &GetEmailsCO{}
	if err := ctx.ShouldBindQuery(co); err != nil {
		core.ResError2(ctx, err, core.Err400BadParam)
		return // 参数错误
	}
	emls, err := aa.GetEmails2(ctx, co)
	if err != nil {
		core.Fix500Logger(ctx, err)
		return // 服务器错误
	}
	ctx.Header("Content-Type", "text/plain; charset=utf-8")
	bts := bytes.Buffer{}
	bts.WriteString("Number of email: " + strconv.Itoa(len(emls)) + "\n")
	for idx, eml := range emls {
		stridx := strconv.Itoa(idx + 1)
		bts.WriteString(stridx + " >> ======================================================================\n")
		bts.WriteString("Msg-Id: " + eml.MsgId + "\n")
		bts.WriteString("Msg-Date: " + eml.Date.Format(time.RFC1123Z) + "\n")
		bts.WriteString("From: " + eml.From + "\n")
		bts.WriteString("To: " + eml.To + "\n")
		bts.WriteString("Subject: " + eml.Subject + "\n")
		bts.WriteString("--------------------------------------------------------------------------\n")
		bts.WriteString(eml.Text + "\n")
		bts.WriteString(stridx + " << ======================================================================\n")
	}
	ctx.String(200, bts.String())
	ctx.Abort()
}

// 获取单个邮件
func (aa *MailManager) GetEmail(ctx *gin.Context) {
	co := &GetEmailCO{}
	if err := ctx.ShouldBindQuery(co); err != nil {
		core.ResError2(ctx, err, core.Err400BadParam)
		return // 参数错误
	}
	eml, err := aa.GetEmail2(ctx, co)
	if err != nil {
		core.Fix500Logger(ctx, err)
		return // 服务器错误
	}
	core.ResSuccess(ctx, eml)
}

// 获取单个邮件的HTML
func (aa *MailManager) GetEmailHtml(ctx *gin.Context) {
	co := &GetEmailCO{}
	if err := ctx.ShouldBindQuery(co); err != nil {
		core.ResError2(ctx, err, core.Err400BadParam)
		return // 参数错误
	}
	eml, err := aa.GetEmail2(ctx, co)
	if err != nil {
		core.Fix500Logger(ctx, err)
		return // 服务器错误
	}
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Header("Msg-Id", eml.MsgId)
	ctx.Header("Msg-Date", eml.Date.Format(time.RFC1123Z))
	ctx.Header("From", eml.From)
	ctx.Header("To", eml.To)
	ctx.Header("Subject", url.QueryEscape(eml.Subject))
	ctx.String(200, eml.Html)
	ctx.Abort()
}

// 获取单个邮件的HTML
func (aa *MailManager) GetEmailText(ctx *gin.Context) {
	co := &GetEmailCO{}
	if err := ctx.ShouldBindQuery(co); err != nil {
		core.ResError2(ctx, err, core.Err400BadParam)
		return // 参数错误
	}
	eml, err := aa.GetEmail2(ctx, co)
	if err != nil {
		core.Fix500Logger(ctx, err)
		return // 服务器错误
	}
	ctx.Header("Content-Type", "text/plain; charset=utf-8")
	ctx.Header("Msg-Id", eml.MsgId)
	ctx.Header("Msg-Date", eml.Date.Format(time.RFC1123Z))
	ctx.Header("From", eml.From)
	ctx.Header("To", eml.To)
	ctx.Header("Subject", url.QueryEscape(eml.Subject))
	ctx.String(200, eml.Text)
	ctx.Abort()
}

// 保存邮件
func (aa *MailManager) InsertEmail(ctx *gin.Context) {
	bts := ctx.Request.Body
	eml, err := aa.InsertEmail2(ctx, bts)
	if err != nil {
		core.Fix500Logger(ctx, err)
		return // 服务器错误
	}
	core.ResSuccess(ctx, eml.MsgId)
}

// 更新邮件, 内容过于简单，不提供单独更新方法
type UpdateEmailCO struct {
	MsgId string `json:"mid"`   // 邮件ID
	State int    `json:"state"` // 邮件状态, 1: 未读, 2: 已读
	// Zone  string `json:"zone"`  // 邮箱域
}

func (aa *MailManager) UpdateEmail(ctx *gin.Context) {
	co := &UpdateEmailCO{}
	if err := ctx.ShouldBindJSON(co); err != nil {
		core.ResError2(ctx, err, core.Err400BadParam)
		return // 参数错误
	}
	if co.MsgId == "" {
		core.ResError(ctx, core.ErrorOf("S_EML_PARAMS", "mid is required"))
		return // 参数错误
	}
	if co.State < 1 || co.State > 2 {
		core.ResError(ctx, core.ErrorOf("S_EML_PARAMS", "state is invalid"))
		return // 参数错误
	}
	fitler := bson.M{"mid": co.MsgId}
	update := bson.M{"$set": bson.M{"state": co.State}}
	// domail := &Mail{To: "demo@" + co.Zone}
	_, err := aa.Coll().UpdateOne(ctx, fitler, update)
	if err != nil {
		core.ResError2(ctx, err, core.Err500InternalServer)
		return // 服务器错误
	}
	core.ResSuccess(ctx, "ok") // 更新成功
}

// 删除邮件, 内容过于简单，不提供单独删除方法
type DeleteEmailCO struct {
	MsgId string `form:"mid"` // 邮件ID
	// Zone  string `form:"zone"` // 邮箱域
}

func (aa *MailManager) DeleteEmail(ctx *gin.Context) {
	co := &DeleteEmailCO{}
	if err := ctx.ShouldBindQuery(co); err != nil {
		core.ResError2(ctx, err, core.Err400BadParam)
		return // 参数错误
	}
	if co.MsgId == "" {
		core.ResError(ctx, core.ErrorOf("S_EML_PARAMS", "mid is required"))
		return // 参数错误
	}
	// if co.Zone == "" {
	// 	core.ResError(ctx, core.ErrorOf("S_EML_PARAMS", "zone is required"))
	// 	return // 参数错误
	// }
	fitler := bson.M{"mid": co.MsgId}
	// domail := &Mail{To: "demo@" + co.Zone}
	_, err := aa.Coll().DeleteOne(ctx, fitler)
	if err != nil {
		core.ResError2(ctx, err, core.Err500InternalServer)
		return // 服务器错误
	}
	core.ResSuccess(ctx, "ok") // 删除成功
}

//============================================================================

var WsMap = sync.Map{}

type SyncEmailCO struct {
	Addr string `form:"addr"` // 邮箱地址
	Zone string `form:"zone"` // 邮箱域
	Html bool   `form:"html"` // 是否显示邮件内容
	Text bool   `form:"text"` // 是否显示邮件内容
}

type SyncEmailHL func(Mail) // 同步邮件回调, Mail是副本， 避免影响原始数据

// 获取同步状态
func (aa *MailManager) SyncEmailWs(ctx *gin.Context, wss *websocket.Upgrader) {
	co := SyncEmailCO{}
	if err := ctx.ShouldBindQuery(&co); err != nil {
		core.ResError(ctx, core.Err400BadParam) // 参数错误
		return
	}
	// 升级连接为ws
	ccc, err := wss.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Errorf("websocket upgrade error: %v", err)
		core.ResError(ctx, core.ErrorOf("S_WS_UPGRADE-ERR", "websocket upgrade error:"+err.Error()))
		return
	}
	defer ccc.Close() // 关闭连接

	// 发布邮件， 没有条件
	eh1 := func(eml Mail) {
		if !co.Html {
			eml.Html = ""
		}
		if !co.Text {
			eml.Text = ""
		}
		bts, _ := json.Marshal(eml)
		ccc.WriteMessage(websocket.TextMessage, bts)
	}

	// 发布邮件, 根据条件
	var ehl SyncEmailHL
	if co.Addr != "" {
		ehl = func(eml Mail) {
			if eml.To == co.Addr {
				eh1(eml)
			}
		}
	} else if co.Zone != "" {
		ehl = func(eml Mail) {
			if eml.Zone == co.Zone {
				eh1(eml)
			}
		}
	} else {
		ehl = eh1
	}
	WsMap.Store(ccc, ehl)   // 保存连接
	defer WsMap.Delete(ccc) // 删除连接

	done := make(chan int)
	go func() { // 监听上传的消息
		defer close(done)
		for {
			msgtype, message, err := ccc.ReadMessage()
			if err != nil {
				logrus.Error("read:", err.Error())
				return
			}
			if msgtype == websocket.PingMessage {
				// 不会处理， 被上级拦截器处理了
				ccc.WriteMessage(websocket.PongMessage, []byte("pong"))
			}
			logrus.Infof("recv: %s", string(message))
		}
	}()
	// 保持连接
	for {
		select {
		case <-done:
			return // 连接断开，读取异常
		case <-ctx.Done():
			return // 连接断开, 链接关闭
		case <-time.After(time.Second * 20):
			ccc.WriteMessage(websocket.PingMessage, []byte("ping")) // 保持连接
		}
	}
}

// 接受邮件通知
func (aa *MailManager) SyncMailNotice(eml *Mail) {
	WsMap.Range(func(key, value interface{}) bool {
		if ehl, ok := value.(SyncEmailHL); ok {
			ehl(*eml) // 发布邮件
		}
		return true
	})
}
