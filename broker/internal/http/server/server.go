package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/melyouz/risala/broker/internal/storage"
	"net/http"
	"reflect"
	"strings"
	"time"
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
		validate:           NewJSONValidator(),
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

func NewJSONValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}
