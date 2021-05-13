package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/component"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/global/response"
	"io/ioutil"
	"time"
)

type Params struct {
	Name string `json:"name" form:"name" `
	Age  int    `json:"age" form:"age"`
}

func TestList(c *gin.Context) {
	cache := component.NewRedisCache()
	cache.Set("name", []string{"maxsf"}, 0)
	err := cache.Set("age", 22, 10*time.Second)
	fmt.Println(err)
	var age int
	err = cache.Get("age", &age)
	fmt.Println(age, err)
	type person struct {
		Name string
		Age  int
	}
	p := person{
		Name: "zhimma",
		Age:  30,
	}

	fmt.Println("client_num", globalInstance.RedisClientNum)

	cache.Set("p", p, 0)
	res := person{}
	err = cache.Get("p", &res)
	fmt.Println(res, err)

	var data []string
	err = cache.Get("name", &data)
	fmt.Println(data, err)
	//var path = "/Users/zhimma/go/src/goin-web/app/http/controllers"
	//readDir(path, 0)
	var params Params
	if err := c.ShouldBind(&params); err != nil {
		response.ValidateFail(err.Error(), c)
		return
	}
	response.OkWithData(params, c)
	return
}
func readDir(path string, curHier int) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range files {
		if v.IsDir() {
			for tmpHier := curHier; tmpHier > 0; tmpHier-- {
				fmt.Printf("|\t")
			}
			fmt.Println(v.Name(), curHier)
			readDir(path+"/"+v.Name(), curHier+1)
		} else {
			for tmpHier := curHier; tmpHier > 0; tmpHier-- {
				fmt.Printf("|\t")
			}
			fmt.Println(v.Name())
		}
	}
}
