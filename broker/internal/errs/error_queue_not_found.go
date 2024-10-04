/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

const QueueNotFoundErrorCode = "QUEUE_NOT_FOUND"

func NewQueueNotFoundError(msg string) *Error {
	return &Error{
		Code:    QueueNotFoundErrorCode,
		Message: msg,
	}
}
