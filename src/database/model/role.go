package model

// 管理后台接口数据库
type Role struct {
	BaseModel
	Name   string `json:"name" gorm:"default:'';comment:接口地址;not null;index"`
	Status int    `json:"status" gorm:"default:1;comment:'启用状态:0-禁用，1-启用';not null"`
	Type   int    `json:"type" gorm:"default:2;comment:'类型：1-超级管理员，2-普通管理员';not null"`
}
