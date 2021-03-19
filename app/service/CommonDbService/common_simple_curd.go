/**

 * @Author: zhimma
 * @Description:
 * @File:  common_simple_curd
 * @Date: 2020/11/18 11:08
 */
package CommonDbService

import (
	globalInstance "github.com/zhimma/goin-web/global"
)

func List(where map[string]interface{}, list interface{}) (err error) {
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
	if err := DetailById(model, id); err != nil {
		return 0, err
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

func FirstOrCreate(where map[string]interface{}, model interface{}) error {
	return globalInstance.DB.FirstOrCreate(model, where).Error
}
