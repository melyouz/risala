/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

const connectionErrorCode = "CONNECTION_ERROR"

func NewConnectionError(msg string) *Error {
	return &Error{
		Code:    connectionErrorCode,
		Message: msg,
	}
}
