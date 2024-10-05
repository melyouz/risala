/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

const ParamInvalidErrorCode = "INVALID_PARAM"

func NewParamInvalidError(param string, msg string) *Error {
	return &Error{
		Code:    ParamInvalidErrorCode,
		Param:   param,
		Message: msg,
	}
}
