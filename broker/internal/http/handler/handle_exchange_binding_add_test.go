/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
	httputil "github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
	"github.com/melyouz/risala/broker/internal/testing/util"
)

func setupExchangeBindingAddTest(t *testing.T, queues map[string]*internal.Queue, exchanges map[string]*internal.Exchange, exchangeName string, bindingBody []byte) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	queueRepository := storage.NewInMemoryQueueRepository(queues)
	exchangeRepository := storage.NewInMemoryExchangeRepository(exchanges)

	path := fmt.Sprintf("%s/exchanges/%s/bindings", util.ApiV1BasePath, exchangeName)
	request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bindingBody))
	response := httptest.NewRecorder()

	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("exchangeName", exchangeName)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routerCtx))

	HandleExchangeBindingAdd(exchangeRepository, queueRepository, httputil.NewJSONValidator())(response, request)

	return response, request
}

func TestHandleExchangeBindingAdd(t *testing.T) {
	t.Parallel()

	exchanges := map[string]*internal.Exchange{
		"app.internal": util.NewTestExchangeWithBindings("app.internal", []*internal.Binding{
			{Id: uuid.New(), Queue: "tmp", RoutingKey: "#"},
		}),
		"app.external": util.NewTestExchangeWithoutBindings("app.external"),
	}
	queues := map[string]*internal.Queue{
		"events": util.NewNewQueueDurableWithoutMessages("events"),
		"tmp":    util.NewQueueTransientWithoutMessages("tmp"),
	}

	t.Run("Adds binding when validations pass", func(t *testing.T) {
		t.Parallel()
		bindingBody, _ := json.Marshal(map[string]interface{}{
			"queue":      "events",
			"routingKey": "#",
		})

		response, _ := setupExchangeBindingAddTest(t, queues, exchanges, "app.internal", bindingBody)

		util.AssertCreated(t, response)
		jsonResponse := util.JSONItemResponse(response)
		assert.NotEmpty(t, jsonResponse["id"])
		assert.Equal(t, "events", jsonResponse["queue"])
		assert.Equal(t, "#", jsonResponse["routingKey"])
	})

	t.Run("Returns not found error when exchange does not exist", func(t *testing.T) {
		t.Parallel()
		bindingBody, _ := json.Marshal(map[string]interface{}{
			"queue":      "events",
			"routingKey": "#",
		})

		response, _ := setupExchangeBindingAddTest(t, queues, exchanges, "nonExistingExchangeName", bindingBody)

		util.AssertNotFound(t, response, "EXCHANGE_NOT_FOUND", "Exchange 'nonExistingExchangeName' not found")
	})

	t.Run("Returns not found error when queue does not exist", func(t *testing.T) {
		t.Parallel()

		bindingBody, _ := json.Marshal(map[string]interface{}{
			"queue":      "nonExistingQueueName",
			"routingKey": "#",
		})

		response, _ := setupExchangeBindingAddTest(t, queues, exchanges, "app.internal", bindingBody)

		util.AssertNotFound(t, response, "QUEUE_NOT_FOUND", "Queue 'nonExistingQueueName' not found")
	})

	t.Run("Returns conflict error when binding already exists", func(t *testing.T) {
		t.Parallel()

		bindingBody, _ := json.Marshal(map[string]interface{}{
			"queue":      "tmp",
			"routingKey": "#",
		})

		response, _ := setupExchangeBindingAddTest(t, queues, exchanges, "app.internal", bindingBody)

		util.AssertConflict(t, response, "BINDING_ALREADY_EXISTS", "Binding to Queue 'tmp' already exists")
	})

	t.Run("Returns validation error when no queue name supplied", func(t *testing.T) {
		t.Parallel()

		bindingBody, _ := json.Marshal(map[string]interface{}{
			"routingKey": "#",
		})

		response, _ := setupExchangeBindingAddTest(t, queues, exchanges, "app.internal", bindingBody)

		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"queue", "This field is required"},
		})
	})
}
