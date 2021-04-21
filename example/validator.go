package main

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	vtzh "github.com/go-playground/validator/v10/translations/zh"
)

type User struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Age       uint8  `validate:"gte=0,lte=130"`
	Email     string `validate:"required,email"`
}

func main() {
	user := &User{
		FirstName: "firstName",
		LastName:  "lastName",
		Age:       136,
		Email:     "f163.com",
	}
	// 创建一个验证器
	validate := validator.New()

	// 创建消息国际化通用翻译器
	zhTranslator := zh.New()
	uniTranslator := ut.New(zhTranslator, zhTranslator)
	translator, found := uniTranslator.GetTranslator("zh")
	if found {
		// 将翻译器和验证器绑定
		err := vtzh.RegisterDefaultTranslations(validate, translator)
		if err != nil {
			fmt.Println("++", err)
		}
	}

	err := validate.Struct(user)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			fmt.Println("--", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range errs {
				fmt.Println("**", err.Translate(translator))
			}
		}
	}
}
