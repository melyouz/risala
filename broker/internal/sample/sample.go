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
	"events": {Name: "events", Durability: internal.Durability.DURABLE, Messages: []*internal.Message{
		{Id: uuid.New(), Payload: "Message 1 (events)"},
		{Id: uuid.New(), Payload: "Message 2 (events)"},
		{Id: uuid.New(), Payload: "Message 3 (events)"},
		{Id: uuid.New(), Payload: "Message 4 (events)"},
		{Id: uuid.New(), Payload: "Message 5 (events)"},
	}},
	"tmp": {Name: "tmp", Durability: internal.Durability.TRANSIENT, Messages: []*internal.Message{
		{Id: uuid.New(), Payload: "Message 1 (tmp)"},
		{Id: uuid.New(), Payload: "Message 2 (tmp)"},
		{Id: uuid.New(), Payload: "Message 3 (tmp)"},
		{Id: uuid.New(), Payload: "Message 4 (tmp)"},
		{Id: uuid.New(), Payload: "Message 5 (tmp)"},
	}},
}
