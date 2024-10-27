/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package internal

import (
	"fmt"
	"sync"

	"github.com/melyouz/risala/broker/internal/errs"
)

type Queue struct {
	sync.RWMutex
	Name       string         `json:"name" validate:"required"`
	Durability DurabilityType `json:"durability" validate:"required,oneof=durable transient"`
	Messages   []*Message     `json:"-" validate:"dive"`
}

func (q *Queue) Enqueue(message *Message) (err errs.AppError) {
	q.Lock()
	defer q.Unlock()

	q.Messages = append(q.Messages, message)

	return nil
}

func (q *Queue) Dequeue() (message *Message, err errs.AppError) {
	q.Lock()
	defer q.Unlock()

	messagesCount := len(q.Messages)
	if messagesCount == 0 {
		return nil, errs.NewQueueEmptyError(fmt.Sprintf("Queue '%s' is empty", q.Name))
	}

	for _, m := range q.Messages {
		if !m.IsAwaiting() {
			m.Await()
			return m, nil
		}
	}

	return nil, errs.NewMessageNotFoundError("No pending messages available")
}

func (q *Queue) Seek(count int) (messages []*Message, err errs.AppError) {
	q.Lock()
	defer q.Unlock()

	messagesCount := len(q.Messages)
	if messagesCount == 0 {
		return messages, nil
	}

	if count > messagesCount {
		count = messagesCount
	}

	var result []*Message
	for _, m := range q.Messages {
		if !m.IsAwaiting() && len(result) < count {
			result = append(result, m)
		}
	}

	return result, nil
}
