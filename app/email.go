package app

import (
	"math/rand"
	"time"
	"vkc/core"
	"vkc/mgo"
	"vkc/xmail"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Email 接口
type Email struct {
	EM *xmail.MailManager
}

// Register 注册路由
func (aa *Email) Register(r gin.IRouter) {
	r.GET("emls", aa.Auth, aa.EM.GetEmails)     // 获取多个邮件
	r.GET("eml", aa.Auth, aa.EM.GetEmail)       // 获取单个邮件
	r.POST("eml", aa.Auth, aa.EM.InsertEmail)   // 保存邮件
	r.PUT("eml", aa.Auth, aa.EM.UpdateEmail)    // 更新邮件
	r.DELETE("eml", aa.Auth, aa.EM.DeleteEmail) // 删除邮件

	r.GET("eml.html", aa.Auth, aa.EM.GetEmailHtml) // 获取单个邮件
	r.GET("eml.text", aa.Auth, aa.EM.GetEmailText) // 获取单个邮件

	r.GET("eml/sync", aa.Auth, aa.EM.SyncEmailGet)  // 获取同步状态
	r.POST("eml/sync", aa.Auth, aa.EM.SyncEmailCtl) // 配置同步任务

	r.GET("init_token", aa.token) // 初始化令牌
}

func (aa *Email) token(ctx *gin.Context) {
	rst := aa.EM.DS.Collection("token").FindOne(ctx, bson.M{})
	if rst.Err() == nil {
		core.ResError(ctx, core.ErrorOf("S_TKN_EXIST", "令牌已经初始化"))
		return
	}
	if rst.Err() != mongo.ErrNoDocuments {
		core.ResError1(ctx, rst.Err(), core.Err500InternalServer)
		return
	}
	sts := "0123456789abcdefghijklmnopqrstuvwxyz"
	sll := len(sts)
	bts := []byte{}
	for i := 0; i < 32; i++ {
		bts = append(bts, sts[rand.Intn(sll)])
	}
	// 保存令牌
	tkn := &mgo.Token{
		Apikey:    "api_" + string(bts), // 36位长度
		Secret:    "none",
		Permis:    []string{"*"},
		Remark:    "初始化令牌",
		CreatedAt: time.Now(),
	}
	_, err := aa.EM.DS.Collection("token").InsertOne(ctx, tkn)
	if err != nil {
		core.ResError1(ctx, err, core.Err500InternalServer)
		return
	}
	core.ResSuccess(ctx, "令牌初始化完成")
}

func (aa *Email) Auth(ctx *gin.Context) {
	// apikey := ctx.GetHeader("x-api-key")
	// if apikey == "" {
	// 	core.ResError(ctx, core.Err403Forbidden)
	// 	return
	// }
	// logrus.Debugf("apikey: %s", apikey)
	apikey := ""
	if ak := ctx.Query("ak"); ak != "" {
		apikey = ak
	} else if ak := ctx.GetHeader("x-api-key"); ak != "" {
		apikey = ak
	} else {
		core.ResError(ctx, core.Err403Forbidden)
		return
	}
	// logrus.Info("apikey: ", apikey)

	// 获取 apikey 对应的权限
	rst := aa.EM.DS.Collection("token").FindOne(ctx, bson.M{"apikey": apikey})
	if rst.Err() != nil {
		if rst.Err() == mongo.ErrNoDocuments {
			core.ResError(ctx, core.Err403Forbidden)
		} else {
			core.ResError1(ctx, rst.Err(), core.Err403Forbidden)
		}
		return
	}
	tkn := &mgo.Token{}
	if err := rst.Decode(tkn); err != nil {
		core.ResError1(ctx, err, core.Err403Forbidden)
		return
	}
	if tkn.Secret != "none" {
		core.ResError(ctx, core.Err403Forbidden)
		// 暂时不支持签名验证的方式
		// 暂时不支持签名验证的方式
		// 暂时不支持签名验证的方式
		return
	}

	// 权限校验
	pkey := ctx.Request.Method + " " + ctx.Request.URL.Path
	pchk := false
	for _, p := range tkn.Permis {
		if p == "*" || p == pkey {
			pchk = true
			break
		}
	}
	if !pchk {
		core.ResError(ctx, core.Err403Forbidden)
		return
	}

	// 权限校验通过
	ctx.Next()
}
