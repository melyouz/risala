/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
)

type Server struct {
	listenAddr         string
	router             *chi.Mux
	validate           *validator.Validate
	queueRepository    storage.QueueRepository
	exchangeRepository storage.ExchangeRepository
}

func NewServer(
	listenAddr string,
	router *chi.Mux,
	queuesRepository storage.QueueRepository,
	exchangesRepository storage.ExchangeRepository,
) *http.Server {
	s := &Server{
		listenAddr:         listenAddr,
		router:             router,
		validate:           util.NewJSONValidator(),
		queueRepository:    queuesRepository,
		exchangeRepository: exchangesRepository,
	}
	s.RegisterRoutes()

	server := &http.Server{
		Addr:         s.listenAddr,
		Handler:      s.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
