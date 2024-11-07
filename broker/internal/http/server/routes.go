/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package server

import (
	"github.com/go-chi/chi/v5"

	"github.com/melyouz/risala/broker/internal/http/handler"
)

const ApiV1BasePath = "/api/v1"

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

	// main router
	s.router.Route(ApiV1BasePath, func(r chi.Router) {
		r.Mount("/queues", queuesRouter)
		r.Mount("/exchanges", exchangesRouter)
	})
}
