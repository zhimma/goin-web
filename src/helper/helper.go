package helper

import (
	"fmt"
	"github.com/zhimma/goin-web/database/model"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
	"time"
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

// 生产hash字符
func GenerateHashString(content string, salt string) (string, error) {
	ps := content + salt
	hash, err := bcrypt.GenerateFromPassword([]byte(ps), bcrypt.DefaultCost)
	return string(hash), err

}

// 比较hash字符串
func CompareHashString(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}
	return true
}

// validator去除tag
func RemoveTopStruct(fields map[string]string) []string {
	// resMap := map[string]string{}
	var resSlice []string
	for _, err := range fields {
		fmt.Println(err)
		// resMap[field[strings.Index(field, ".")+1:]] = err
		resSlice = append(resSlice, err)
	}
	return resSlice
}

// 生成随机字符
func RandStringBytes(n int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func List2Tree(data []model.Category, pk interface{}, pid uint, child interface{}, root uint) (result map[interface{}]interface {
}) {
	result = map[interface{}]interface{}{}
	for _, v := range data {
		result[v.ID] = v
	}
	return result

}
