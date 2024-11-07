/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package internal

import (
	"github.com/google/uuid"
)

type RawMessage struct {
	Id      uuid.UUID `json:"id"`
	Payload string    `json:"payload"`
}
