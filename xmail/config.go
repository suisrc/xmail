package xmail

import (
	"vkc/conf"
)

var (
	C = new(Config)
)

func init() {
	conf.CS = append(conf.CS, C)
}

// Config 配置参数
type Config struct {
	XMail MailConfig
}

type MailConfig struct {
	CollName string `default:"xmail"`
	SyncUri  string // 同步邮件的地址
	SyncSec  int    `default:"20"` // 同步邮件的时间间隔
}
