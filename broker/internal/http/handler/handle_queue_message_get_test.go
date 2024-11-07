/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
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
	"github.com/melyouz/risala/broker/internal/storage"
	"github.com/melyouz/risala/broker/internal/testing/util"
)

func setupQueueMessageGetTest(t *testing.T, queues map[string]*internal.Queue, queueName string) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	queueRepository := storage.NewInMemoryQueueRepository(queues)

	path := fmt.Sprintf("%s/queues/%s/messages/get", util.ApiV1BasePath, queueName)
	request := httptest.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()

	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("queueName", queueName)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routerCtx))

	HandleQueueMessageGet(queueRepository)(response, request)

	return response, request
}

func TestHandleQueueMessageGet(t *testing.T) {

	queues := map[string]*internal.Queue{
		"events": util.NewTestQueueDurableWithoutMessages("events"),
		"tmp":    util.NewTestQueueTransientWithoutMessages("tmp"),
	}

	t.Run("Marks as processing & returns first message when available", func(t *testing.T) {
		messages := []*internal.Message{
			{Id: uuid.New(), Payload: "Message 1"},
			{Id: uuid.New(), Payload: "Message 2"},
			{Id: uuid.New(), Payload: "Message 3"},
			{Id: uuid.New(), Payload: "Message 4"},
			{Id: uuid.New(), Payload: "Message 5"},
		}
		queues["events"].Messages = messages
		initialMessageCount := len(messages)
		firstMessage := messages[0]

		response, _ := setupQueueMessageGetTest(t, queues, "events")

		assert.True(t, firstMessage.IsProcessing())
		util.AssertOk(t, response)
		var jsonResponse map[string]interface{}
		_ = json.Unmarshal(response.Body.Bytes(), &jsonResponse)
		assert.Equal(t, firstMessage.Id.String(), jsonResponse["id"])
		assert.Equal(t, firstMessage.Payload, jsonResponse["payload"])
		assert.Len(t, queues["events"].Messages, initialMessageCount)
	})

	t.Run("Returns no content when no messages", func(t *testing.T) {
		response, _ := setupQueueMessageGetTest(t, queues, "tmp")

		util.AssertNoContent(t, response)
	})

	t.Run("Returns not found when queue does not exist", func(t *testing.T) {

		response, _ := setupQueueMessageGetTest(t, queues, "nonExistingQueueName")

		util.AssertNotFound(t, response, "QUEUE_NOT_FOUND", "Queue 'nonExistingQueueName' not found")
	})
}
