/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

const ReadErrorCode = "READ_ERROR"

func NewReadError(msg string) *Error {
	return &Error{
		Code:    ReadErrorCode,
		Message: msg,
	}
}
