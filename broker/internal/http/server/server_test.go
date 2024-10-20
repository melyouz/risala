/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package server

import (
	"bytes"
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
	"github.com/melyouz/risala/broker/internal/sample"
	"github.com/melyouz/risala/broker/internal/storage"
	"github.com/melyouz/risala/broker/internal/testing/util"
)

var testEventsQueue = sample.Queues["events"]
var testTmpQueue = sample.Queues["tmp"]

var testInternalExchange = sample.Exchanges["app.internal"]
var testExternalExchange = sample.Exchanges["app.external"]

type ServerSampleData struct {
	queues    map[string]internal.Queue
	exchanges map[string]internal.Exchange
}

func createTestServer(sampleData ServerSampleData) (server *Server) {
	queues := map[string]internal.Queue{}
	if sampleData.queues != nil {
		queues = sampleData.queues
	}

	exchanges := map[string]internal.Exchange{}
	if sampleData.exchanges != nil {
		exchanges = sampleData.exchanges
	}

	queueRepository := storage.NewInMemoryQueueRepository(queues)
	exchangeRepository := storage.NewInMemoryExchangeRepository(exchanges)
	server = &Server{
		router:             chi.NewRouter(),
		validate:           NewJSONValidator(),
		queueRepository:    queueRepository,
		exchangeRepository: exchangeRepository,
	}
	server.RegisterRoutes()

	return server
}

func TestHandleQueuesFind(t *testing.T) {
	t.Run("Returns list when queues exist", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		request := httptest.NewRequest(http.MethodGet, ApiV1BasePath+"/queues", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		util.AssertOk(t, response)
		jsonResponse := util.JSONCollectionResponse(response)
		assert.Len(t, jsonResponse, 2)
		assert.Equal(t, testEventsQueue.Name, jsonResponse[0]["name"])
		assert.Equal(t, testEventsQueue.Durability.String(), jsonResponse[0]["durability"])
		assert.Equal(t, testTmpQueue.Name, jsonResponse[1]["name"])
		assert.Equal(t, testTmpQueue.Durability.String(), jsonResponse[1]["durability"])
	})

	t.Run("Returns empty list when no queues", func(t *testing.T) {
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodGet, ApiV1BasePath+"/queues", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		util.AssertOk(t, response)
		jsonResponse := util.JSONCollectionResponse(response)
		assert.Len(t, jsonResponse, 0)
	})
}

func TestHandleQueueGet(t *testing.T) {
	t.Run("Returns queue when exists", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		path := fmt.Sprintf("%s/queues/%s", ApiV1BasePath, testEventsQueue.Name)
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		util.AssertOk(t, response)
		jsonResponse := util.JSONItemResponse(response)
		assert.Equal(t, testEventsQueue.Name, jsonResponse["name"])
		assert.Equal(t, testEventsQueue.Durability.String(), jsonResponse["durability"])
	})

	t.Run("Returns not found when queue does not exist", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		request := httptest.NewRequest(http.MethodGet, ApiV1BasePath+"/queues/nonExistingQueueName", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertNotFound(t, response, "QUEUE_NOT_FOUND", "Queue 'nonExistingQueueName' not found")
	})
}

func TestHandleQueueCreate(t *testing.T) {
	t.Run("Creates durable queue when validations pass", func(t *testing.T) {
		queueBody, _ := json.Marshal(map[string]interface{}{
			"name":       "testDurableQueueName",
			"durability": "durable",
		})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/queues", bytes.NewReader(queueBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertCreated(t, response)
		jsonResponse := util.JSONItemResponse(response)
		assert.Equal(t, "testDurableQueueName", jsonResponse["name"])
		assert.Equal(t, internal.Durability.DURABLE.String(), jsonResponse["durability"])
	})

	t.Run("Creates transient queue when validations pass", func(t *testing.T) {
		queueBody, _ := json.Marshal(map[string]interface{}{
			"name":       "testTransientQueueName",
			"durability": "transient",
		})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/queues", bytes.NewReader(queueBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertCreated(t, response)
		jsonResponse := util.JSONItemResponse(response)
		assert.Equal(t, "testTransientQueueName", jsonResponse["name"])
		assert.Equal(t, internal.Durability.TRANSIENT.String(), jsonResponse["durability"])
	})

	t.Run("Returns validation error when unknown queue durability", func(t *testing.T) {
		queueBody, _ := json.Marshal(map[string]interface{}{
			"name":       "testUnknownQueueName",
			"durability": "whatever",
		})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/queues", bytes.NewReader(queueBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"durability", "Invalid value 'whatever'. Must be one of: durable transient"},
		})
	})

	t.Run("Returns validation error when no queue name supplied", func(t *testing.T) {
		queueBody, _ := json.Marshal(map[string]interface{}{
			"durability": "durable",
		})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/queues", bytes.NewReader(queueBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"name", "This field is required"},
		})
	})

	t.Run("Returns validation error when no queue durability supplied", func(t *testing.T) {
		queueBody, _ := json.Marshal(map[string]interface{}{
			"name": "unknownQueueType",
		})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/queues", bytes.NewReader(queueBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"durability", "This field is required"},
		})
	})

	t.Run("Returns validation errors when no queue name nor durability supplied", func(t *testing.T) {
		queueBody, _ := json.Marshal(map[string]interface{}{
			"fullName": "nonMappedField",
		})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/queues", bytes.NewReader(queueBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"name", "This field is required"},
			{"durability", "This field is required"},
		})
	})
}

