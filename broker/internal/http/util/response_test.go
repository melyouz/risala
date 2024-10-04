/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/melyouz/risala/broker/internal"
)

func TestRespond(t *testing.T) {
	t.Run("Responds to HTTP request", func(t *testing.T) {
		queue := internal.Queue{Name: "testQueueName", Durability: internal.Durability.DURABLE}

		response := httptest.NewRecorder()
		Respond(response, queue, http.StatusOK)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, "{\"name\":\"testQueueName\",\"durability\":\"durable\"}", response.Body.String())
	})
}
