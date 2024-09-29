package sample

import "github.com/melyouz/risala/broker/internal"

var Exchanges = map[string]internal.Exchange{
	"app.internal": {Name: "app.internal", Bindings: []internal.Binding{}},
	"app.external": {Name: "app.external", Bindings: []internal.Binding{}},
}

var Queues = map[string]internal.Queue{
	"events": {Name: "events", Durability: internal.Durability.DURABLE, Messages: []internal.Message{}},
	"tmp":    {Name: "tmp", Durability: internal.Durability.TRANSIENT, Messages: []internal.Message{}},
}
