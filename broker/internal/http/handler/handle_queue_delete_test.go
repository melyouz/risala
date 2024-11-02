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

func setupQueueDeleteTest(t *testing.T, queues map[string]*internal.Queue, exchangeName string) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	queueRepository := storage.NewInMemoryQueueRepository(queues)
	path := fmt.Sprintf("%s/exchanges/%s", util.ApiV1BasePath, exchangeName)
	request := httptest.NewRequest(http.MethodDelete, path, nil)
	response := httptest.NewRecorder()

	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("queueName", exchangeName)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routerCtx))

	HandleQueueDelete(queueRepository)(response, request)

	return response, request
}

func TestHandleQueueDelete(t *testing.T) {
	t.Parallel()

	queues := map[string]*internal.Queue{
		"events": util.NewNewQueueDurableWithoutMessages("events"),
		"tmp":    util.NewQueueTransientWithoutMessages("tmp"),
	}

	t.Run("Returns accepted when queue exists", func(t *testing.T) {
		t.Parallel()

		response, _ := setupQueueDeleteTest(t, queues, "events")

		util.AssertAccepted(t, response)
	})

	t.Run("Returns not found when queue does not exist", func(t *testing.T) {
		t.Parallel()

		response, _ := setupQueueDeleteTest(t, queues, "nonExistingQueueName")

		util.AssertNotFound(t, response, "QUEUE_NOT_FOUND", "Queue 'nonExistingQueueName' not found")
	})
}
