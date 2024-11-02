/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package util

import (
	"github.com/melyouz/risala/broker/internal"
)

func NewNewQueueDurableWithoutMessages(name string) (queue *internal.Queue) {
	return &internal.Queue{
		Name:       name,
		Durability: internal.Durability.DURABLE,
		Messages:   []*internal.Message{},
	}
}

func NewQueueTransientWithoutMessages(name string) (queue *internal.Queue) {
	return &internal.Queue{
		Name:       name,
		Durability: internal.Durability.TRANSIENT,
		Messages:   []*internal.Message{},
	}
}

func NewQueueTransientWithMessages(name string, messages []*internal.Message) (queue *internal.Queue) {
	return &internal.Queue{
		Name:       name,
		Durability: internal.Durability.TRANSIENT,
		Messages:   messages,
	}
}
