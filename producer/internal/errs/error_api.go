/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

const apiErrorCode = "API_ERROR"

func NewApiError(msg string) *Error {
	return &Error{
		Code:    apiErrorCode,
		Message: msg,
	}
}
