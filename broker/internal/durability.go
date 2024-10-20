/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package internal

type DurabilityType string

var Durability = struct {
	DURABLE   DurabilityType
	TRANSIENT DurabilityType
}{
	DURABLE:   "durable",
	TRANSIENT: "transient",
}

func (d *DurabilityType) String() string {
	return string(*d)
}
