/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package action

import (
	"github.com/melyouz/risala/consumer/internal"
	"github.com/melyouz/risala/consumer/internal/errs"
)

type Action interface {
	SupportedType() string
	Handle(action internal.Event) errs.AppError
}
