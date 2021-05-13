package helper

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/database/model"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
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

func GetCurrentManagerInfo(c *gin.Context) (data model.Manager, err error) {
	data = model.Manager{}
	managerInfo, exist := c.Get("managerInfo")
	if !exist {
		return data, errors.New("用户登陆状态获取失败或获取UID出错[not exists]")
	}
	userInfo, ok := managerInfo.(model.Manager)
	if ok {
		userInfo = managerInfo.(model.Manager)
	} else {
		return data, errors.New("用户登陆状态获取失败或获取UID出错[type error]")
	}
	data.ID = userInfo.ID
	data.RoleId = userInfo.RoleId
	data.Account = userInfo.Account
	data.Salt = userInfo.Salt
	data.Password = userInfo.Password
	data.Avatar = userInfo.Avatar
	data.Name = userInfo.Name
	data.Phone = userInfo.Phone
	data.Email = userInfo.Email
	data.Status = userInfo.Status
	data.IsSuper = userInfo.IsSuper
	data.LastLoginIp = userInfo.LastLoginIp
	data.LastLoginAt = userInfo.LastLoginAt
	data.LoginTimes = userInfo.LoginTimes
	data.OperatorId = userInfo.OperatorId
	return data, nil
}
func ValueInterfaceToString(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

//序列化
func Serialize(value interface{}) ([]byte, error) {
	if bytes, ok := value.([]byte); ok {
		return bytes, nil
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10)), nil
	}
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	if err := encoder.Encode(value); err != nil { //编码
		return nil, err
	}
	return b.Bytes(), nil
}

//反序列化
func Unserialize(byt []byte, ptr interface{}) (err error) {
	if b, ok := ptr.(*[]byte); ok {
		*b = byt

		return nil
	}

	if v := reflect.ValueOf(ptr); v.Kind() == reflect.Ptr { // 通过反射得到ptr类型，判断ptr是指针类型
		switch p := v.Elem(); p.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: //符号整型
			var i int64
			i, err = strconv.ParseInt(string(byt), 10, 64)
			if err != nil {
				return err
			} else {
				p.SetInt(i)
			}
			return nil

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: //无符号整型
			var i uint64
			i, err = strconv.ParseUint(string(byt), 10, 64)
			if err != nil {
				return err
			} else {
				p.SetUint(i)
			}
			return nil
		}
	}

	b := bytes.NewBuffer(byt)
	decoder := gob.NewDecoder(b)
	if err = decoder.Decode(ptr); err != nil { //解码
		return err
	}
	return nil
}
