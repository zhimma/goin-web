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
	ID   uint
	UUID string
	jwt.StandardClaims
}
