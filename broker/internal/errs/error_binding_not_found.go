/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

const BindingNotFoundErrorCode = "BINDING_NOT_FOUND"

func NewBindingNotFoundError(msg string) *Error {
	return &Error{
		Code:    BindingNotFoundErrorCode,
		Message: msg,
	}
}
