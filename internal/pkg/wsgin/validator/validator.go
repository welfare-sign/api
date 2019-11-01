package validator

import (
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

func isMobile(level validator.FieldLevel) bool {
	reg := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(level.Field().String())
}