func TestHandleQueueDelete(t *testing.T) {
	t.Run("Returns accepted when queue exists", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		path := fmt.Sprintf("%s/queues/%s", ApiV1BasePath, testEventsQueue.Name)
		request := httptest.NewRequest(http.MethodDelete, path, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertAccepted(t, response)
	})

	t.Run("Returns not found when queue does not exist", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		request := httptest.NewRequest(http.MethodDelete, ApiV1BasePath+"/queues/nonExistingQueueName", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertNotFound(t, response, "QUEUE_NOT_FOUND", "Queue 'nonExistingQueueName' not found")
	})
}

func TestHandleQueueMessagePublish(t *testing.T) {
	t.Run("Publishes message when validations pass", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": "Hello world!",
		})
		server := createTestServer(ServerSampleData{
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		path := fmt.Sprintf("%s/queues/%s/publish", ApiV1BasePath, testTmpQueue.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		util.AssertOk(t, response)
		assert.Empty(t, response.Body)
	})

	t.Run("Returns validation error when no message payload supplied", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{})
		server := createTestServer(ServerSampleData{})
		path := fmt.Sprintf("%s/queues/%s/publish", ApiV1BasePath, testTmpQueue.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"payload", "This field is required"},
		})
	})

	t.Run("Returns validation error when message payload is empty", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": "",
		})

		server := createTestServer(ServerSampleData{})
		path := fmt.Sprintf("%s/queues/%s/publish", ApiV1BasePath, testTmpQueue.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"payload", "This field is required"},
		})
	})

	t.Run("Returns validation error when message payload is nil", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": nil,
		})

		server := createTestServer(ServerSampleData{})
		path := fmt.Sprintf("%s/queues/%s/publish", ApiV1BasePath, testTmpQueue.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"payload", "This field is required"},
		})
	})

	t.Run("Returns not found when queue does not exist", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": "Hello world!",
		})
		server := createTestServer(ServerSampleData{
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
			},
		})
		path := fmt.Sprintf("%s/queues/%s/publish", ApiV1BasePath, testTmpQueue.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertNotFound(t, response, "QUEUE_NOT_FOUND", "Queue 'tmp' not found")
	})
}

