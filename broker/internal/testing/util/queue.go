/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package util

import (
	"github.com/melyouz/risala/broker/internal"
)

func NewTestQueueDurableWithoutMessages(name string) (queue *internal.Queue) {
	return &internal.Queue{
		Name:       name,
		Durability: internal.Durability.DURABLE,
		Messages:   []*internal.Message{},
	}
}

func NewTestQueueTransientWithoutMessages(name string) (queue *internal.Queue) {
	return &internal.Queue{
		Name:       name,
		Durability: internal.Durability.TRANSIENT,
		Messages:   []*internal.Message{},
	}
}

func NewTestQueueTransientWithMessages(name string, messages []*internal.Message) (queue *internal.Queue) {
	return &internal.Queue{
		Name:       name,
		Durability: internal.Durability.TRANSIENT,
		Messages:   messages,
	}
}

func NewTestSystemQueueWithoutMessages(name string) (queue *internal.Queue) {
	return &internal.Queue{
		Name:       name,
		Durability: internal.Durability.DURABLE,
		Messages:   []*internal.Message{},
		System:     true,
	}
}
