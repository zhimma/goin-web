package jwtLibrary

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/zhimma/goin-web/database/structure"
	globalInstance "github.com/zhimma/goin-web/global"
	"net/http"
	"strings"
	"time"
)

type JWT struct {
	SigningKey []byte
	Expires    int64
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenMalformed   = errors.New("Token is malformed")
	TokenNotValidYet = errors.New("Token is not valid yet")
	TokenInvalid     = errors.New("Token is invalid")
)

func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(globalInstance.BaseConfig.Jwt.JwtSecret),
		Expires:    globalInstance.BaseConfig.Jwt.JwtTtl,
	}
}

// 生产jwt token
func (j *JWT) GenerateJwtToken(identifier interface{}) (*structure.JwtTokenDetails, error) {
	// token detail struct
	tokenDetail := &structure.JwtTokenDetails{}
	// 赋值过期时间
	tokenDetail.Expires = j.Expires
	//  赋值uuid
	accessTokenClaims := j.BuildClaims(identifier)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, accessTokenErr := accessToken.SignedString(j.SigningKey)
	if accessTokenErr != nil {
		return nil, accessTokenErr
	}
	// 赋值access_token
	tokenDetail.Token = accessTokenString
	return tokenDetail, nil
}

// 解析token
func (j *JWT) ParseJwtToken(tokenString string) (*structure.JwtClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &structure.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*structure.JwtClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

// 构造jwt claims数据
func (j *JWT) BuildClaims(identifier interface{}) structure.JwtClaims {
	claimsData := structure.JwtClaims{
		IDENTIFIER: identifier,
		StandardClaims: jwt.StandardClaims{
			// Audience:  "",
			// Id:        "",
			// IssuedAt:  0,
			// Subject:   "",
			NotBefore: time.Now().Unix() - 1000,      // 签名生效时间
			ExpiresAt: time.Now().Unix() + j.Expires, // 过期时间 7天
			Issuer:    "goin-web-admin",
		},
	}
	return claimsData
}

// 从request中提取token
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
