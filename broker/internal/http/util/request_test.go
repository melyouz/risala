/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/melyouz/risala/broker/internal"
)

func TestDecode(t *testing.T) {
	t.Run("Decodes JSON body from HTTP request", func(t *testing.T) {
		queueBody, _ := json.Marshal(map[string]interface{}{
			"name":       "testQueueName",
			"durability": "transient",
		})
		request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(queueBody))

		var queue internal.Queue
		Decode(request, &queue)

		assert.Equal(t, "testQueueName", queue.Name)
		assert.Equal(t, internal.Durability.TRANSIENT, queue.Durability)
		assert.Empty(t, queue.Messages)
	})
}
