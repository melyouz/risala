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

const DeadLetterQueueName = "system.dead-letter"

type Queue struct {
	sync.RWMutex
	Name       string         `json:"name" validate:"required"`
	Durability DurabilityType `json:"durability" validate:"required,oneof=durable transient"`
	Messages   []*Message     `json:"-" validate:"dive"`
	System     bool           `json:"isSystem"`
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
		if !m.IsProcessing() {
			m.MarkProcessing()
			return m
		}
	}

	return nil
}

func (q *Queue) Ack(messageId uuid.UUID) (err errs.AppError) {
	q.Lock()
	defer q.Unlock()

	for i, m := range q.Messages {
		if m.Id == messageId && m.IsProcessing() {
			q.Messages = slices.Delete(q.Messages, i, i+1)
			return nil
		}
	}

	return errs.NewMessageNotFoundError(fmt.Sprintf("Message '%s' not found", messageId.String()))
}

func (q *Queue) Nack(messageId uuid.UUID) (message *Message, err errs.AppError) {
	q.Lock()
	defer q.Unlock()

	for i, m := range q.Messages {
		if m.Id == messageId && m.IsProcessing() {
			m.UnmarkProcessing()
			q.Messages = slices.Delete(q.Messages, i, i+1)
			return m, nil
		}
	}

	return nil, errs.NewMessageNotFoundError(fmt.Sprintf("Message '%s' not found", messageId.String()))
}

func (q *Queue) Peek(limit int) (messages []*Message, err errs.AppError) {
	q.Lock()
	defer q.Unlock()

	messagesCount := len(q.Messages)
	if messagesCount == 0 || limit <= 0 {
		return make([]*Message, 0), nil
	}

	if limit > messagesCount {
		limit = messagesCount
	}

	return q.Messages[:limit], nil
}

func (q *Queue) Purge() (err errs.AppError) {
	q.Lock()
	defer q.Unlock()

	q.Messages = []*Message{}

	return nil
}

func (q *Queue) IsSystem() bool {
	return q.System
}
