/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

const ValidationErrorCode = "VALIDATION_ERROR"

func NewValidationError(errors []validator.FieldError) *Error {
	validationErrors := make([]ValidationError, len(errors))
	for i, fe := range errors {
		validationErrors[i] = ValidationError{fe.Field(), validationMessageForTag(fe)}
	}

	return &Error{
		Code:   ValidationErrorCode,
		Errors: validationErrors,
	}
}

func validationMessageForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	//case "email":
	//	return "Invalid email"
	case "oneof":
		return fmt.Sprintf("Invalid value '%s'. Must be one of: %s", fe.Value(), fe.Param())
	default:
		return fe.Error()
	}
}
