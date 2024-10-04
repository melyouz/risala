/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

const ExchangeExistsErrorCode = "EXCHANGE_ALREADY_EXISTS"

func NewExchangeExistsError(msg string) *Error {
	return &Error{
		Code:    ExchangeExistsErrorCode,
		Message: msg,
	}
}
