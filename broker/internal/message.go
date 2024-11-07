/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package internal

import (
	"sync"

	"github.com/google/uuid"
)

type Message struct {
	sync.Mutex
	Id         uuid.UUID `json:"id" validate:"required"`
	Payload    string    `json:"payload" validate:"required"`
	Processing bool      `json:"isProcessing"`
}

func (m *Message) MarkProcessing() {
	m.Lock()
	defer m.Unlock()

	m.Processing = true
}

func (m *Message) UnmarkProcessing() {
	m.Lock()
	defer m.Unlock()

	m.Processing = false
}

func (m *Message) IsProcessing() bool {
	m.Lock()
	defer m.Unlock()

	return m.Processing
}
