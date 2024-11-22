/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/melyouz/risala/broker/internal/http/handler"
)

const ApiV1BasePath = "/api/v1"
const apiV1DocsBasePath = "/api/v1/docs"
const ApiV1OpenApiSpecJsonFilePath = "docs/api/openapi3_0.json"

func (s *Server) RegisterRoutes() {
	// queues
	queuesRouter := chi.NewRouter()
	queuesRouter.Post("/", handler.HandleQueueCreate(s.queueRepository, s.validate))
	queuesRouter.Get("/", handler.HandleQueueFind(s.queueRepository))
	queuesRouter.Get("/{queueName}", handler.HandleQueueGet(s.queueRepository))
	queuesRouter.Delete("/{queueName}", handler.HandleQueueDelete(s.queueRepository))
	queuesRouter.Post("/{queueName}/messages/publish", handler.HandleQueueMessagePublish(s.queueRepository, s.validate))
	queuesRouter.Get("/{queueName}/messages/peek", handler.HandleQueueMessagePeek(s.queueRepository))
	queuesRouter.Post("/{queueName}/messages/consume", handler.HandleQueueMessageConsume(s.queueRepository))
	queuesRouter.Post("/{queueName}/messages/purge", handler.HandleQueueMessagePurge(s.queueRepository))
	queuesRouter.Post("/{queueName}/messages/get", handler.HandleQueueMessageGet(s.queueRepository))
	queuesRouter.Post("/{queueName}/messages/{messageId}/ack", handler.HandleQueueMessageAck(s.queueRepository))
	queuesRouter.Post("/{queueName}/messages/{messageId}/nack", handler.HandleQueueMessageNack(s.queueRepository))

	// exchanges
	exchangesRouter := chi.NewRouter()
	exchangesRouter.Post("/", handler.HandleExchangeCreate(s.exchangeRepository, s.validate))
	exchangesRouter.Get("/", handler.HandleExchangeFind(s.exchangeRepository))
	exchangesRouter.Get("/{exchangeName}", handler.HandleExchangeGet(s.exchangeRepository))
	exchangesRouter.Delete("/{exchangeName}", handler.HandleExchangeDelete(s.exchangeRepository))
	exchangesRouter.Post("/{exchangeName}/bindings", handler.HandleExchangeBindingAdd(s.exchangeRepository, s.queueRepository, s.validate))
	exchangesRouter.Delete("/{exchangeName}/bindings/{bindingId}", handler.HandleExchangeBindingDelete(s.exchangeRepository))
	exchangesRouter.Post("/{exchangeName}/messages/publish", handler.HandleExchangeMessagePublish(s.exchangeRepository, s.queueRepository, s.validate))

	// v1 routes group
	s.router.Route(ApiV1BasePath, func(r chi.Router) {
		r.Mount("/queues", queuesRouter)
		r.Mount("/exchanges", exchangesRouter)
	})

	// docs
	s.router.Get(fmt.Sprintf("/%s", ApiV1OpenApiSpecJsonFilePath), func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, ApiV1OpenApiSpecJsonFilePath)
	})
	s.router.Get("/api*", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fmt.Sprintf("%s/index.html", apiV1DocsBasePath), http.StatusPermanentRedirect)
	})
	s.router.Get(fmt.Sprintf("%s/*", apiV1DocsBasePath), httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:8000/%s", ApiV1OpenApiSpecJsonFilePath)),
		httpSwagger.AfterScript(`document.querySelectorAll(".topbar")[0].remove();`),
	))
}
