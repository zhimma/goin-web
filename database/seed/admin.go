package seed

import (
	"fmt"
	"github.com/zhimma/goin-web/database/model"
	"github.com/zhimma/goin-web/helper"
	"gorm.io/gorm"
	"time"
)

func Admin(db *gorm.DB) {

	user := model.Admin{
		Account:     "zhimma",
		Salt:        "1256",
		Password:    "123456",
		Avatar:      "",
		Name:        "马雄飞",
		Phone:       "18710830000",
		Status:      1,
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
	if result := db.Create(&user); result.Error != nil {
		fmt.Printf("填充数据失败：「%v」\n", result.Error)
	}
	fmt.Printf("「%v」数据表填充完成\n", "admin")
}
