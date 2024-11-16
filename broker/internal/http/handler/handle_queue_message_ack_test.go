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
	"github.com/melyouz/risala/broker/internal/errs"
	"github.com/melyouz/risala/broker/internal/storage"
	"github.com/melyouz/risala/broker/internal/testing/util"
)

func setupQueueMessageAckTest(t *testing.T, queues map[string]*internal.Queue, queueName string, messageId uuid.UUID) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	queueRepository := storage.NewInMemoryQueueRepository(queues)

	path := fmt.Sprintf("%s/queues/%s/messages/%s/ack", util.ApiV1BasePath, queueName, messageId.String())
	request := httptest.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()

	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("queueName", queueName)
	routerCtx.URLParams.Add("messageId", messageId.String())
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routerCtx))

	HandleQueueMessageAck(queueRepository)(response, request)

	return response, request
}

func TestHandleQueueMessageAck(t *testing.T) {

	queues := map[string]*internal.Queue{
		"events":                     util.NewTestQueueDurableWithoutMessages("events"),
		"tmp":                        util.NewTestQueueTransientWithoutMessages("tmp"),
		internal.DeadLetterQueueName: util.NewTestSystemQueueWithoutMessages(internal.DeadLetterQueueName),
	}

	t.Run("Acknowledges a message", func(t *testing.T) {
		messages := []*internal.Message{
			{Id: uuid.New(), Payload: "Message 1", Processing: true},
			{Id: uuid.New(), Payload: "Message 2"},
			{Id: uuid.New(), Payload: "Message 3"},
			{Id: uuid.New(), Payload: "Message 4"},
			{Id: uuid.New(), Payload: "Message 5"},
		}
		queues["events"].Messages = messages
		initialMessageCount := len(messages)
		messageId := messages[0].Id

		response, _ := setupQueueMessageAckTest(t, queues, "events", messageId)

		util.AssertNoContent(t, response)
		assert.Len(t, queues["events"].Messages, initialMessageCount-1)
		assert.Len(t, queues[internal.DeadLetterQueueName].Messages, 0)
	})

	t.Run("Returns not found when message is not being processed", func(t *testing.T) {
		messages := []*internal.Message{
			{Id: uuid.New(), Payload: "Message 1"},
			{Id: uuid.New(), Payload: "Message 2", Processing: true},
		}
		queues["events"].Messages = messages
		messageId := messages[0].Id

		response, _ := setupQueueMessageAckTest(t, queues, "events", messageId)

		util.AssertNotFound(t, response, errs.MessageNotFoundErrorCode, fmt.Sprintf("Message '%s' not found", messageId.String()))
	})

	t.Run("Returns not found when message does not exist", func(t *testing.T) {
		messageId := uuid.New()

		response, _ := setupQueueMessageAckTest(t, queues, "events", messageId)

		util.AssertNotFound(t, response, errs.MessageNotFoundErrorCode, fmt.Sprintf("Message '%s' not found", messageId.String()))
	})

	t.Run("Returns not found when queue does not exist", func(t *testing.T) {
		messageId := uuid.New()

		response, _ := setupQueueMessageAckTest(t, queues, "nonExistingQueueName", messageId)

		util.AssertNotFound(t, response, "QUEUE_NOT_FOUND", "Queue 'nonExistingQueueName' not found")
	})
}
