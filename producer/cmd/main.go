/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package main

import (
	"flag"
	"log"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/google/uuid"

	"github.com/melyouz/risala/producer/internal"
	"github.com/melyouz/risala/producer/internal/sender"
)

func main() {
	log.Println("Sending events...")

	eventsCount := flag.Int("events-count", 1000, "Send EVENTS-COUNT events")
	flag.Parse()

	log.Printf("eventsCount: %d", *eventsCount)

	successCount := 0
	eventSender := sender.NewHttpEventSender()
	for i := 0; i < *eventsCount; i++ {
		event := internal.Event{
			Id:        uuid.New(),
			EventType: "product.published",
			Data: map[string]interface{}{
				"productId": uuid.New().String(),
				"actorId":   uuid.New().String(),
			},
			Timestamp: time.Now().Unix(),
		}

		err := eventSender.Send(event)
		if err != nil {
			log.Printf("failed to send event: %v", err)
		} else {
			successCount++
		}
	}

	log.Printf("Successfully sent %d events", successCount)
}
