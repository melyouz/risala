/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package internal

import (
	"slices"
	"sync"
)

var queueLock = &sync.RWMutex{}

type Queue struct {
	Name       string         `json:"name" validate:"required"`
	Durability DurabilityType `json:"durability" validate:"required,oneof=durable transient"`
	Messages   []Message      `json:"-" validate:"dive"`
}

func (q *Queue) PublishMessage(message Message) {
	queueLock.Lock()
	defer queueLock.Unlock()

	q.Messages = append(q.Messages, message)
}

func (q *Queue) GetMessages(desiredCount int) (messages []Message) {
	queueLock.Lock()
	defer queueLock.Unlock()

	messagesCount := len(q.Messages)
	if messagesCount == 0 {
		return []Message{}
	}

	if desiredCount > messagesCount {
		desiredCount = messagesCount
	}

	tmp := q.Messages[messagesCount-desiredCount:]

	messages = make([]Message, len(tmp))
	copy(messages, tmp)
	slices.Reverse(messages)

	return messages
}
