/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package sample

import (
	"github.com/google/uuid"

	"github.com/melyouz/risala/broker/internal"
)

var Exchanges = map[string]*internal.Exchange{
	"app.internal": {Name: "app.internal", Bindings: []*internal.Binding{
		{Id: uuid.New(), Queue: "events", RoutingKey: "#"},
	}},
	"app.external": {Name: "app.external", Bindings: []*internal.Binding{
		{Id: uuid.New(), Queue: "tmp", RoutingKey: "#"},
	}},
}

var Queues = map[string]*internal.Queue{
	"events": {
		Name:       "events",
		Durability: internal.Durability.DURABLE,
		Messages:   []*internal.Message{},
	},
	"tmp": {
		Name:       "tmp",
		Durability: internal.Durability.TRANSIENT,
		Messages: []*internal.Message{
			{Id: uuid.New(), Payload: "Message 1 (tmp)"},
			{Id: uuid.New(), Payload: "Message 2 (tmp)"},
			{Id: uuid.New(), Payload: "Message 3 (tmp)"},
			{Id: uuid.New(), Payload: "Message 4 (tmp)"},
			{Id: uuid.New(), Payload: "Message 5 (tmp)"},
		},
	},
	internal.DeadLetterQueueName: {
		Name:       internal.DeadLetterQueueName,
		Durability: internal.Durability.DURABLE,
		System:     true,
		Messages:   []*internal.Message{},
	},
}
