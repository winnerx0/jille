package utils

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type XValidator struct {
	Validator *validator.Validate
}

var Validate *validator.Validate

func init() {
	Validate = validator.New()
	Validate.RegisterValidation("optionlistmin", func(fl validator.FieldLevel) bool {

		param, err := strconv.Atoi(fl.Param())

		if err != nil {
			return false
		}
		return fl.Field().Len() >= param
	})

	Validate.RegisterValidation("optionlistmax", func(fl validator.FieldLevel) bool {

		param, err := strconv.Atoi(fl.Param())

		if err != nil {
			return false
		}
		return fl.Field().Len() <= param
	})
}

func (v XValidator) Validate(data any) error {
	var validationError error

	err := Validate.Struct(data)
	if err != nil {
		err := err.(validator.ValidationErrors)
		switch err[0].Tag() {
		case "required":
			return fmt.Errorf("%s required", err[0].StructField())
		case "email":
			return errors.New("Invalid email address")
		case "min":
			return fmt.Errorf("%s must be %s characters", err[0].StructField(), err[0].Param())

		case "max":
			return fmt.Errorf("%s must be at most %s characters", err[0].StructField(), err[0].Param())

		case "optionlistmin":
			return fmt.Errorf("Options must be at least %s", err[0].Param())

		case "optionlistmax":
			return fmt.Errorf("Options must be at most %s", err[0].Param())
		default:
			return errors.New(err[0].Error())
		}
	}

	return validationError
}
