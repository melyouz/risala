/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/storage"
	"github.com/melyouz/risala/broker/internal/testing/util"
)

func setupQueueFindTest(t *testing.T, queues map[string]*internal.Queue) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	exchangeRepository := storage.NewInMemoryQueueRepository(queues)
	request := httptest.NewRequest(http.MethodGet, util.ApiV1BasePath+"/exchanges", nil)
	response := httptest.NewRecorder()

	HandleQueueFind(exchangeRepository)(response, request)

	return response, request
}

func TestHandleQueueFind(t *testing.T) {
	t.Run("Returns list when queues exist", func(t *testing.T) {

		queues := map[string]*internal.Queue{
			"events": util.NewTestQueueDurableWithoutMessages("events"),
			"tmp":    util.NewTestQueueTransientWithoutMessages("tmp"),
		}

		response, _ := setupQueueFindTest(t, queues)

		util.AssertOk(t, response)
		jsonResponse := util.JSONCollectionResponse(response)
		assert.Len(t, jsonResponse, 2)
		assert.Equal(t, "events", jsonResponse[0]["name"])
		assert.Equal(t, internal.Durability.DURABLE.String(), jsonResponse[0]["durability"])
		assert.Equal(t, "tmp", jsonResponse[1]["name"])
		assert.Equal(t, internal.Durability.TRANSIENT.String(), jsonResponse[1]["durability"])
	})

	t.Run("Returns empty list when no queues", func(t *testing.T) {

		queues := map[string]*internal.Queue{}

		response, _ := setupQueueFindTest(t, queues)

		util.AssertOk(t, response)
		jsonResponse := util.JSONCollectionResponse(response)
		assert.Empty(t, jsonResponse)
	})
}
