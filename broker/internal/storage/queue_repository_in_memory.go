/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package storage

import (
	"fmt"
	"slices"
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
	if !ok {
		return nil, errs.NewQueueNotFoundError(fmt.Sprintf("Queue '%s' not found", name))
	}

	return &q, err
}

func (r *InMemoryQueueRepository) DeleteQueue(name string) (err errs.AppError) {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.QueueList[name]
	if !ok {
		return errs.NewQueueNotFoundError(fmt.Sprintf("Queue '%s' not found", name))
	}

	delete(r.QueueList, name)
	return err
}

func (r *InMemoryQueueRepository) PublishMessage(name string, message internal.Message) (err errs.AppError) {
	r.lock.Lock()
	defer r.lock.Unlock()

	q, ok := r.QueueList[name]
	if !ok {
		return errs.NewQueueNotFoundError(fmt.Sprintf("Queue '%s' not found", name))
	}

	q.Messages = append(q.Messages, message)
	r.QueueList[name] = q

	return nil
}

func (r *InMemoryQueueRepository) GetMessages(name string, count int) (messages []internal.Message, err errs.AppError) {
	r.lock.Lock()
	defer r.lock.Unlock()

	q, ok := r.QueueList[name]
	if !ok {
		return nil, errs.NewQueueNotFoundError(fmt.Sprintf("Queue '%s' not found", name))
	}

	messagesCount := len(q.Messages)
	if messagesCount == 0 {
		return []internal.Message{}, nil
	}

	if count > messagesCount {
		count = messagesCount
	}

	tmp := q.Messages[messagesCount-count:]

	messages = make([]internal.Message, len(tmp))
	copy(messages, tmp)
	slices.Reverse(messages)

	return messages, nil
}
