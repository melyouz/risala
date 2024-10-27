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
	Id       uuid.UUID `json:"id" validate:"required"`
	Payload  string    `json:"payload" validate:"required"`
	Awaiting bool      `json:"-"`
}

func (m *Message) Await() {
	m.Lock()
	defer m.Unlock()

	m.Awaiting = true
}

func (m *Message) IsAwaiting() bool {
	m.Lock()
	defer m.Unlock()

	return m.Awaiting
}
