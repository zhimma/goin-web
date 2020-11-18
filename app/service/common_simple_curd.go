/**

 * @Author: zhimma
 * @Description:
 * @File:  common_simple_curd
 * @Date: 2020/11/18 11:08
 */
package service

import (
	globalInstance "github.com/zhimma/goin-web/global"
)

func Create(data interface{}) (err error) {
	err = globalInstance.DB.Create(data).Error
	return
}
