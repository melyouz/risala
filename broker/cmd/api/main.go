/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package main

import (
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-playground/validator/v10"

	"github.com/melyouz/risala/broker/internal/http/server"
	"github.com/melyouz/risala/broker/internal/sample"
	"github.com/melyouz/risala/broker/internal/storage"
)

func main() {
	listenAddr := "localhost:8000"
	router := chi.NewRouter()
	queueRepository := storage.NewInMemoryQueueRepository(sample.Queues)
	exchangeRepository := storage.NewInMemoryExchangeRepository(sample.Exchanges)

	s := server.NewServer(listenAddr, router, queueRepository, exchangeRepository)
	fmt.Printf("Listening on: http://%s\n", listenAddr)
	log.Fatal(s.ListenAndServe())
}
