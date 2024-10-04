/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package errs

import "fmt"

type AppError interface {
	error
	GetCode() string
	GetMessage() string
}

type Error struct {
	Code    string            `json:"code,omitempty"`
	Message string            `json:"message,omitempty"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) GetCode() string {
	return e.Code
}

func (e *Error) GetMessage() string {
	return e.Message
}
