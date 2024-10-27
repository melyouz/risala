/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

const QueueEmptyErrorCode = "QUEUE_EMPTY"

func NewQueueEmptyError(msg string) *Error {
	return &Error{
		Code:    QueueEmptyErrorCode,
		Message: msg,
	}
}
