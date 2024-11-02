/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
	httputil "github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
	"github.com/melyouz/risala/broker/internal/testing/util"
)

func setupExchangeCreateTest(t *testing.T, exchanges map[string]*internal.Exchange, body map[string]interface{}) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	exchangeRepository := storage.NewInMemoryExchangeRepository(exchanges)
	exchangeBody, _ := json.Marshal(body)
	request := httptest.NewRequest(http.MethodPost, util.ApiV1BasePath+"/exchanges", bytes.NewReader(exchangeBody))
	response := httptest.NewRecorder()

	HandleExchangeCreate(exchangeRepository, httputil.NewJSONValidator())(response, request)

	return response, request
}

func TestHandleExchangeCreate(t *testing.T) {
	t.Run("Creates exchange when validations pass", func(t *testing.T) {
		exchanges := map[string]*internal.Exchange{
			"app.internal": util.NewTestExchangeWithoutBindings("app.internal"),
			"app.external": util.NewTestExchangeWithoutBindings("app.external"),
		}
		exchangeBody := map[string]interface{}{
			"name": "app.tmp",
		}

		response, _ := setupExchangeCreateTest(t, exchanges, exchangeBody)

		util.AssertCreated(t, response)
		jsonResponse := util.JSONItemResponse(response)
		assert.Equal(t, "app.tmp", jsonResponse["name"])
		assert.Len(t, jsonResponse["bindings"], 0)
	})

	t.Run("Returns validation error when no exchange name supplied", func(t *testing.T) {

		exchanges := map[string]*internal.Exchange{}
		exchangeBody := map[string]interface{}{}

		response, _ := setupExchangeCreateTest(t, exchanges, exchangeBody)

		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"name", "This field is required"},
		})
	})

	t.Run("Returns validation error when exchange name is empty", func(t *testing.T) {

		exchanges := map[string]*internal.Exchange{}
		exchangeBody := map[string]interface{}{
			"name": "",
		}

		response, _ := setupExchangeCreateTest(t, exchanges, exchangeBody)

		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"name", "This field is required"},
		})
	})

	t.Run("Returns validation error when exchange name is nil", func(t *testing.T) {

		exchanges := map[string]*internal.Exchange{}
		exchangeBody := map[string]interface{}{
			"name": nil,
		}

		response, _ := setupExchangeCreateTest(t, exchanges, exchangeBody)

		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"name", "This field is required"},
		})
	})

	t.Run("Returns conflict error when exchange already exists", func(t *testing.T) {

		exchanges := map[string]*internal.Exchange{
			"app.internal": util.NewTestExchangeWithoutBindings("app.internal"),
			"app.external": util.NewTestExchangeWithoutBindings("app.external"),
		}
		exchangeBody := map[string]interface{}{
			"name": "app.internal",
		}

		response, _ := setupExchangeCreateTest(t, exchanges, exchangeBody)

		util.AssertConflict(t, response, "EXCHANGE_ALREADY_EXISTS", "Exchange 'app.internal' already exists")
	})
}
