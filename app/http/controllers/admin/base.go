package admin

import (
	"github.com/dgrijalva/jwt-go"
	globalInstance "github.com/zhimma/goin-web/global"
	"time"
)

var StandardClaims jwt.StandardClaims

func NewStandardClaims() {
	StandardClaims = jwt.StandardClaims{
		//Audience:  "",
		//Id:        "",
		//IssuedAt:  0,
		//Subject:   "",
		NotBefore: time.Now().Unix() - 1000,                                 // 签名生效时间
		ExpiresAt: time.Now().Unix() + globalInstance.BaseConfig.Jwt.JwtTtl, // 过期时间 7天
		Issuer:    "goin-web-admin",
	}
}
