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
