package errs

const ExchangeExistsErrorCode = "EXCHANGE_ALREADY_EXISTS"

func NewExchangeExistsError(msg string) *Error {
	return &Error{
		Code:    ExchangeExistsErrorCode,
		Message: msg,
	}
}
