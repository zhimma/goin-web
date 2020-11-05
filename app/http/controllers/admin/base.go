package admin

import (
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate
var trans ut.Translator

func init() {
	binding.Validator.Engine()

	validate = validator.New()
}
