package server

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type validateErr struct {
	Errors map[string]string `json:"errors"`
}

func Validate(req interface{}) (*validateErr, error) {
	val := validator.New()

	if err := val.Struct(req); err != nil {
		var ive *validator.InvalidValidationError
		if errors.As(err, &ive) {
			return nil, ive
		}

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			valErrs := map[string]string{}
			for _, err := range ve {
				valErrs[err.Field()] = err.Error()
			}

			return &validateErr{Errors: valErrs}, nil
		}
	}

	return nil, nil
}
