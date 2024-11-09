/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/melyouz/risala/producer/internal"
	"github.com/melyouz/risala/producer/internal/errs"
	"github.com/melyouz/risala/producer/internal/util"
)

type HTTPEventSender struct{}

func NewHttpEventSender() *HTTPEventSender {
	return &HTTPEventSender{}
}

func (s *HTTPEventSender) Send(event internal.Event) errs.AppError {
	internalExchangeEndpoint := util.GetEnvVarStringRequired("EXCHANGE_INTERNAL_ENDPOINT")
	messagePublishEndpoint := fmt.Sprintf("%s/messages/publish", internalExchangeEndpoint)

	encodedEvent, eventEncodeErr := json.Marshal(event)
	if eventEncodeErr != nil {
		return errs.NewEncodeError(fmt.Sprintf("Error encoding event: %s", eventEncodeErr))
	}

	message := internal.Message{Payload: string(encodedEvent)}
	encodedMessage, messageEncodeErr := json.Marshal(message)
	if messageEncodeErr != nil {
		return errs.NewEncodeError(fmt.Sprintf("Error encoding message: %s", messageEncodeErr))
	}

	response, connectionErr := http.Post(messagePublishEndpoint, "application/json", bytes.NewBuffer(encodedMessage))
	if connectionErr != nil {
		return errs.NewConnectionError(fmt.Sprintf("Error connecting to Broker: %s", connectionErr))
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			log.Println("[Sender] Error closing response body:", closeErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(response.Body)
		if readErr != nil {
			return errs.NewReadError(fmt.Sprintf("Error reading response body: %s", readErr))
		}

		var apiErr errs.Error
		if decodeErr := json.Unmarshal(body, &apiErr); decodeErr != nil {
			return errs.NewDecodeError(fmt.Sprintf("Error decoding response body: %s", decodeErr))
		}

		return &apiErr
	}

	return nil
}
