/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

const QueueExistsErrorCode = "QUEUE_ALREADY_EXISTS"

func NewQueueExistsError(msg string) *Error {
	return &Error{
		Code:    QueueExistsErrorCode,
		Message: msg,
	}
}
