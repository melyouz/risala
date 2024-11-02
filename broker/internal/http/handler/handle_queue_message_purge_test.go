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

func setupQueueMessagePurgeTest(t *testing.T, queues map[string]*internal.Queue, queueName string) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	queueRepository := storage.NewInMemoryQueueRepository(queues)

	path := fmt.Sprintf("%s/queues/%s/messages/purge", util.ApiV1BasePath, queueName)
	request := httptest.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()

	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("queueName", queueName)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routerCtx))

	HandleQueueMessagePurge(queueRepository)(response, request)

	return response, request
}

func TestHandleQueueMessagePurge(t *testing.T) {

	queues := map[string]*internal.Queue{
		"events": util.NewTestQueueDurableWithoutMessages("events"),
		"tmp":    util.NewTestQueueTransientWithoutMessages("tmp"),
	}

	t.Run("Returns accepted on success", func(t *testing.T) {
		messages := []*internal.Message{
			{Id: uuid.New(), Payload: "Message 1"},
			{Id: uuid.New(), Payload: "Message 2"},
			{Id: uuid.New(), Payload: "Message 3"},
			{Id: uuid.New(), Payload: "Message 4"},
			{Id: uuid.New(), Payload: "Message 5"},
		}
		queues["events"].Messages = messages

		response, _ := setupQueueMessagePurgeTest(t, queues, "events")

		util.AssertAccepted(t, response)
		assert.Empty(t, response.Body)
		assert.Empty(t, queues["events"].Messages)
	})

	t.Run("Returns not found when queue does not exist", func(t *testing.T) {

		response, _ := setupQueueMessagePurgeTest(t, queues, "nonExistingQueueName")

		util.AssertNotFound(t, response, "QUEUE_NOT_FOUND", "Queue 'nonExistingQueueName' not found")
	})
}
