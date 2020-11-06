package helper

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/zhimma/goin-web/database/structure"
	globalInstance "github.com/zhimma/goin-web/global"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GenerateHashString(password string, salt string) (string, error) {
	ps := password + salt
	hash, err := bcrypt.GenerateFromPassword([]byte(ps), bcrypt.DefaultCost)
	return string(hash), err

}

// 生成jwt-token
func GenerateJwtToken(claims jwt.Claims) (token string, err error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString([]byte(globalInstance.BaseConfig.Jwt.JwtSecret))
	return
}

func ParseJwtToken(tokenString string) (*structure.AdminClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &structure.AdminClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(globalInstance.BaseConfig.Jwt.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*structure.AdminClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func ValidateJwtToken() {

}
