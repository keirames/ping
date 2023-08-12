package validator

import "github.com/go-playground/validator/v10"

var V *validator.Validate

func New() {
	V = validator.New()
}
