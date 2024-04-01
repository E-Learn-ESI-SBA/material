package utils

import "github.com/go-playground/validator/v10"

func Validator(object interface{}) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(object)

}
