// 数据填充文件
package seed

import (
	"fmt"
	"github.com/zhimma/goin-web/app/service/CommonDbService"
	"github.com/zhimma/goin-web/database/model"
	"github.com/zhimma/goin-web/helper"
	"gorm.io/gorm"
	"time"
)

func Manager(db *gorm.DB) {
	user := model.Manager{
		Account:     "zhimma",
		Salt:        "1256",
		Password:    "123456",
		Email:       "mma5694@gmail.com",
		Avatar:      "",
		Name:        "马雄飞",
		Phone:       "18710830000",
		Status:      1,
		IsSuper:     1,
		LastLoginIp: "127.0.0.1",
		LastLoginAt: time.Now(),
		LoginTimes:  0,
	}
	// https://my.oschina.net/u/4296417/blog/4329615
	hash, err := helper.GenerateHashString(user.Password, user.Salt)
	if err != nil {
		fmt.Printf("密码加密失败：「%v」\n", err)
	}
	user.Password = hash
	where := map[string]interface{}{
		"account": user.Account,
		"email":   user.Email,
	}
	if errors := CommonDbService.FirstOrCreate(where, &user); errors != nil {
		fmt.Printf("填充数据失败：「%v」\n", errors)
	}
	fmt.Printf("「%v」数据表填充完成\n", "admin")
}
