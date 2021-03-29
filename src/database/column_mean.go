package database

import "gorm.io/gorm"

type ColumnTypeMean struct {
	gorm.Model
	// int 和 unit 区别在于 是否勾选无符号，uint == 勾选无符号
	// int == bigint
	StatusInt int `gorm:"column:status_int"`
	// int8 == tinyint
	StatusInt8 int8 `gorm:"column:status_int8"`
	// int16 == smallint
	StatusInt16 int16 `gorm:"column:status_int16"`
	// int32 == int
	StatusInt32 int32 `gorm:"column:status_int32"`
	// int64 == bigint
	StatusInt64   int64 `gorm:"column:status_int64"`
	StatusIntSize int64 `gorm:"column:status_int64;size:10"`
	// string 默认191个长度 DefaultStringSize配置参数
	StringColumn     string `gorm:"column:string_column_rename;size:50"`
	StringColumnSize string `gorm:"column:string_column_size;comment:这是本字段的描述;type:varchar(100);not null"`
	// 创建索引
	IndexColumn string `gorm:"column:index_col;index"`
	// text字段 多个tag存在
	TextColumn string `gorm:"type:text" json:"content"`
	// 设置默认值
	DefaultColumn string `gorm:"default:''" json:"default_value"`
}
