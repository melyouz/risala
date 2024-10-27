/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package internal

import (
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestQueueConcurrency(t *testing.T) {
	t.Parallel()
	t.Run("Enqueues & dequeues messages concurrently", func(t *testing.T) {
		t.Parallel()
		q := &Queue{Name: "testQueue", Durability: Durability.DURABLE}
		var wg sync.WaitGroup
		numOperations := 1000

		wg.Add(numOperations)
		for i := 0; i < numOperations; i++ {
			go func(i int) {
				defer wg.Done()
				message := &Message{Id: uuid.New(), Payload: fmt.Sprintf("Message %d", i)}
				enqueueErr := q.Enqueue(message)
				assert.Nil(t, enqueueErr)
			}(i)
		}

		wg.Add(numOperations)
		for i := 0; i < numOperations; i++ {
			go func() {
				defer wg.Done()
				message, dequeueErr := q.Dequeue()
				assert.Nil(t, dequeueErr)
				assert.NotNil(t, message)
				assert.True(t, message.IsAwaiting())
			}()
		}

		wg.Wait()

		assert.Len(t, q.Messages, numOperations)
		for _, m := range q.Messages {
			assert.True(t, m.IsAwaiting(), "All messages should be marked as awaiting")
		}
	})
}
