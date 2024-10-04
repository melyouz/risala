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
	queuesRouter.Get("/", handler.HandleQueueFind(s.queueRepository))
	queuesRouter.Get("/{queueName}", handler.HandleQueueGet(s.queueRepository))
	queuesRouter.Post("/", handler.HandleQueueCreate(s.queueRepository, s.validate))
	queuesRouter.Delete("/{queueName}", handler.HandleQueueDelete(s.queueRepository))
	queuesRouter.Post("/{queueName}/messages", handler.HandleQueueMessagePublish(s.queueRepository, s.validate))

	// exchanges
	exchangesRouter := chi.NewRouter()
	exchangesRouter.Get("/", handler.HandleExchangeFind(s.exchangeRepository))
	exchangesRouter.Get("/{exchangeName}", handler.HandleExchangeGet(s.exchangeRepository))
	exchangesRouter.Post("/", handler.HandleExchangeCreate(s.exchangeRepository, s.validate))
	exchangesRouter.Delete("/{exchangeName}", handler.HandleExchangeDelete(s.exchangeRepository))
	exchangesRouter.Post("/{exchangeName}/bindings", handler.HandleExchangeBindingAdd(s.exchangeRepository, s.queueRepository, s.validate))
	exchangesRouter.Delete("/{exchangeName}/bindings/{bindingId}", handler.HandleExchangeBindingDelete(s.exchangeRepository))

	// main router
	s.router.Route(ApiV1BasePath, func(r chi.Router) {
		r.Mount("/queues", queuesRouter)
		r.Mount("/exchanges", exchangesRouter)
	})
}
