package errs

const QueueNotFoundErrorCode = "QUEUE_NOT_FOUND"

func NewQueueNotFoundError(msg string) *Error {
	return &Error{
		Code:    QueueNotFoundErrorCode,
		Message: msg,
	}
}
