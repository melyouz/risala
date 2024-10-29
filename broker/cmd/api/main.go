/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package main

import (
	"flag"
	"log"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-playground/validator/v10"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/http/server"
	"github.com/melyouz/risala/broker/internal/sample"
	"github.com/melyouz/risala/broker/internal/storage"
)

func main() {
	listenAddr := "localhost:8000"
	router := chi.NewRouter()

	withSampleData := flag.Bool("with-sample-data", false, "Initialize API with sample data")
	flag.Parse()

	queues := map[string]*internal.Queue{}
	exchanges := map[string]*internal.Exchange{}
	if *withSampleData {
		queues = sample.Queues
		exchanges = sample.Exchanges
	}

	queueRepository := storage.NewInMemoryQueueRepository(queues)
	exchangeRepository := storage.NewInMemoryExchangeRepository(exchanges)

	s := server.NewServer(listenAddr, router, queueRepository, exchangeRepository)
	log.Printf("Listening on: http://%s\n", listenAddr)
	log.Printf("With sample data: %v", *withSampleData)
	log.Fatal(s.ListenAndServe())
}
