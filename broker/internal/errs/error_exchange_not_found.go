package errs

const ExchangeNotFoundErrorCode = "EXCHANGE_NOT_FOUND"

func NewExchangeNotFoundError(msg string) *Error {
	return &Error{
		Code:    ExchangeNotFoundErrorCode,
		Message: msg,
	}
}
