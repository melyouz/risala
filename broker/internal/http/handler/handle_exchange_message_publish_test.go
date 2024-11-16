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

func setupExchangeMessagePublishTest(t *testing.T, queues map[string]*internal.Queue, exchanges map[string]*internal.Exchange, exchangeName string, messageBody []byte) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	queueRepository := storage.NewInMemoryQueueRepository(queues)
	exchangeRepository := storage.NewInMemoryExchangeRepository(exchanges)

	path := fmt.Sprintf("%s/exchanges/%s/messages/publish", util.ApiV1BasePath, exchangeName)
	request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
	response := httptest.NewRecorder()

	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("exchangeName", exchangeName)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routerCtx))

	HandleExchangeMessagePublish(exchangeRepository, queueRepository, httputil.NewJSONValidator())(response, request)

	return response, request
}

func TestHandleExchangeMessagePublish(t *testing.T) {

	exchanges := map[string]*internal.Exchange{
		"app.internal": util.NewTestExchangeWithBindings("app.internal", []*internal.Binding{
			{Id: uuid.New(), Queue: "tmp", RoutingKey: "#"},
		}),
		"app.external": util.NewTestExchangeWithoutBindings("app.external"),
	}
	queues := map[string]*internal.Queue{
		"events": util.NewTestQueueDurableWithoutMessages("events"),
		"tmp":    util.NewTestQueueTransientWithoutMessages("tmp"),
	}

	t.Run("Publishes message when validations pass & queue binding exists", func(t *testing.T) {

		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": "Hello world from Exchange",
		})

		tmpQueueMessagesCount := len(queues["tmp"].Messages)
		response, _ := setupExchangeMessagePublishTest(t, queues, exchanges, "app.internal", messageBody)

		util.AssertCreated(t, response)
		var jsonResponse map[string]interface{}
		_ = json.Unmarshal(response.Body.Bytes(), &jsonResponse)
		assert.NotEmpty(t, jsonResponse["id"])
		assert.Equal(t, "Hello world from Exchange", jsonResponse["payload"])
		assert.Len(t, queues["tmp"].Messages, tmpQueueMessagesCount+1)
	})

	t.Run("Returns validation error when no message payload supplied", func(t *testing.T) {

		messageBody, _ := json.Marshal(map[string]interface{}{})

		response, _ := setupExchangeMessagePublishTest(t, queues, exchanges, "app.internal", messageBody)

		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"payload", "This field is required"},
		})
	})

	t.Run("Returns validation error when message payload is empty", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": "",
		})

		response, _ := setupExchangeMessagePublishTest(t, queues, exchanges, "app.internal", messageBody)

		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"payload", "This field is required"},
		})
	})

	t.Run("Returns validation error when message payload is nil", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": nil,
		})

		response, _ := setupExchangeMessagePublishTest(t, queues, exchanges, "app.internal", messageBody)

		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"payload", "This field is required"},
		})
	})

	t.Run("Returns not found when exchange does not exist", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": "Hello world!",
		})

		response, _ := setupExchangeMessagePublishTest(t, queues, exchanges, "nonExistingExchangeName", messageBody)

		util.AssertNotFound(t, response, "EXCHANGE_NOT_FOUND", "Exchange 'nonExistingExchangeName' not found")
	})
}
