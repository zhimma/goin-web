package code

import "github.com/zhimma/goin-web/config"

type Failure struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

const (
	TooManyRequests = 10102
)

func Text(code int) string {
	lang := config.Get().Language.Local

	if lang == config.ZhCN {
		return zhCNText[code]
	}

	if lang == config.EnUS {
		return enUSText[code]
	}

	return zhCNText[code]
}
