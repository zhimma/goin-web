package model

type Article struct {
	BaseModel
	CategoryId    uint   `json:"category_id" gorm:"default:0;comment:'分类id';index"`
	Title         string `json:"title" gorm:"type:varchar(100);default:'';comment:'标题';index"`
	CoverImageUrl string `json:"cover_image_url" gorm:"type:varchar(255);default:'';comment:'封面图片'"`
	Author        string `json:"author" gorm:"type:varchar(100);default:'';comment:'作者';index"`
	Content       string `json:"avatar" gorm:"type:text;comment:'正文'"`
	OperatorId    uint   `json:"operator_id" gorm:"default:0;comment:'操作人id'"`
	Status        int8   `json:"status" gorm:"comment:状态：1-上架，0-下架"`
}
