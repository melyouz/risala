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

func setupQueueMessageConsumeTest(t *testing.T, queues map[string]*internal.Queue, queueName string, limit string) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	queueRepository := storage.NewInMemoryQueueRepository(queues)

	path := fmt.Sprintf("%s/queues/%s/messages/consume", util.ApiV1BasePath, queueName)
	if limit != "" {
		path += fmt.Sprintf("?limit=%s", limit)
	}
	request := httptest.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()

	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("queueName", queueName)
	routerCtx.URLParams.Add("limit", limit)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routerCtx))

	HandleQueueMessageConsume(queueRepository)(response, request)

	return response, request
}

func TestHandleQueueMessageConsume(t *testing.T) {

	queues := map[string]*internal.Queue{
		"events": util.NewNewQueueDurableWithoutMessages("events"),
		"tmp":    util.NewQueueTransientWithoutMessages("tmp"),
	}

	t.Run("Returns empty list when no messages", func(t *testing.T) {
		response, _ := setupQueueMessageConsumeTest(t, queues, "events", "10")

		util.AssertOk(t, response)
		jsonResponse := util.JSONCollectionResponse(response)
		assert.Len(t, jsonResponse, 0)
	})

	t.Run("Returns one message when no limit supplied", func(t *testing.T) {
		messages := []*internal.Message{
			{Id: uuid.New(), Payload: "Message 1"},
			{Id: uuid.New(), Payload: "Message 2"},
			{Id: uuid.New(), Payload: "Message 3"},
			{Id: uuid.New(), Payload: "Message 4"},
			{Id: uuid.New(), Payload: "Message 5"},
		}
		queues["events"].Messages = messages
		initialMessageCount := len(messages)

		response, _ := setupQueueMessageConsumeTest(t, queues, "events", "1")

		util.AssertOk(t, response)
		var jsonResponse []map[string]interface{}
		_ = json.Unmarshal(response.Body.Bytes(), &jsonResponse)
		assert.Len(t, jsonResponse, 1)
		assert.NotEmpty(t, jsonResponse[0]["id"])
		assert.Equal(t, "Message 1", jsonResponse[0]["payload"])
		assert.Len(t, queues["events"].Messages, initialMessageCount-1)
	})

	t.Run("Returns N messages when limit=N and available messages > N", func(t *testing.T) {
		messages := []*internal.Message{
			{Id: uuid.New(), Payload: "Message 1"},
			{Id: uuid.New(), Payload: "Message 2"},
			{Id: uuid.New(), Payload: "Message 3"},
			{Id: uuid.New(), Payload: "Message 4"},
			{Id: uuid.New(), Payload: "Message 5"},
		}
		queues["events"].Messages = messages
		initialMessageCount := len(messages)

		response, _ := setupQueueMessageConsumeTest(t, queues, "events", "2")

		util.AssertOk(t, response)
		var jsonResponse []map[string]interface{}
		_ = json.Unmarshal(response.Body.Bytes(), &jsonResponse)
		assert.Len(t, jsonResponse, 2)
		assert.Equal(t, "Message 1", jsonResponse[0]["payload"])
		assert.Equal(t, "Message 2", jsonResponse[1]["payload"])
		assert.Len(t, queues["events"].Messages, initialMessageCount-2)
	})

	t.Run("Returns all messages when limit=N and available messages < N", func(t *testing.T) {
		messages := []*internal.Message{
			{Id: uuid.New(), Payload: "Message 1"},
			{Id: uuid.New(), Payload: "Message 2"},
			{Id: uuid.New(), Payload: "Message 3"},
			{Id: uuid.New(), Payload: "Message 4"},
			{Id: uuid.New(), Payload: "Message 5"},
		}
		queues["events"].Messages = messages
		initialMessageCount := len(messages)

		response, _ := setupQueueMessageConsumeTest(t, queues, "events", "200")

		util.AssertOk(t, response)
		var jsonResponse []map[string]interface{}
		_ = json.Unmarshal(response.Body.Bytes(), &jsonResponse)
		assert.Len(t, jsonResponse, initialMessageCount)
		for i, message := range jsonResponse {
			assert.NotEmpty(t, message["id"])
			assert.Equal(t, fmt.Sprintf("Message %d", i+1), message["payload"])
		}
		assert.Len(t, queues["events"].Messages, 0)
	})

	t.Run("Returns not found when queue does not exist", func(t *testing.T) {

		response, _ := setupQueueMessageConsumeTest(t, queues, "nonExistingQueueName", "200")

		util.AssertNotFound(t, response, "QUEUE_NOT_FOUND", "Queue 'nonExistingQueueName' not found")
	})
}
