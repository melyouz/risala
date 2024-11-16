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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/storage"
	"github.com/melyouz/risala/broker/internal/testing/util"
)

func setupExchangeBindingDeleteTest(t *testing.T, exchanges map[string]*internal.Exchange, exchangeName string, bindingId string) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	exchangeRepository := storage.NewInMemoryExchangeRepository(exchanges)

	path := fmt.Sprintf("%s/exchanges/%s/bindings/%s", util.ApiV1BasePath, exchangeName, bindingId)
	request := httptest.NewRequest(http.MethodDelete, path, nil)
	response := httptest.NewRecorder()

	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("exchangeName", exchangeName)
	routerCtx.URLParams.Add("bindingId", bindingId)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routerCtx))

	HandleExchangeBindingDelete(exchangeRepository)(response, request)

	return response, request
}

func TestHandleExchangeBindingDelete(t *testing.T) {

	exchanges := map[string]*internal.Exchange{
		"app.internal": util.NewTestExchangeWithBindings("app.internal", []*internal.Binding{
			{Id: uuid.New(), Queue: "tmp", RoutingKey: "#"},
		}),
		"app.external": util.NewTestExchangeWithoutBindings("app.external"),
	}

	t.Run("Deletes binding when validations pass", func(t *testing.T) {

		response, _ := setupExchangeBindingDeleteTest(t, exchanges, "app.internal", exchanges["app.internal"].Bindings[0].Id.String())

		util.AssertNoContent(t, response)
	})

	t.Run("Returns not found error when exchange does not exist", func(t *testing.T) {

		response, _ := setupExchangeBindingDeleteTest(t, exchanges, "nonExistingExchangeName", uuid.New().String())

		util.AssertNotFound(t, response, "EXCHANGE_NOT_FOUND", "Exchange 'nonExistingExchangeName' not found")
	})

	t.Run("Returns not found error when binding does not exist", func(t *testing.T) {

		nonExistingBindingId := uuid.New().String()

		response, _ := setupExchangeBindingDeleteTest(t, exchanges, "app.internal", nonExistingBindingId)

		util.AssertNotFound(t, response, "BINDING_NOT_FOUND", fmt.Sprintf("Binding '%s' not found", nonExistingBindingId))
	})

	t.Run("Returns validation error when wrong binding id format supplied", func(t *testing.T) {

		wrongBindingId := "123"

		response, _ := setupExchangeBindingDeleteTest(t, exchanges, "app.internal", wrongBindingId)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		jsonResponse := util.JSONItemResponse(response)
		assert.Equal(t, "INVALID_PARAM", jsonResponse["code"])
		assert.Equal(t, "bindingId", jsonResponse["param"])
		assert.Equal(t, "invalid UUID length: 3", jsonResponse["message"])
	})
}
