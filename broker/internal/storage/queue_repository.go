package storage

import (
	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
)

type QueueRepository interface {
	StoreQueue(queue *internal.Queue)
	GetQueue(name string) (queue *internal.Queue, err errs.AppError)
	FindQueues() []internal.Queue
	DeleteQueue(name string) (err errs.AppError)
}
