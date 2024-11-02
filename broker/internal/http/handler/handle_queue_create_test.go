/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
	httputil "github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
	"github.com/melyouz/risala/broker/internal/testing/util"
)

func setupQueueCreateTest(t *testing.T, queues map[string]*internal.Queue, body map[string]interface{}) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	queueRepository := storage.NewInMemoryQueueRepository(queues)
	requestBody, _ := json.Marshal(body)
	request := httptest.NewRequest(http.MethodPost, util.ApiV1BasePath+"/queues", bytes.NewReader(requestBody))
	response := httptest.NewRecorder()

	HandleQueueCreate(queueRepository, httputil.NewJSONValidator())(response, request)

	return response, request
}

func TestHandleQueueCreate(t *testing.T) {
	t.Parallel()
	t.Run("Creates durable queue when validations pass", func(t *testing.T) {
		t.Parallel()

		queues := map[string]*internal.Queue{}
		queueBody := map[string]interface{}{
			"name":       "testDurableQueueName",
			"durability": internal.Durability.DURABLE.String(),
		}

		response, _ := setupQueueCreateTest(t, queues, queueBody)

		util.AssertCreated(t, response)
		jsonResponse := util.JSONItemResponse(response)
		assert.Equal(t, "testDurableQueueName", jsonResponse["name"])
		assert.Equal(t, internal.Durability.DURABLE.String(), jsonResponse["durability"])
	})

	t.Run("Creates transient queue when validations pass", func(t *testing.T) {
		t.Parallel()

		queues := map[string]*internal.Queue{}
		queueBody := map[string]interface{}{
			"name":       "testTransientQueueName",
			"durability": internal.Durability.TRANSIENT.String(),
		}

		response, _ := setupQueueCreateTest(t, queues, queueBody)

		util.AssertCreated(t, response)
		jsonResponse := util.JSONItemResponse(response)
		assert.Equal(t, "testTransientQueueName", jsonResponse["name"])
		assert.Equal(t, internal.Durability.TRANSIENT.String(), jsonResponse["durability"])
	})

	t.Run("Returns validation error when unknown queue durability", func(t *testing.T) {
		t.Parallel()

		queues := map[string]*internal.Queue{}
		queueBody := map[string]interface{}{
			"name":       "testInvalidDurabilityQueueName",
			"durability": "whatever",
		}

		response, _ := setupQueueCreateTest(t, queues, queueBody)

		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"durability", "Invalid value 'whatever'. Must be one of: durable transient"},
		})
	})

	t.Run("Returns validation error when no queue name supplied", func(t *testing.T) {
		t.Parallel()

		queues := map[string]*internal.Queue{}
		queueBody := map[string]interface{}{
			"durability": internal.Durability.DURABLE.String(),
		}

		response, _ := setupQueueCreateTest(t, queues, queueBody)

		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"name", "This field is required"},
		})
	})

	t.Run("Returns validation error when no queue durability supplied", func(t *testing.T) {
		t.Parallel()

		queues := map[string]*internal.Queue{}
		queueBody := map[string]interface{}{
			"name": "testInvalidDurabilityQueueName",
		}

		response, _ := setupQueueCreateTest(t, queues, queueBody)

		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"durability", "This field is required"},
		})
	})

	t.Run("Returns validation errors when no queue name nor durability supplied", func(t *testing.T) {
		t.Parallel()

		queues := map[string]*internal.Queue{}
		queueBody := map[string]interface{}{
			"shortname": "nonMappedField",
		}

		response, _ := setupQueueCreateTest(t, queues, queueBody)

		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"name", "This field is required"},
			{"durability", "This field is required"},
		})
	})

	t.Run("Returns conflict error when queue already exists", func(t *testing.T) {
		t.Parallel()

		queues := map[string]*internal.Queue{
			"events": util.NewNewQueueDurableWithoutMessages("events"),
			"tmp":    util.NewQueueTransientWithoutMessages("tmp"),
		}
		queueBody := map[string]interface{}{
			"name":       "events",
			"durability": internal.Durability.DURABLE.String(),
		}

		response, _ := setupQueueCreateTest(t, queues, queueBody)

		util.AssertConflict(t, response, "QUEUE_ALREADY_EXISTS", "Queue 'events' already exists")
	})
}
