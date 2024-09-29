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
	queuesRouter.Get("/{name}", handler.HandleQueueGet(s.queueRepository))
	queuesRouter.Post("/", handler.HandleQueueCreate(s.queueRepository, s.validate))
	queuesRouter.Delete("/{name}", handler.HandleQueueDelete(s.queueRepository))

	// exchanges
	exchangesRouter := chi.NewRouter()
	exchangesRouter.Get("/", handler.HandleExchangeFind(s.exchangeRepository))
	exchangesRouter.Get("/{name}", handler.HandleExchangeGet(s.exchangeRepository))
	exchangesRouter.Post("/", handler.HandleExchangeCreate(s.exchangeRepository, s.validate))
	exchangesRouter.Delete("/{name}", handler.HandleExchangeDelete(s.exchangeRepository))
	exchangesRouter.Post("/{name}/bindings", handler.HandleExchangeBindingAdd(s.exchangeRepository, s.queueRepository, s.validate))
	exchangesRouter.Delete("/{name}/bindings/{bindingId}", handler.HandleExchangeBindingDelete(s.exchangeRepository))

	// main router
	s.router.Route(ApiV1BasePath, func(r chi.Router) {
		r.Mount("/queues", queuesRouter)
		r.Mount("/exchanges", exchangesRouter)
	})
}
