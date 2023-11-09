package xmail

import (
	"vkc/core"

	"github.com/gin-gonic/gin"
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
	core.ResSuccess(ctx, eml)
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
	ctx.String(200, eml.Html)
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

// 接受邮件通知
func (aa *MailManager) SyncMailNotice(eml *Mail) {
	// do nohing
}