func TestHandleQueueMessageGet(t *testing.T) {
	t.Run("Returns empty list when no messages", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		path := fmt.Sprintf("%s/queues/%s/messages?count=1", ApiV1BasePath, testTmpQueue.Name)
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		util.AssertOk(t, response)
		jsonResponse := util.JSONCollectionResponse(response)
		assert.Len(t, jsonResponse, 0)
	})

	t.Run("Returns messages when exist", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})

		path := fmt.Sprintf("%s/queues/%s/publish", ApiV1BasePath, testTmpQueue.Name)
		for i := 1; i <= 3; i++ {
			messageBody, _ := json.Marshal(map[string]interface{}{
				"payload": fmt.Sprintf("Message %d", i),
			})
			request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)
			util.AssertOk(t, response)
			assert.Empty(t, response.Body)
		}

		path = fmt.Sprintf("%s/queues/%s/messages", ApiV1BasePath, testTmpQueue.Name)
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		util.AssertOk(t, response)
		var jsonResponse1 []map[string]interface{}
		_ = json.Unmarshal([]byte(response.Body.String()), &jsonResponse1)
		assert.Len(t, jsonResponse1, 1)
		assert.NotEmpty(t, jsonResponse1[0]["id"])
		assert.Equal(t, "Message 3", jsonResponse1[0]["payload"])

		path = fmt.Sprintf("%s/queues/%s/messages?count=2", ApiV1BasePath, testTmpQueue.Name)
		request = httptest.NewRequest(http.MethodGet, path, nil)
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)
		util.AssertOk(t, response)
		var jsonResponse2 []map[string]interface{}
		_ = json.Unmarshal([]byte(response.Body.String()), &jsonResponse2)
		assert.Len(t, jsonResponse2, 2)
		assert.NotEmpty(t, jsonResponse2[0]["id"])
		assert.Equal(t, "Message 3", jsonResponse2[0]["payload"])
		assert.NotEmpty(t, jsonResponse2[0]["id"])
		assert.Equal(t, "Message 2", jsonResponse2[1]["payload"])

		path = fmt.Sprintf("%s/queues/%s/messages?count=200", ApiV1BasePath, testTmpQueue.Name)
		request = httptest.NewRequest(http.MethodGet, path, nil)
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)
		util.AssertOk(t, response)
		var jsonResponse3 []map[string]interface{}
		_ = json.Unmarshal([]byte(response.Body.String()), &jsonResponse3)
		assert.Len(t, jsonResponse3, 3)
		assert.NotEmpty(t, jsonResponse3[0]["id"])
		assert.Equal(t, "Message 3", jsonResponse3[0]["payload"])
		assert.NotEmpty(t, jsonResponse3[0]["id"])
		assert.Equal(t, "Message 2", jsonResponse3[1]["payload"])
		assert.NotEmpty(t, jsonResponse3[0]["id"])
		assert.Equal(t, "Message 1", jsonResponse3[2]["payload"])
	})

	t.Run("Returns not found when queue does not exist", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": "Hello world!",
		})
		server := createTestServer(ServerSampleData{
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
			},
		})
		path := fmt.Sprintf("%s/queues/%s/messages", ApiV1BasePath, testTmpQueue.Name)
		request := httptest.NewRequest(http.MethodGet, path, bytes.NewReader(messageBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertNotFound(t, response, "QUEUE_NOT_FOUND", "Queue 'tmp' not found")
	})
}

func TestHandleExchangesFind(t *testing.T) {
	t.Run("Returns list when exchanges exist", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
		})
		request := httptest.NewRequest(http.MethodGet, ApiV1BasePath+"/exchanges", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		util.AssertOk(t, response)
		jsonResponse := util.JSONCollectionResponse(response)
		assert.Len(t, jsonResponse, 2)
		assert.Equal(t, testExternalExchange.Name, jsonResponse[0]["name"])
		assert.Len(t, jsonResponse[0]["bindings"], 0)
		assert.Equal(t, testInternalExchange.Name, jsonResponse[1]["name"])
		assert.Len(t, jsonResponse[1]["bindings"], 0)
	})

	t.Run("Returns empty list when no exchanges", func(t *testing.T) {
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodGet, ApiV1BasePath+"/exchanges", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		util.AssertOk(t, response)
		jsonResponse := util.JSONCollectionResponse(response)
		assert.Len(t, jsonResponse, 0)
	})
}

func TestHandleExchangeGet(t *testing.T) {
	t.Run("Returns exchange when exists", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
		})
		path := fmt.Sprintf("%s/exchanges/%s", ApiV1BasePath, testExternalExchange.Name)
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		util.AssertOk(t, response)
		jsonResponse := util.JSONItemResponse(response)
		assert.Equal(t, testExternalExchange.Name, jsonResponse["name"])
		assert.Len(t, jsonResponse["bindings"], 0)
	})

	t.Run("Returns not found when exchange does not exist", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
		})
		request := httptest.NewRequest(http.MethodGet, ApiV1BasePath+"/exchanges/nonExistingExchangeName", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertNotFound(t, response, "EXCHANGE_NOT_FOUND", "Exchange 'nonExistingExchangeName' not found")
	})
}

