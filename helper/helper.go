package helper

import (
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

/*func GenerateJwtString(data interface{}) (token string, err error) {
	c := MyClaims{
		"username", // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "my-project",                               // 签发人
		},
	}
	tokenData := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenData.SignedString(globalInstance.BaseConfig.Jwt.JwtSecret)
	return
}*/

func ParseJwtString() {

}

func ValidateJwtString() {

}
