/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

const decodeErrorCode = "DECODE_ERROR"

func NewDecodeError(msg string) *Error {
	return &Error{
		Code:    decodeErrorCode,
		Message: msg,
	}
}
