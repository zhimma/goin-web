// 此文件用于定于jwt 数据的 struct集合
package structure

import "github.com/dgrijalva/jwt-go"

type JwtTokenDetails struct {
	AccessToken         string
	RefreshToken        string
	AccessTokenUuid     string
	RefreshTokenUuid    string
	AccessTokenExpires  int64
	RefreshTokenExpires int64
}

type JwtClaims struct {
	UID  int64
	UUID string
	jwt.StandardClaims
}
