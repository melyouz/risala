package internal

type DurabilityType string

var Durability = struct {
	DURABLE   DurabilityType
	TRANSIENT DurabilityType
}{
	DURABLE:   "durable",
	TRANSIENT: "transient",
}
