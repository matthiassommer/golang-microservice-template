package api

import (
	"gopkg.in/go-playground/validator.v9"
)

// CustomValidator extends the base validator.
type CustomValidator struct {
	Validator *validator.Validate
}

// NewValidator returns a new instance of the validator.
func NewValidator() *CustomValidator {
	v9Validator := validator.New()

	return &CustomValidator{Validator: v9Validator}
}

// Validate validates structs.
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
