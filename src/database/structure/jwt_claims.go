// 此文件用于定于jwt 数据的 struct集合
package structure

import "github.com/dgrijalva/jwt-go"

// JWT签名结构
type JwtTokenDetails struct {
	Token   string
	Expires int64
}

// 荷载 payload
type JwtClaims struct {
	IDENTIFIER int64
	jwt.StandardClaims
}
