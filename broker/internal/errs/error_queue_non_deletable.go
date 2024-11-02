/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

const QueueNonDeletableErrorCode = "QUEUE_NON_DELETABLE"

func NewCannotDeleteSystemQueueError(msg string) *Error {
	return &Error{
		Code:    QueueNonDeletableErrorCode,
		Message: msg,
	}
}
