/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package internal

type Queue struct {
	Name       string         `json:"name" validate:"required"`
	Durability DurabilityType `json:"durability" validate:"required,oneof=durable transient"`
	Messages   []Message      `json:"-" validate:"dive"`
}
