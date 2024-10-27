/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

const MessageNotFoundErrorCode = "MESSAGE_NOT_FOUND"

func NewMessageNotFoundError(msg string) *Error {
	return &Error{
		Code:    MessageNotFoundErrorCode,
		Message: msg,
	}
}
