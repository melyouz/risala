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

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/storage"
	"github.com/melyouz/risala/broker/internal/testing/util"
)

func setupExchangeDeleteTest(t *testing.T, exchanges map[string]*internal.Exchange, exchangeName string) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	exchangeRepository := storage.NewInMemoryExchangeRepository(exchanges)
	path := fmt.Sprintf("%s/exchanges/%s", util.ApiV1BasePath, exchangeName)
	request := httptest.NewRequest(http.MethodDelete, path, nil)
	response := httptest.NewRecorder()

	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("exchangeName", exchangeName)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routerCtx))

	HandleExchangeDelete(exchangeRepository)(response, request)

	return response, request
}

func TestHandleExchangeDelete(t *testing.T) {
	exchanges := map[string]*internal.Exchange{
		"app.internal": util.NewTestExchangeWithoutBindings("app.internal"),
		"app.external": util.NewTestExchangeWithoutBindings("app.external"),
	}

	t.Run("Returns accepted when exchange exists", func(t *testing.T) {
		response, _ := setupExchangeDeleteTest(t, exchanges, "app.external")

		util.AssertAccepted(t, response)
	})

	t.Run("Returns not found when exchange does not exist", func(t *testing.T) {
		response, _ := setupExchangeDeleteTest(t, exchanges, "nonExistingExchangeName")

		util.AssertNotFound(t, response, "EXCHANGE_NOT_FOUND", "Exchange 'nonExistingExchangeName' not found")
	})
}
