/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package storage

import (
	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
)

type QueueRepository interface {
	StoreQueue(queue *internal.Queue)
	FindQueues() []*internal.Queue
	GetQueue(name string) (queue *internal.Queue, err errs.AppError)
	DeleteQueue(name string) (err errs.AppError)
}
