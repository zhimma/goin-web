package service

import (
	"fmt"
	"github.com/zhimma/goin-web/database/structure"
)

func AdminUserInfo(id uint) {
	fmt.Println(id)
}

func CacheUserToken(uid uint, tokenDetail *structure.JwtTokenDetails) {

}
