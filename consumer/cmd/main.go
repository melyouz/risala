/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/melyouz/risala/consumer/internal/worker"
)

func main() {
	eventWorker := worker.NewEventWorker()
	eventWorker.Start()
}
