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

func setupQueueGetTest(t *testing.T, queues map[string]*internal.Queue, queueName string) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	queueRepository := storage.NewInMemoryQueueRepository(queues)

	path := fmt.Sprintf("%s/queues/%s", util.ApiV1BasePath, queueName)
	request := httptest.NewRequest(http.MethodGet, path, nil)
	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("queueName", queueName)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routerCtx))

	response := httptest.NewRecorder()

	HandleQueueGet(queueRepository)(response, request)

	return response, request
}

func TestHandleQueueGet(t *testing.T) {
	t.Parallel()

	queues := map[string]*internal.Queue{
		"events": util.NewNewQueueDurableWithoutMessages("events"),
		"tmp":    util.NewQueueTransientWithoutMessages("tmp"),
	}

	t.Run("Returns queue when exists", func(t *testing.T) {
		t.Parallel()

		response, _ := setupQueueGetTest(t, queues, "events")

		util.AssertOk(t, response)
		jsonResponse := util.JSONItemResponse(response)
		assert.Equal(t, "events", jsonResponse["name"])
		assert.Equal(t, internal.Durability.DURABLE.String(), jsonResponse["durability"])
	})

	t.Run("Returns not found when queue does not exist", func(t *testing.T) {
		t.Parallel()

		response, _ := setupQueueGetTest(t, queues, "nonExistingQueueName")

		util.AssertNotFound(t, response, "QUEUE_NOT_FOUND", "Queue 'nonExistingQueueName' not found")
	})
}
