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
	GetQueue(name string) (queue *internal.Queue, err errs.AppError)
	FindQueues() []internal.Queue
	DeleteQueue(name string) (err errs.AppError)
	PublishMessage(name string, message internal.Message) (err errs.AppError)
	GetMessages(name string, count int) (messages []internal.Message, err errs.AppError)
}
