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
	"github.com/melyouz/risala/broker/internal/sample"
	"github.com/melyouz/risala/broker/internal/storage"
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

		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, "[{\"name\":\"events\",\"durability\":\"durable\"},{\"name\":\"tmp\",\"durability\":\"transient\"}]", response.Body.String())
	})

	t.Run("Returns empty list when no queues", func(t *testing.T) {
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodGet, ApiV1BasePath+"/queues", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, "[]", response.Body.String())
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

		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, "{\"name\":\"events\",\"durability\":\"durable\"}", response.Body.String())
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

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.JSONEq(t, "{\"code\":\"QUEUE_NOT_FOUND\",\"message\":\"Queue 'nonExistingQueueName' not found\"}", response.Body.String())
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

		assert.Equal(t, http.StatusCreated, response.Code)
		assert.JSONEq(t, "{\"name\":\"testDurableQueueName\",\"durability\":\"durable\"}", response.Body.String())
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

		assert.Equal(t, http.StatusCreated, response.Code)
		assert.JSONEq(t, "{\"name\":\"testTransientQueueName\",\"durability\":\"transient\"}", response.Body.String())
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

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, "{\"code\":\"VALIDATION_ERROR\",\"errors\":[{\"field\":\"durability\",\"message\":\"Invalid value 'whatever'. Must be one of: durable transient\"}]}", response.Body.String())
	})

	t.Run("Returns validation error when no queue name supplied", func(t *testing.T) {
		queueBody, _ := json.Marshal(map[string]interface{}{
			"durability": "durable",
		})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/queues", bytes.NewReader(queueBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, "{\"code\":\"VALIDATION_ERROR\",\"errors\":[{\"field\":\"name\",\"message\":\"This field is required\"}]}", response.Body.String())
	})

	t.Run("Returns validation error when no queue durability supplied", func(t *testing.T) {
		queueBody, _ := json.Marshal(map[string]interface{}{
			"name": "unknownQueueType",
		})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/queues", bytes.NewReader(queueBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, "{\"code\":\"VALIDATION_ERROR\",\"errors\":[{\"field\":\"durability\",\"message\":\"This field is required\"}]}", response.Body.String())
	})

	t.Run("Returns validation errors when no queue name nor durability supplied", func(t *testing.T) {
		queueBody, _ := json.Marshal(map[string]interface{}{
			"fullName": "nonMappedField",
		})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/queues", bytes.NewReader(queueBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, "{\"code\":\"VALIDATION_ERROR\",\"errors\":[{\"field\":\"name\",\"message\":\"This field is required\"},{\"field\":\"durability\",\"message\":\"This field is required\"}]}", response.Body.String())
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

		assert.Equal(t, http.StatusAccepted, response.Code)
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

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.JSONEq(t, "{\"code\":\"QUEUE_NOT_FOUND\",\"message\":\"Queue 'nonExistingQueueName' not found\"}", response.Body.String())
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

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Empty(t, response.Body)
	})

	t.Run("Returns validation error when no message payload supplied", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{})
		server := createTestServer(ServerSampleData{})
		path := fmt.Sprintf("%s/queues/%s/publish", ApiV1BasePath, testTmpQueue.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, "{\"code\":\"VALIDATION_ERROR\",\"errors\":[{\"field\":\"payload\",\"message\":\"This field is required\"}]}", response.Body.String())
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

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, "{\"code\":\"VALIDATION_ERROR\",\"errors\":[{\"field\":\"payload\",\"message\":\"This field is required\"}]}", response.Body.String())
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

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, "{\"code\":\"VALIDATION_ERROR\",\"errors\":[{\"field\":\"payload\",\"message\":\"This field is required\"}]}", response.Body.String())
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

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.JSONEq(t, "{\"code\":\"QUEUE_NOT_FOUND\",\"message\":\"Queue 'tmp' not found\"}", response.Body.String())
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

		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, "[]", response.Body.String())
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
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Empty(t, response.Body)
		}

		path = fmt.Sprintf("%s/queues/%s/messages", ApiV1BasePath, testTmpQueue.Name)
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, "[{\"payload\":\"Message 3\"}]", response.Body.String())

		path = fmt.Sprintf("%s/queues/%s/messages?count=2", ApiV1BasePath, testTmpQueue.Name)
		request = httptest.NewRequest(http.MethodGet, path, nil)
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, "[{\"payload\":\"Message 3\"},{\"payload\":\"Message 2\"}]", response.Body.String())

		path = fmt.Sprintf("%s/queues/%s/messages?count=200", ApiV1BasePath, testTmpQueue.Name)
		request = httptest.NewRequest(http.MethodGet, path, nil)
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, "[{\"payload\":\"Message 3\"},{\"payload\":\"Message 2\"},{\"payload\":\"Message 1\"}]", response.Body.String())
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

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.JSONEq(t, "{\"code\":\"QUEUE_NOT_FOUND\",\"message\":\"Queue 'tmp' not found\"}", response.Body.String())
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

		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, "[{\"name\":\"app.external\",\"bindings\":[]},{\"name\":\"app.internal\",\"bindings\":[]}]", response.Body.String())
	})

	t.Run("Returns empty list when no exchanges", func(t *testing.T) {
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodGet, ApiV1BasePath+"/exchanges", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, "[]", response.Body.String())
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

		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, "{\"name\":\"app.external\",\"bindings\":[]}", response.Body.String())
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

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.JSONEq(t, "{\"code\":\"EXCHANGE_NOT_FOUND\",\"message\":\"Exchange 'nonExistingExchangeName' not found\"}", response.Body.String())
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

		assert.Equal(t, http.StatusCreated, response.Code)
		assert.JSONEq(t, "{\"name\":\"app.tmp\",\"bindings\":[]}", response.Body.String())
	})

	t.Run("Returns validation error when no exchange name supplied", func(t *testing.T) {
		exchangeBody, _ := json.Marshal(map[string]interface{}{})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/exchanges", bytes.NewReader(exchangeBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, "{\"code\":\"VALIDATION_ERROR\",\"errors\":[{\"field\":\"name\",\"message\":\"This field is required\"}]}", response.Body.String())
	})

	t.Run("Returns validation error when exchange name is empty", func(t *testing.T) {
		exchangeBody, _ := json.Marshal(map[string]interface{}{
			"name": "",
		})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/exchanges", bytes.NewReader(exchangeBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, "{\"code\":\"VALIDATION_ERROR\",\"errors\":[{\"field\":\"name\",\"message\":\"This field is required\"}]}", response.Body.String())
	})

	t.Run("Returns validation error when exchange name is nil", func(t *testing.T) {
		exchangeBody, _ := json.Marshal(map[string]interface{}{
			"name": nil,
		})
		server := createTestServer(ServerSampleData{})
		request := httptest.NewRequest(http.MethodPost, ApiV1BasePath+"/exchanges", bytes.NewReader(exchangeBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, "{\"code\":\"VALIDATION_ERROR\",\"errors\":[{\"field\":\"name\",\"message\":\"This field is required\"}]}", response.Body.String())
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

		assert.Equal(t, http.StatusConflict, response.Code)
		assert.JSONEq(t, "{\"code\":\"EXCHANGE_ALREADY_EXISTS\",\"message\":\"Exchange 'app.internal' already exists\"}", response.Body.String())
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

		assert.Equal(t, http.StatusAccepted, response.Code)
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

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.JSONEq(t, "{\"code\":\"EXCHANGE_NOT_FOUND\",\"message\":\"Exchange 'nonExistingExchangeName' not found\"}", response.Body.String())
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

		assert.Equal(t, http.StatusCreated, response.Code)

		respJson := response.Body.String()
		respBinding := &internal.Binding{}
		_ = json.Unmarshal([]byte(respJson), respBinding)
		assert.Equal(t, testEventsQueue.Name, respBinding.Queue)
		assert.Equal(t, "#", respBinding.RoutingKey)
		assert.NotEmpty(t, respBinding.Id)

		expectedJson := fmt.Sprintf("{\"id\":\"%s\",\"queue\":\"events\",\"routingKey\":\"#\"}", respBinding.Id)
		assert.JSONEq(t, expectedJson, respJson)
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

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.JSONEq(t, "{\"code\":\"EXCHANGE_NOT_FOUND\",\"message\":\"Exchange 'nonExistingExchangeName' not found\"}", response.Body.String())
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

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.JSONEq(t, "{\"code\":\"QUEUE_NOT_FOUND\",\"message\":\"Queue 'nonExistingQueueName' not found\"}", response.Body.String())
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
		exchange, _ := server.exchangeRepository.GetExchange(testInternalExchange.Name)
		exchange.AddBinding(internal.Binding{Queue: testEventsQueue.Name})
		server.exchangeRepository.StoreExchange(exchange)

		bindingBody, _ := json.Marshal(map[string]interface{}{
			"queue":      testEventsQueue.Name,
			"routingKey": "#",
		})
		path := fmt.Sprintf("%s/exchanges/%s/bindings", ApiV1BasePath, testInternalExchange.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(bindingBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assert.Equal(t, http.StatusConflict, response.Code)
		assert.JSONEq(t, "{\"code\":\"BINDING_ALREADY_EXISTS\",\"message\":\"Binding to Queue 'events' already exists\"}", response.Body.String())
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
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, "{\"code\":\"VALIDATION_ERROR\",\"errors\":[{\"field\":\"queue\",\"message\":\"This field is required\"}]}", response.Body.String())
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
		exchange, _ := server.exchangeRepository.GetExchange(testInternalExchange.Name)
		binding := internal.Binding{Id: uuid.New(), Queue: testEventsQueue.Name}
		exchange.AddBinding(binding)
		server.exchangeRepository.StoreExchange(exchange)

		path := fmt.Sprintf("%s/exchanges/%s/bindings/%s", ApiV1BasePath, testInternalExchange.Name, binding.Id)
		request := httptest.NewRequest(http.MethodDelete, path, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusAccepted, response.Code)
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

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.JSONEq(t, "{\"code\":\"EXCHANGE_NOT_FOUND\",\"message\":\"Exchange 'nonExistingExchangeName' not found\"}", response.Body.String())
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

		assert.Equal(t, http.StatusNotFound, response.Code)
		expectedJson := fmt.Sprintf("{\"code\":\"BINDING_NOT_FOUND\",\"message\":\"Binding '%s' not found\"}", nonExistingBindingId)
		assert.JSONEq(t, expectedJson, response.Body.String())
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
		assert.JSONEq(t, "{\"code\":\"INVALID_PARAM\",\"param\":\"bindingId\",\"message\":\"invalid UUID length: 3\"}", response.Body.String())
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
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Empty(t, response.Body)
	})

	t.Run("Returns validation error when no message payload supplied", func(t *testing.T) {
		messageBody, _ := json.Marshal(map[string]interface{}{})
		server := createTestServer(ServerSampleData{})
		path := fmt.Sprintf("%s/exchanges/%s/publish", ApiV1BasePath, testTmpQueue.Name)
		request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(messageBody))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, "{\"code\":\"VALIDATION_ERROR\",\"errors\":[{\"field\":\"payload\",\"message\":\"This field is required\"}]}", response.Body.String())
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

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, "{\"code\":\"VALIDATION_ERROR\",\"errors\":[{\"field\":\"payload\",\"message\":\"This field is required\"}]}", response.Body.String())
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

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, "{\"code\":\"VALIDATION_ERROR\",\"errors\":[{\"field\":\"payload\",\"message\":\"This field is required\"}]}", response.Body.String())
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

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.JSONEq(t, "{\"code\":\"EXCHANGE_NOT_FOUND\",\"message\":\"Exchange 'app.external' not found\"}", response.Body.String())
	})

	// TODO: Check message is routed to queue
}
