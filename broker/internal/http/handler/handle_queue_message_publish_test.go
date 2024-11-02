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
	"github.com/stretchr/testify/assert"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
	httputil "github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
	"github.com/melyouz/risala/broker/internal/testing/util"
)

func setupQueueMessagePublishTest(t *testing.T, queues map[string]*internal.Queue, queueName string, messageBody []byte) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	queueRepository := storage.NewInMemoryQueueRepository(queues)

	path := fmt.Sprintf("%s/queues/%s/messages/publish", util.ApiV1BasePath, queueName)
	request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
	response := httptest.NewRecorder()

	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("queueName", queueName)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routerCtx))

	HandleQueueMessagePublish(queueRepository, httputil.NewJSONValidator())(response, request)

	return response, request
}

func TestHandleQueueMessagePublish(t *testing.T) {
	t.Parallel()

	queues := map[string]*internal.Queue{
		"events": util.NewNewQueueDurableWithoutMessages("events"),
		"tmp":    util.NewQueueTransientWithoutMessages("tmp"),
	}

	t.Run("Publishes message when validations pass", func(t *testing.T) {
		t.Parallel()
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": "Hello world!",
		})

		response, _ := setupQueueMessagePublishTest(t, queues, "tmp", messageBody)

		util.AssertOk(t, response)
		assert.Empty(t, response.Body)
	})

	t.Run("Returns validation error when no message payload supplied", func(t *testing.T) {
		t.Parallel()
		messageBody, _ := json.Marshal(map[string]interface{}{})

		response, _ := setupQueueMessagePublishTest(t, queues, "tmp", messageBody)

		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"payload", "This field is required"},
		})
	})

	t.Run("Returns validation error when message payload is empty", func(t *testing.T) {
		t.Parallel()
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": "",
		})

		response, _ := setupQueueMessagePublishTest(t, queues, "tmp", messageBody)

		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"payload", "This field is required"},
		})
	})

	t.Run("Returns validation error when message payload is nil", func(t *testing.T) {
		t.Parallel()
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": nil,
		})

		response, _ := setupQueueMessagePublishTest(t, queues, "tmp", messageBody)

		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"payload", "This field is required"},
		})
	})

	t.Run("Returns not found when queue does not exist", func(t *testing.T) {
		t.Parallel()
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": "Hello world!",
		})

		response, _ := setupQueueMessagePublishTest(t, queues, "nonExistingQueueName", messageBody)

		util.AssertNotFound(t, response, "QUEUE_NOT_FOUND", "Queue 'nonExistingQueueName' not found")
	})
}
