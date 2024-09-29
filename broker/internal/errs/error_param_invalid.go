package errs

const ParamInvalidErrorCode = "PARAM_INVALID"

func NewParamInvalidError(msg string) *Error {
	return &Error{
		Code:    ParamInvalidErrorCode,
		Message: msg,
	}
}
