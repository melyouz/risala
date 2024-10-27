/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package storage

import (
	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
)

type ExchangeRepository interface {
	StoreExchange(exchange *internal.Exchange)
	GetExchange(name string) (queue *internal.Exchange, err errs.AppError)
	FindExchanges() []*internal.Exchange
	DeleteExchange(name string) (err errs.AppError)
}
