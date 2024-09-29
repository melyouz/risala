package errs

const BindingExistsErrorCode = "BINDING_ALREADY_EXISTS"

func NewBindingExistsError(msg string) *Error {
	return &Error{
		Code:    BindingExistsErrorCode,
		Message: msg,
	}
}
