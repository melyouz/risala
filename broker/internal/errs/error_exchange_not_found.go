/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

const ExchangeNotFoundErrorCode = "EXCHANGE_NOT_FOUND"

func NewExchangeNotFoundError(msg string) *Error {
	return &Error{
		Code:    ExchangeNotFoundErrorCode,
		Message: msg,
	}
}
