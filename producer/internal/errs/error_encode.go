/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

const encodeErrorCode = "ENCODE_ERROR"

func NewEncodeError(msg string) *Error {
	return &Error{
		Code:    encodeErrorCode,
		Message: msg,
	}
}
