/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package storage

import (
	"github.com/google/uuid"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
)

type ExchangeRepository interface {
	StoreExchange(exchange *internal.Exchange)
	GetExchange(name string) (queue *internal.Exchange, err errs.AppError)
	FindExchanges() []internal.Exchange
	DeleteExchange(name string) (err errs.AppError)
	AddBinding(name string, binding internal.Binding) (err errs.AppError)
	DeleteBinding(name string, bindingId uuid.UUID) (err errs.AppError)
}
