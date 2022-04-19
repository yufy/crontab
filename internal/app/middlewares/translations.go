package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
	zh_translation "github.com/go-playground/validator/v10/translations/zh"
)

func Translations() gin.HandlerFunc {
	return func(c *gin.Context) {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			locale := c.GetHeader("locale")
			uni := ut.New(zh.New(), en.New())

			trans, _ := uni.GetTranslator(locale)
			switch locale {
			case "zh":
				zh_translation.RegisterDefaultTranslations(v, trans)
			default:
				en_translation.RegisterDefaultTranslations(v, trans)
			}

			c.Set("trans", trans)
		}
		c.Next()
	}
}
