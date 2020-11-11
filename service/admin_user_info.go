package service

import (
	"fmt"
	"github.com/zhimma/goin-web/database/model"
	"github.com/zhimma/goin-web/database/structure"
	globalInstance "github.com/zhimma/goin-web/global"
)

func AdminUserInfo(id uint) {
	fmt.Println(id)
}

func AdminLogin(admin *model.Admin) (user model.Admin, err error) {
	var data model.Admin
	err = globalInstance.DB.Where("account = ?", admin.Account).Find(&data).Error
	return data, err
}

func CacheAdminUserToken(uid uint, tokenDetail *structure.JwtTokenDetails) {

}
