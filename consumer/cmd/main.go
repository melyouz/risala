/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package main

import (
	"github.com/melyouz/risala/consumer/internal/worker"
)

func main() {
	eventWorker := worker.NewEventWorker()
	eventWorker.Start()
}
