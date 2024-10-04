/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package util

import (
	"net/http"

	"github.com/melyouz/risala/broker/internal/errs"
)

var httpDefaultStatusCode = http.StatusInternalServerError
var httpStatusCodes = map[string]int{
	errs.ExchangeNotFoundErrorCode: http.StatusNotFound,
	errs.ExchangeExistsErrorCode:   http.StatusConflict,
	errs.QueueNotFoundErrorCode:    http.StatusNotFound,
	errs.BindingNotFoundErrorCode:  http.StatusNotFound,
	errs.BindingExistsErrorCode:    http.StatusConflict,
	errs.ParamInvalidErrorCode:     http.StatusBadRequest,
}

func HttpStatusCodeFromAppError(err errs.AppError) int {
	statusCode, ok := httpStatusCodes[err.GetCode()]
	if ok {
		return statusCode
	}

	return httpDefaultStatusCode
}
