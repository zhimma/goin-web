package model

import (
	"time"
)

type Admin struct {
	BaseModel
	Account     string    `json:"account" gorm:"type:varchar(100);default:'';comment:'账号';index"`
	Salt        string    `json:"salt" gorm:"type:varchar(10);default:'';comment:'盐'"`
	Password    string    `json:"_" gorm:"type:varchar(100);default:'';comment:'密码'"`
	Avatar      string    `json:"avatar" gorm:"type:varchar(150);default:'';comment:'头像'"`
	Name        string    `json:"name" gorm:"type:varchar(150);default:'';comment:'名称'"`
	Phone       string    `json:"phone" gorm:"type:varchar(20);default:'';comment:'手机号';index"`
	Status      int8      `json:"status" gorm:"comment:状态：1-启用，0-禁用"`
	LastLoginIp string    `json:"last_login_ip" gorm:"comment:最后登陆ip"`
	LastLoginAt time.Time `json:"last_login_at" gorm:"comment:最后登陆时间"`
	LoginTimes  int64     `json:"login_time" gorm:"comment:登陆次数"`
}
