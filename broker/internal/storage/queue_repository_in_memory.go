/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package storage

import (
	"fmt"
	"sort"
	"sync"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
)

type InMemoryQueueRepository struct {
	lock      *sync.RWMutex
	QueueList map[string]internal.Queue
}

func NewInMemoryQueueRepository(queueList map[string]internal.Queue) *InMemoryQueueRepository {
	return &InMemoryQueueRepository{
		lock:      &sync.RWMutex{},
		QueueList: queueList,
	}
}

func (r *InMemoryQueueRepository) FindQueues() []internal.Queue {
	r.lock.Lock()
	defer r.lock.Unlock()

	result := make([]internal.Queue, 0, len(r.QueueList))

	for _, value := range r.QueueList {
		result = append(result, value)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

func (r *InMemoryQueueRepository) StoreQueue(queue *internal.Queue) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.QueueList[queue.Name] = *queue
}

func (r *InMemoryQueueRepository) GetQueue(name string) (queue *internal.Queue, err errs.AppError) {
	r.lock.Lock()
	defer r.lock.Unlock()

	q, ok := r.QueueList[name]
	if ok {
		return &q, err
	}

	return nil, errs.NewQueueNotFoundError(fmt.Sprintf("Queue '%s' not found", name))
}

func (r *InMemoryQueueRepository) DeleteQueue(name string) (err errs.AppError) {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.QueueList[name]
	if ok {
		delete(r.QueueList, name)
		return err
	}

	return errs.NewQueueNotFoundError(fmt.Sprintf("Queue '%s' not found", name))
}
