/**

 * @Author: zhimma
 * @Description:
 * @File:  common_simple_curd
 * @Date: 2020/11/18 11:08
 */
package CommonDbService

import (
	"errors"
	"fmt"
	globalInstance "github.com/zhimma/goin-web/global"
	"gorm.io/gorm"
)

// 分页参数和条件
type PageStruct struct {
	Page         int                    `json:"page"`
	PageSize     int                    `json:"page_size"`
	MapWhere     map[string]interface{} `json:"map_where"`
	LikeMapWhere map[string]interface{} `json:"like_map_where"`
}

// 分页meta参数
type MetaResult struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

// 分页结果
type PageResult struct {
	List interface{} `json:"list"`
	Meta MetaResult  `json:"meta"`
}

// 列表
func List(where map[string]interface{}, list interface{}) (err error) {
	err = globalInstance.DB.Where(where).Find(list).Error
	return err
}
func ListByIdSlice(where []int64, list interface{}) (err error) {
	err = globalInstance.DB.Where(where).Find(list).Error
	return err
}

// 创建一条数据
func InsertOne(data interface{}) (err error) {
	err = globalInstance.DB.Create(data).Error
	return err
}

// 根据id获取详情
func DetailById(data interface{}, id int64) (err error) {
	err = globalInstance.DB.First(data, id).Error
	return err
}

// 根据id更新
func UpdateById(model interface{}, id int64, data interface{}) (number int64, err error) {
	err1 := DetailById(model, id)
	if err1 != nil {
		return 0, err1
	}
	if errors.Is(err1, gorm.ErrRecordNotFound) {
		return 0, errors.New("记录不存在")
	}
	number = globalInstance.DB.Model(model).Updates(data).RowsAffected
	return number, nil
}

// 根据id删除
func DeleteById(model interface{}) (err error) {
	err = globalInstance.DB.Delete(model).Error
	return nil
}

// 物理删除
func ForceDelete(model interface{}) (err error) {
	err = globalInstance.DB.Unscoped().Delete(model).Error
	return nil
}

// 根据条件查询详情
func DetailByMapWhere(where map[string]interface{}, model interface{}) error {
	return globalInstance.DB.Where(where).First(model).Error
}

// firstorcreate
func FirstOrCreate(where map[string]interface{}, model interface{}) error {
	return globalInstance.DB.FirstOrCreate(model, where).Error
}

// 分页
func toPage(page int, pageSize int, where map[string]interface{}, likeWhereString string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		if len(where) > 0 {
			return db.Where(where).Where(likeWhereString).Offset(offset).Limit(pageSize)
		}
		return db.Offset(offset).Where(likeWhereString).Limit(pageSize)
	}
}

// 分页
func Paginate(where PageStruct, model interface{}) (result PageResult) {
	if where.Page == 0 {
		where.Page = 1
	}
	switch {
	case where.PageSize > 1000:
		where.PageSize = 1000
	case where.PageSize <= 0:
		where.PageSize = 10
	}
	result = PageResult{
		List: model,
		Meta: MetaResult{
			Page:     where.Page,
			PageSize: where.PageSize,
		},
	}
	var likeWhereString string
	if len(where.LikeMapWhere) > 0 {
		lenNum := len(where.LikeMapWhere)
		for k, v := range where.LikeMapWhere {
			lenNum--
			vStr := fmt.Sprintf("%v", v)
			if len(vStr) <= 0 {
				continue
			}
			if lenNum > 0 {
				likeWhereString += k + " LIKE '%" + vStr + "%' AND "
			} else {
				likeWhereString += k + " like '%" + vStr + "%' "
			}
		}
	}
	globalInstance.DB.Scopes(toPage(where.Page, where.PageSize, where.MapWhere, likeWhereString)).Find(result.List)
	globalInstance.DB.Model(model).Where(where.MapWhere).Where(likeWhereString).Count(&result.Meta.Total)
	return result
}
