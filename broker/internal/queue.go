/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package internal

import (
	"fmt"
	"slices"
	"sync"

	"github.com/google/uuid"

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

func (q *Queue) Dequeue() (message *Message) {
	q.Lock()
	defer q.Unlock()

	messagesCount := len(q.Messages)
	if messagesCount == 0 {
		return nil
	}

	for _, m := range q.Messages {
		if !m.IsAwaiting() {
			m.Await()
			return m
		}
	}

	return nil
}

func (q *Queue) Ack(messageId uuid.UUID) (err errs.AppError) {
	q.Lock()
	defer q.Unlock()

	for i, m := range q.Messages {
		if m.Id == messageId && m.IsAwaiting() {
			q.Messages = slices.Delete(q.Messages, i, i+1)
			return nil
		}
	}

	return errs.NewMessageNotFoundError(fmt.Sprintf("Message '%s' not found", messageId.String()))
}

func (q *Queue) Peek(limit int) (messages []*Message, err errs.AppError) {
	q.Lock()
	defer q.Unlock()

	result := make([]*Message, 0)

	messagesCount := len(q.Messages)
	if messagesCount == 0 {
		return result, nil
	}

	if limit > messagesCount {
		limit = messagesCount
	}

	for _, m := range q.Messages {
		if !m.IsAwaiting() && len(result) < limit {
			result = append(result, m)
		}
	}

	return result, nil
}
