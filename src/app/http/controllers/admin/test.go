package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhimma/goin-web/global/response"
	"io/ioutil"
)

type Params struct {
	Name string `json:"name" form:"name" `
	Age  int    `json:"age" form:"age"`
}

func TestList(c *gin.Context) {

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
