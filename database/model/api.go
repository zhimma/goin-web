package model

// 管理后台接口数据库
type Api struct {
	BaseModel
	ApiGroupId  int64  `json:"api_group_id" gorm:"default:0;comment:'接口组id';not null;index"`
	Path        string `json:"path" gorm:"default:'';comment:接口地址;not null;index"`
	Method      string `json:"method" gorm:"default:'';comment:方法;not null;index"`
	Description string `json:"description" gorm:"default:'';comment:'接口描述';not null"`
}
