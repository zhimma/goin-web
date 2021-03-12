package jwtLibrary

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-basic/uuid"
	"github.com/zhimma/goin-web/database/structure"
	globalInstance "github.com/zhimma/goin-web/global"
	"net/http"
	"strings"
	"time"
)

type JWT struct {
	SigningKey          []byte
	AccessTokenExpires  int64
	RefreshTokenExpires int64
}

const (
	AccessToken  = 1
	RefreshToken = 2
)

var (
	TokenExpired     = errors.New("Token is expired")
	TokenMalformed   = errors.New("Token is malformed")
	TokenNotValidYet = errors.New("Token is not valid yet")
	TokenInvalid     = errors.New("Token is invalid")
)

func NewJWT() *JWT {
	return &JWT{
		SigningKey:          []byte(globalInstance.BaseConfig.Jwt.JwtSecret),
		AccessTokenExpires:  globalInstance.BaseConfig.Jwt.JwtTtl,
		RefreshTokenExpires: globalInstance.BaseConfig.Jwt.JwtTtl,
	}
}

// 生产jwt token
func (j *JWT) GenerateJwtToken(Id int64) (*structure.JwtTokenDetails, error) {
	// token detail struct

	tokenDetail := &structure.JwtTokenDetails{}

	tokenDetail.AccessTokenExpires = j.AccessTokenExpires
	tokenDetail.RefreshTokenExpires = j.RefreshTokenExpires

	//  access token
	accessTokenClaims := j.BuildClaims(Id, AccessToken)
	tokenDetail.AccessTokenUuid = accessTokenClaims.UUID
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, accessTokenErr := accessToken.SignedString(j.SigningKey)

	if accessTokenErr != nil {
		return nil, accessTokenErr
	}
	tokenDetail.AccessToken = accessTokenString

	// refresh token
	refreshTokenClaims := j.BuildClaims(Id, RefreshToken)
	tokenDetail.RefreshTokenUuid = refreshTokenClaims.UUID
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, refreshTokenErr := refreshToken.SignedString(j.SigningKey)
	if refreshTokenErr != nil {
		return nil, refreshTokenErr
	}
	tokenDetail.RefreshToken = refreshTokenString

	return tokenDetail, nil
}

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
func (j *JWT) BuildClaims(Id int64, model int) structure.JwtClaims {
	var expiresAt int64
	if model == AccessToken {
		expiresAt = j.AccessTokenExpires
	} else if model == RefreshToken {
		expiresAt = j.RefreshTokenExpires
	}
	claimsData := structure.JwtClaims{
		UID:  Id,
		UUID: uuid.New(),
		StandardClaims: jwt.StandardClaims{
			// Audience:  "",
			// Id:        "",
			// IssuedAt:  0,
			// Subject:   "",
			NotBefore: time.Now().Unix() - 1000,      // 签名生效时间
			ExpiresAt: time.Now().Unix() + expiresAt, // 过期时间 7天
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
