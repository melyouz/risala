/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package internal

import (
	"github.com/google/uuid"
)

type Message struct {
	Id      uuid.UUID `json:"id" validate:"required"`
	Payload string    `json:"payload" validate:"required"`
}
