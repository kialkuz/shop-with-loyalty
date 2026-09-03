package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	return &Validator{validate: v}
}

func (v *Validator) ValidateStructDTO(s any) error {
	return v.validate.Struct(s)
}
