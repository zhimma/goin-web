package env

import (
	"flag"
	"fmt"
	"strings"
)

var nowEnv Environment
var local Environment = &environment{value: "local"}
var dev Environment = &environment{value: "dev"}
var uat Environment = &environment{value: "uat"}
var prod Environment = &environment{value: "prod"}

var _ Environment = (*environment)(nil)

type Environment interface {
	Value() string
	IsLocal() bool
	IsDev() bool
	IsUat() bool
	IsProd() bool
	e()
}

type environment struct {
	value string
}

func (e *environment) Value() string {
	return e.value
}

func (e *environment) IsLocal() bool {
	return e.value == "local"
}

func (e *environment) IsDev() bool {
	return e.value == "dev"
}

func (e *environment) IsUat() bool {
	return e.value == "uat"
}

func (e *environment) IsProd() bool {
	return e.value == "prod"
}

func (e *environment) e() {}

func init() {
	env := flag.String("env", "", "请输入运行环境:\n local:本地环境\n dev:开发环境\n uat:预上线环境\n pro:正式环境\n")
	flag.Parse()

	switch strings.ToLower(strings.TrimSpace(*env)) {
	case "local":
		nowEnv = local
	case "dev":
		nowEnv = dev
	case "uat":
		nowEnv = uat
	case "pro":
		nowEnv = prod
	default:
		nowEnv = dev
		fmt.Println("-env 未找到对应环境，已默认使用dev环境")
	}
}

func NowEnv() Environment {
	return nowEnv
}
