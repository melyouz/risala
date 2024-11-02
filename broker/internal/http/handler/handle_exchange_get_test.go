/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/storage"
	"github.com/melyouz/risala/broker/internal/testing/util"
)

func setupExchangeGetTest(t *testing.T, exchanges map[string]*internal.Exchange, exchangeName string) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	exchangeRepository := storage.NewInMemoryExchangeRepository(exchanges)

	path := fmt.Sprintf("%s/exchanges/%s", util.ApiV1BasePath, exchangeName)
	request := httptest.NewRequest(http.MethodGet, path, nil)
	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("exchangeName", exchangeName)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routerCtx))

	response := httptest.NewRecorder()

	HandleExchangeGet(exchangeRepository)(response, request)

	return response, request
}

func TestHandleExchangeGet(t *testing.T) {
	t.Parallel()

	exchanges := map[string]*internal.Exchange{
		"app.internal": util.NewTestExchangeWithoutBindings("app.internal"),
		"app.external": util.NewTestExchangeWithoutBindings("app.external"),
	}

	t.Run("Returns exchange when exists", func(t *testing.T) {
		t.Parallel()

		response, _ := setupExchangeGetTest(t, exchanges, "app.external")

		util.AssertOk(t, response)
		jsonResponse := util.JSONItemResponse(response)
		assert.Equal(t, "app.external", jsonResponse["name"])
		assert.Len(t, jsonResponse["bindings"], 0)
	})

	t.Run("Returns not found when exchange does not exist", func(t *testing.T) {
		t.Parallel()

		response, _ := setupExchangeGetTest(t, exchanges, "nonExistingExchangeName")

		util.AssertNotFound(t, response, "EXCHANGE_NOT_FOUND", "Exchange 'nonExistingExchangeName' not found")
	})
}
