package ginmiddleware

import (
	"github.com/go-playground/validator/v10"
)

func RequiredIf(fl validator.FieldLevel) bool {
	sessionType := fl.Parent().FieldByName("SessionType").Int()
	switch sessionType {
	default:
		return true
	}
	return true
}
