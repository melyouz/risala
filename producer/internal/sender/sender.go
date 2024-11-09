/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package sender

import (
	"github.com/melyouz/risala/producer/internal"
	"github.com/melyouz/risala/producer/internal/errs"
)

type EventSender interface {
	Send(event internal.Event) errs.AppError
}
