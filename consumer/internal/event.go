/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package internal

import (
	"github.com/google/uuid"
)

type Event struct {
	Id        uuid.UUID              `json:"id"`
	EventType string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
}
