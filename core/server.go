package core

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	globalInstance "github.com/zhimma/goin-web/global"
	"github.com/zhimma/goin-web/router"
	"time"
)

type server interface {
	ListenAndServe() error
}

func Run() {
	router := router.Routers()
	address := fmt.Sprintf(":%d", globalInstance.BaseConfig.System.Addr)
	fmt.Printf("服务开始运行，地址为:「%v」\n", address)
	// 加载数据验证器
	s := serverStart(address, router)
	fmt.Println(s.ListenAndServe().Error())
}

func serverStart(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}
