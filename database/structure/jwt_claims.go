// 此文件用于定于jwt 数据的 struct集合
package structure

import "github.com/dgrijalva/jwt-go"

type JwtTokenDetails struct {
	Token   string
	Uuid    string
	Expires int64
}

type JwtClaims struct {
	IDENTIFIER interface{}
	UUID       string
	jwt.StandardClaims
}
