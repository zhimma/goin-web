package model

import (
	"time"
)

type Manager struct {
	BaseModel
	Account     string    `json:"account" gorm:"type:varchar(100);default:'';not null;comment:'账号';index"`
	Salt        string    `json:"-" gorm:"type:varchar(10);default:'';not null;comment:'盐'"`
	Password    string    `json:"-" gorm:"type:varchar(100);default:'';not null;comment:'密码'"`
	Avatar      string    `json:"avatar" gorm:"type:varchar(150);default:'';not null;comment:'头像'"`
	Name        string    `json:"name" gorm:"type:varchar(150);default:'';not null;comment:'名称'"`
	Phone       string    `json:"phone" gorm:"type:varchar(20);default:'';not null;comment:'手机号';index"`
	Email       string    `json:"email" gorm:"type:varchar(50);default:'';not null;comment:'邮箱';index"`
	Status      int8      `json:"status" gorm:"comment:状态：1-启用，0-禁用;default:0;index"`
	LastLoginIp string    `json:"last_login_ip" gorm:"comment:最后登陆ip;default:'';not null"`
	LastLoginAt time.Time `json:"last_login_at" gorm:"comment:最后登陆时间'"`
	LoginTimes  int64     `json:"login_times" gorm:"comment:登陆次数;default:0;not null"`
	OperatorId  int64     `json:"operator_id" gorm:"comment:添加人id;default:0;not null"`
}
