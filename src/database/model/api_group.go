package model

// 管理后台接口数据库
type ApiGroup struct {
	BaseModel
	Name        string `json:"path" gorm:"default:'';comment:接口地址;not null;index"`
	Description string `json:"description" gorm:"default:'';comment:'接口描述';not null"`
	Status      int    `json:"status" gorm:"default:1;comment:'组启用状态:0-禁用，1-启用';not null"`
}
