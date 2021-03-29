package model

// 客户端数据库表
type Client struct {
	BaseModel
	ClientId       int64  `json:"client_id" gorm:"default:0;comment:客户端id;not null;index"`
	Salt           string `json:"salt" gorm:"type:varchar(100);default:'';comment:'盐';not null"`
	ClientSecret   string `json:"client_secret" gorm:"default:'';comment:'客户端secret';not null"`
	Status         int8   `json:"status" gorm:"default:0;comment:状态：0-未启用，1-启用，2-禁用;not null;index"`
	ContactPhone   string `json:"contact_phone" gorm:"type:varchar(20);default:'';comment:联系电话;not null;index"`
	ContactAddress string `json:"contact_address" gorm:"default:'';comment:联系地址;not null"`
}
