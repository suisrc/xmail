package xmail

import (
	"strings"
	"time"
)

// 邮件实体

type Mail struct {
	MsgId   string    `json:"mid,omitempty" bson:"mid,omitempty"`
	Zone    string    `json:"zone,omitempty" bson:"zone,omitempty"`
	From    string    `json:"from,omitempty" bson:"from,omitempty"`
	To      string    `json:"to,omitempty" bson:"to,omitempty"`
	Subject string    `json:"subject,omitempty" bson:"subject,omitempty"`
	Date    time.Time `json:"date,omitempty" bson:"date,omitempty"`
	Html    string    `json:"html,omitempty" bson:"html,omitempty"`
	Text    string    `json:"text,omitempty" bson:"text,omitempty"`
	State   int       `json:"state,omitempty" bson:"state,omitempty"` // 状态
}

type MailRaw struct {
	Mail
	Raw []byte `json:"raw,omitempty" bson:"raw_data,omitempty"` // raw data
}

// Zone return the zone of the mail
func (aa *Mail) SetZone1() string {
	if aa.To == "" {
		return ""
	}
	idx := strings.Index(aa.To, "@") + 1
	end := strings.Index(aa.To[idx:], ">")
	if end < 0 {
		return aa.To[strings.Index(aa.To, "@")+1:]
	} else {
		return aa.To[idx : idx+end]
	}
}

// func (aa *Mail) Coll() string {
// 	// return strings.ReplaceAll(aa.Zone(), ".", "_")
// 	return vpp.C.MailsColl
// }
