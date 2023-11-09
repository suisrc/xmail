package mgo

import "time"

type Token struct {
	Apikey    string    `bson:"apikey" json:"apikey"`         // 令牌
	Secret    string    `bson:"secret" json:"secret"`         // 签名密钥
	Permis    []string  `bson:"permis" json:"permis"`         // 权限
	Remark    string    `bson:"remark" json:"remark"`         // 备注
	CreatedAt time.Time `bson:"created_at" json:"created_at"` // 创建时间
}
