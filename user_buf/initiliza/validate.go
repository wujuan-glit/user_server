package initiliza

import (
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"user/user_buf/global"
	"user/user_buf/validate"
)

func InitRegisterValidator() {
	//自定义手机验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		_ = v.RegisterValidation("mobile", validate.ValidateMobile)

		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {

			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details

		}, func(ut ut.Translator, fe validator.FieldError) string {

			t, _ := ut.T("mobile", fe.Field())

			return t
		})
	}
}
