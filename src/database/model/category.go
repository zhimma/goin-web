package model

type Category struct {
	BaseModel
	Pid         int64  `json:"pid" gorm:"default:0;comment:'pid'"`
	Name        string `json:"name" gorm:"type:varchar(100);default:'';comment:'分类名称';index"`
	Description string `json:"description" gorm:"type:varchar(255);default:'';comment:'描述'"`
}