func TestHandleExchangeCreate(t *testing.T) {
	t.Run("Creates exchange when validations pass", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
		})
		exchangeBody, _ := json.Marshal(map[string]interface{}{
			"name": "app.tmp",
		})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/exchanges", bytes.NewReader(exchangeBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertCreated(t, response)
		jsonResponse := util.JSONItemResponse(response)
		assert.Equal(t, "app.tmp", jsonResponse["name"])
		assert.Len(t, jsonResponse["bindings"], 0)
	})

	t.Run("Returns validation error when no exchange name supplied", func(t *testing.T) {
		exchangeBody, _ := json.Marshal(map[string]interface{}{})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/exchanges", bytes.NewReader(exchangeBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"name", "This field is required"},
		})
	})

	t.Run("Returns validation error when exchange name is empty", func(t *testing.T) {
		exchangeBody, _ := json.Marshal(map[string]interface{}{
			"name": "",
		})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/exchanges", bytes.NewReader(exchangeBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"name", "This field is required"},
		})
	})

	t.Run("Returns validation error when exchange name is nil", func(t *testing.T) {
		exchangeBody, _ := json.Marshal(map[string]interface{}{
			"name": nil,
		})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/exchanges", bytes.NewReader(exchangeBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"name", "This field is required"},
		})
	})

	t.Run("Returns conflict error when exchange already exists", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
		})
		exchangeBody, _ := json.Marshal(map[string]interface{}{
			"name": testInternalExchange.Name,
		})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/exchanges", bytes.NewReader(exchangeBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertConflict(t, response, "EXCHANGE_ALREADY_EXISTS", "Exchange 'app.internal' already exists")
	})
}

func TestHandleExchangeDelete(t *testing.T) {
	t.Run("Returns accepted when exchange exists", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
		})
		path := fmt.Sprintf("%s/exchanges/%s", ApiV1BasePath, testExternalExchange.Name)
		request := httptest.NewRequest(http.MethodDelete, path, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertAccepted(t, response)
	})

	t.Run("Returns not found when exchange does not exist", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
		})
		request := httptest.NewRequest(http.MethodDelete, ApiV1BasePath+"/exchanges/nonExistingExchangeName", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertNotFound(t, response, "EXCHANGE_NOT_FOUND", "Exchange 'nonExistingExchangeName' not found")
	})
}

func TestHandleExchangeBindingAdd(t *testing.T) {
	t.Run("Adds binding when validations pass", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		bindingBody, _ := json.Marshal(map[string]interface{}{
			"queue":      testEventsQueue.Name,
			"routingKey": "#",
		})
		path := fmt.Sprintf("%s/exchanges/%s/bindings", ApiV1BasePath, testInternalExchange.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bindingBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertCreated(t, response)
		jsonResponse := util.JSONItemResponse(response)
		assert.NotEmpty(t, jsonResponse["id"])
		assert.Equal(t, testEventsQueue.Name, jsonResponse["queue"])
		assert.Equal(t, "#", jsonResponse["routingKey"])
	})

	t.Run("Returns not found error when exchange does not exist", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		bindingBody, _ := json.Marshal(map[string]interface{}{
			"queue":      testEventsQueue.Name,
			"routingKey": "#",
		})
		path := fmt.Sprintf("%s/exchanges/nonExistingExchangeName/bindings", ApiV1BasePath)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bindingBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertNotFound(t, response, "EXCHANGE_NOT_FOUND", "Exchange 'nonExistingExchangeName' not found")
	})

	t.Run("Returns not found error when queue does not exist", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		bindingBody, _ := json.Marshal(map[string]interface{}{
			"queue":      "nonExistingQueueName",
			"routingKey": "#",
		})
		path := fmt.Sprintf("%s/exchanges/%s/bindings", ApiV1BasePath, testInternalExchange.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bindingBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertNotFound(t, response, "QUEUE_NOT_FOUND", "Queue 'nonExistingQueueName' not found")
	})

	t.Run("Returns conflict error when binding already exists", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		_ = server.exchangeRepository.AddBinding(testInternalExchange.Name, internal.Binding{Id: uuid.New(), Queue: testEventsQueue.Name})

		bindingBody, _ := json.Marshal(map[string]interface{}{
			"queue":      testEventsQueue.Name,
			"routingKey": "#",
		})
		path := fmt.Sprintf("%s/exchanges/%s/bindings", ApiV1BasePath, testInternalExchange.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bindingBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertConflict(t, response, "BINDING_ALREADY_EXISTS", "Binding to Queue 'events' already exists")
	})

	t.Run("Returns validation error when no queue name supplied", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		bindingBody, _ := json.Marshal(map[string]interface{}{
			"routingKey": "#",
		})
		path := fmt.Sprintf("%s/exchanges/%s/bindings", ApiV1BasePath, testInternalExchange.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bindingBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"queue", "This field is required"},
		})
	})
}

