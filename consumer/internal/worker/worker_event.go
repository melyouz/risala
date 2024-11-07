/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package worker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"

	"github.com/melyouz/risala/consumer/internal"
	"github.com/melyouz/risala/consumer/internal/action"
	"github.com/melyouz/risala/consumer/internal/util"
)

type EventWorker struct {
}

func NewEventWorker() *EventWorker {
	return &EventWorker{}
}

func (w *EventWorker) Start() {
	for {
		consumeMessages()
	}
}

func consumeMessages() {
	eventsQueueEndpoint := util.GetEnvVarStringRequired("EVENTS_QUEUE_ENDPOINT")
	messageConsumeEndpoint := fmt.Sprintf("%s/messages/get", eventsQueueEndpoint)
	response, connectionErr := http.Post(messageConsumeEndpoint, "application/json", nil)
	if connectionErr != nil {
		log.Println("[Worker] Error connecting to Broker:", connectionErr)
		return
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			log.Println("[Worker] Error closing response body:", closeErr)
		}
	}()

	if response.StatusCode == http.StatusNoContent {
		return
	}

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		log.Println("[Worker] Error reading response body:", readErr)
		return
	}

	var rawMessage internal.RawMessage
	if rawMessageDecodeErr := json.Unmarshal(body, &rawMessage); rawMessageDecodeErr != nil {
		log.Println("[Worker] Error decoding response body:", rawMessageDecodeErr)
		return
	}

	var event action.Event
	if deserializeErr := json.Unmarshal([]byte(rawMessage.Payload), &event); deserializeErr != nil {
		log.Println("[Worker] Error deserializing message payload:", deserializeErr)
		return
	}

	processEvent(rawMessage.Id, event)
}

func processEvent(messageId uuid.UUID, event action.Event) {
	fmt.Println("")
	log.Println("[Worker] Event process INIT:", event)

	eventHandled := handleEvent(event)

	if eventHandled {
		log.Println("[Worker] Event handled:", event.EventType)
		sendAcknowledgement(messageId, "ack")
	} else {
		log.Println("[Worker] Event not handled:", event.EventType)
		sendAcknowledgement(messageId, "nack")
	}

	log.Println("[Worker] Event process END:", event)
}

func handleEvent(event action.Event) (eventProcessed bool) {
	var eventActionFound bool

	for _, eventAction := range action.Actions {
		if util.WildcardMatch(eventAction.SupportedType(), event.EventType) {
			eventActionFound = true

			if err := eventAction.Handle(event); err != nil {
				log.Println("[Worker] Error handling event:", err)
			} else {
				eventProcessed = true
			}
		}
	}

	if !eventActionFound {
		log.Println("[Worker] Event action not found for event type:", event.EventType)
	}

	return eventProcessed
}

func sendAcknowledgement(messageId uuid.UUID, ackType string) {
	eventsQueueEndpoint := util.GetEnvVarStringRequired("EVENTS_QUEUE_ENDPOINT")
	messageEndpoint := fmt.Sprintf("%s/messages/%s/%s", eventsQueueEndpoint, messageId.String(), ackType)
	response, connectionErr := http.Post(messageEndpoint, "application/json", nil)
	if connectionErr != nil {
		log.Println("[Worker] An error occurred while connecting to Broker:", connectionErr)
		return
	}
	if response.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(response.Body)
		log.Printf("[Worker] Error %s-ing message: %s %s %s", ackType, messageEndpoint, response.Status, string(body))
		return
	}

	log.Printf("[Worker] Message %s %s-ed: %s %s", messageId, ackType, messageEndpoint, response.Status)
}
