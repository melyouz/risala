/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package util

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/melyouz/risala/broker/internal/errs"
)

func TestHttpStatusCodeFromAppError(t *testing.T) {
	t.Run("Returns mapped status code", func(t *testing.T) {
		const testErrorCode = "TEST_ERROR_CODE"
		httpStatusCodes[testErrorCode] = http.StatusGone

		msg := "Whatever error message..."
		err := errs.Error{Code: testErrorCode, Message: msg}
		statusCode := HttpStatusCodeFromAppError(&err)
		assert.Equal(t, testErrorCode, err.GetCode())
		assert.Equal(t, msg, err.GetMessage())
		assert.Equal(t, http.StatusGone, statusCode)
		assert.Equal(t, fmt.Sprintf("%s: %s", err.GetCode(), err.GetMessage()), err.Error())
	})
	t.Run("Returns default status code", func(t *testing.T) {
		msg := "Whatever error message..."
		err := errs.Error{Code: "NON_MAPPED_ERROR_CODE", Message: msg}
		statusCode := HttpStatusCodeFromAppError(&err)
		assert.Equal(t, httpDefaultStatusCode, statusCode)
		assert.Equal(t, msg, err.GetMessage())
		assert.Equal(t, fmt.Sprintf("%s: %s", err.GetCode(), err.GetMessage()), err.Error())
	})
}
