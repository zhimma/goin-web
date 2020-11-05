package structure

import "github.com/dgrijalva/jwt-go"

type AdminClaims struct {
	ID       uint
	Username string
	NickName string
	jwt.StandardClaims
}
