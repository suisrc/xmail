package xmail

// 邮件管理器

import (
	"context"
	"io"
	"strings"
	"time"
	"vkc/core"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jhillyerd/enmime"
)

type MailManager struct {
	DS *mongo.Database

	ST SyncMailTask // 同步任务信息
}

func (aa *MailManager) Coll() *mongo.Collection {
	if C.XMail.CollName == "" {
		return aa.DS.Collection("mails")
	} else {
		return aa.DS.Collection(C.XMail.CollName)
	}
}

// ==================================================================

type GetEmailsCO struct {
	Zone  string `form:"zone"`  // 邮箱域
	Addr  string `form:"addr"`  // 邮箱地址
	Skip  int    `form:"skip"`  // 跳过多少条, 默认0条
	Limit int    `form:"limit"` // 限制多少条, 默认10条
	Show  bool   `form:"show"`  // 是否显示邮件内容
}

func (aa *MailManager) GetEmails2(ctx context.Context, co *GetEmailsCO) ([]*Mail, error) {
	fitler := bson.M{}
	if co.Zone != "" {
		fitler["zone"] = co.Zone
	} else if co.Addr != "" {
		fitler["to"] = co.Addr
	}
	// 获取邮件
	sorter := bson.M{"date": -1}
	option := options.Find().SetSort(sorter)
	if co.Skip > 0 {
		option.SetSkip(int64(co.Skip))
	}
	if co.Limit > 0 {
		option.SetLimit(int64(co.Limit))
	} else {
		option.SetLimit(10) // 默认10条
	}
	option.SetSort(bson.M{"date": -1}) // 按时间倒序
	// 查询邮件
	cursor, err := aa.Coll().Find(ctx, fitler, option)
	if err != nil {
		return nil, core.ErrorOf("S_EML_GET_FIND", "Error: "+err.Error())
	}
	defer cursor.Close(ctx)
	result := []*Mail{}
	for cursor.Next(ctx) {
		item := &Mail{}
		err := cursor.Decode(item)
		if err != nil {
			continue
		}
		result = append(result, item)
	}

	// 返回邮件
	return result, nil
}

// ==================================================================

type GetEmailCO struct {
	Zone  string `form:"zone"`  // 邮箱域
	Addr  string `form:"addr"`  // 邮箱地址
	MsgId string `form:"mid"`   // 邮件ID
	State int    `form:"state"` // 邮件状态, 1: 未读, 2: 已读, 0: 所有
}

func (aa *MailManager) GetEmail2(ctx context.Context, co *GetEmailCO) (*Mail, error) {
	fitler := bson.M{}
	if co.Zone != "" {
		fitler["zone"] = co.Zone
	} else if co.Addr != "" {
		fitler["to"] = co.Addr
	}
	if co.MsgId != "" {
		fitler["mid"] = co.MsgId // 指定邮件
		// logrus.Info("mid: ", co.MsgId)
	}
	if co.State > 0 {
		fitler["state"] = co.State // 指定邮件

	}
	option := options.FindOneAndUpdate().SetSort(bson.M{"date": -1})
	update := bson.M{"$set": bson.M{"state": 2}}
	result := aa.Coll().FindOneAndUpdate(ctx, fitler, update, option)
	if result.Err() != nil {
		return nil, core.ErrorOf("S_EML_GET_FIND", "Error: "+result.Err().Error())
	}

	rst := &Mail{}
	err := result.Decode(rst)
	if err != nil {
		return nil, core.ErrorOf("S_EML_GET_DECODE", "Error: "+err.Error())
	}

	// 返回邮件
	return rst, nil
}

// ==================================================================

func (aa *MailManager) InsertEmail2(ctx context.Context, bts io.Reader) (*Mail, error) {
	// ctx.GetRawData() // 获取原始数据
	eml, err := enmime.ReadEnvelope(bts)
	if err != nil {
		return nil, core.ErrorOf("S_EML_MAIL_INFO", "Error: "+err.Error())
	}
	date0 := eml.GetHeader("Date")
	if date0 == "" {
		return nil, core.ErrorOf("S_EML_MAIL_INFO", "Date is required")
	}
	data := &Mail{
		MsgId:   eml.GetHeader("Message-ID"),
		From:    eml.GetHeader("From"),
		To:      eml.GetHeader("To"),
		Subject: eml.GetHeader("Subject"),
		Html:    eml.HTML,
		Text:    eml.Text,
		State:   1,
	}
	if data.MsgId == "" || data.From == "" || data.To == "" || data.Subject == "" {
		return nil, core.ErrorOf("S_EML_MAIL_INFO", "Message-ID, From, To, Subject is required")
	}
	if strings.HasPrefix(data.MsgId, "<") && strings.HasSuffix(data.MsgId, ">") {
		data.MsgId = data.MsgId[1 : len(data.MsgId)-1]
	}
	data.MsgId = strings.ReplaceAll(data.MsgId, "+", "-")

	date, err := time.Parse("Mon, _2 Jan 2006 15:04:05 -0700", date0)
	if err != nil {
		date, err = time.Parse("Mon, _2 Jan 2006 15:04:05 -0700 (MST)", date0)
		if err != nil {
			return nil, core.ErrorOf("S_EML_MAIL_INFO", "Date is invalid: "+err.Error())
		}
	}
	data.Date = date
	data.SetZone1()
	// 插入数据
	_, err = aa.Coll().InsertOne(ctx, data)
	if err != nil {
		return nil, core.ErrorOf("S_EML_MAIL_INFO", "Insert Error: "+err.Error())
	}
	return data, nil
}