func TestHandleExchangeBindingDelete(t *testing.T) {
	t.Run("Deletes binding when validations pass", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		binding := internal.Binding{Id: uuid.New(), Queue: testEventsQueue.Name}
		_ = server.exchangeRepository.AddBinding(testInternalExchange.Name, binding)

		path := fmt.Sprintf("%s/exchanges/%s/bindings/%s", ApiV1BasePath, testInternalExchange.Name, binding.Id)
		request := httptest.NewRequest(http.MethodDelete, path, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertAccepted(t, response)
	})

	t.Run("Returns not found error when exchange does not exist", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})

		path := fmt.Sprintf("%s/exchanges/nonExistingExchangeName/bindings/%s", ApiV1BasePath, uuid.New().String())
		request := httptest.NewRequest(http.MethodDelete, path, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertNotFound(t, response, "EXCHANGE_NOT_FOUND", "Exchange 'nonExistingExchangeName' not found")
	})

	t.Run("Returns not found error when binding does not exist", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})

		nonExistingBindingId := uuid.New().String()
		path := fmt.Sprintf("%s/exchanges/%s/bindings/%s", ApiV1BasePath, testInternalExchange.Name, nonExistingBindingId)
		request := httptest.NewRequest(http.MethodDelete, path, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertNotFound(t, response, "BINDING_NOT_FOUND", fmt.Sprintf("Binding '%s' not found", nonExistingBindingId))
	})

	t.Run("Returns validation error when wrong binding id format supplied", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})

		wrongBindingId := "123"
		path := fmt.Sprintf("%s/exchanges/%s/bindings/%s", ApiV1BasePath, testInternalExchange.Name, wrongBindingId)
		request := httptest.NewRequest(http.MethodDelete, path, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		jsonResponse := util.JSONItemResponse(response)
		assert.Equal(t, "INVALID_PARAM", jsonResponse["code"])
		assert.Equal(t, "bindingId", jsonResponse["param"])
		assert.Equal(t, "invalid UUID length: 3", jsonResponse["message"])
	})
}

func TestHandleExchangeMessagePublish(t *testing.T) {
	t.Run("Publishes message when validations pass & queue binding exists", func(t *testing.T) {
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
				testExternalExchange.Name: testExternalExchange,
			},
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		bindingBody, _ := json.Marshal(map[string]interface{}{
			"queue":      testTmpQueue.Name,
			"routingKey": "#",
		})
		path := fmt.Sprintf("%s/exchanges/%s/bindings", ApiV1BasePath, testInternalExchange.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bindingBody))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": "Hello world!",
		})
		path = fmt.Sprintf("%s/exchanges/%s/publish", ApiV1BasePath, testInternalExchange.Name)
		request = httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)
		util.AssertOk(t, response)
		assert.Empty(t, response.Body)
	})

	t.Run("Returns validation error when no message payload supplied", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{})
		server := createTestServer(ServerSampleData{})
		path := fmt.Sprintf("%s/exchanges/%s/publish", ApiV1BasePath, testTmpQueue.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"payload", "This field is required"},
		})
	})

	t.Run("Returns validation error when message payload is empty", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": "",
		})

		server := createTestServer(ServerSampleData{})
		path := fmt.Sprintf("%s/exchanges/%s/publish", ApiV1BasePath, testTmpQueue.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"payload", "This field is required"},
		})
	})

	t.Run("Returns validation error when message payload is nil", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": nil,
		})

		server := createTestServer(ServerSampleData{})
		path := fmt.Sprintf("%s/exchanges/%s/publish", ApiV1BasePath, testTmpQueue.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertValidationErrors(t, response, []errs.ValidationError{
			{"payload", "This field is required"},
		})
	})

	t.Run("Returns not found when exchange does not exist", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{
			"payload": "Hello world!",
		})
		server := createTestServer(ServerSampleData{
			exchanges: map[string]internal.Exchange{
				testInternalExchange.Name: testInternalExchange,
			},
			queues: map[string]internal.Queue{
				testEventsQueue.Name: testEventsQueue,
				testTmpQueue.Name:    testTmpQueue,
			},
		})
		path := fmt.Sprintf("%s/exchanges/%s/publish", ApiV1BasePath, testExternalExchange.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		util.AssertNotFound(t, response, "EXCHANGE_NOT_FOUND", "Exchange 'app.external' not found")
	})

	// TODO: Check message is routed to queue
}
