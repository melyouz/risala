/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package util

import (
	"fmt"
	"github.com/melyouz/risala/broker/internal/errs"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHandleQueueFind(t *testing.T) {
	t.Run("Returns mapped status code", func(t *testing.T) {
		msg := "Whatever error message..."
		err := errs.Error{Code: errs.QueueNotFoundErrorCode, Message: msg}
		statusCode := HttpStatusCodeFromAppError(&err)
		assert.Equal(t, http.StatusNotFound, statusCode)
		assert.Equal(t, msg, err.GetMessage())
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
