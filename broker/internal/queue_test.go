/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package internal

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestQueueConcurrency(t *testing.T) {
	t.Parallel()
	t.Run("Messages are marked as processing when queued & dequeued", func(t *testing.T) {
		t.Parallel()
		q := &Queue{Name: "testQueue", Durability: Durability.DURABLE}
		var wg sync.WaitGroup
		numOperations := 1000
		dequeueMaxRetries := 10

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

				var message *Message
				for retries := 0; retries < dequeueMaxRetries; retries++ {
					message = q.Dequeue()
					if message != nil {
						break
					}
					time.Sleep(10 * time.Millisecond)
				}
				assert.NotNil(t, message)
				assert.True(t, message.IsProcessing())
			}()
		}

		wg.Wait()

		assert.Len(t, q.Messages, numOperations)
		for _, m := range q.Messages {
			assert.True(t, m.IsProcessing(), "All messages should be marked as processing")
		}
	})

	t.Run("Messages are removed from the queue when acknowledged", func(t *testing.T) {
		t.Parallel()
		q := &Queue{Name: "testQueue", Durability: Durability.DURABLE}
		var wg sync.WaitGroup
		numOperations := 1000
		dequeueMaxRetries := 10

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

				var message *Message
				for retries := 0; retries < dequeueMaxRetries; retries++ {
					message = q.Dequeue()
					if message != nil {
						break
					}
					time.Sleep(10 * time.Millisecond)
				}
				assert.NotNil(t, message)
				assert.True(t, message.IsProcessing())

				ackErr := q.Ack(message.Id)
				assert.Nil(t, ackErr)
			}()
		}

		wg.Wait()

		assert.Empty(t, q.Messages)
	})
}
