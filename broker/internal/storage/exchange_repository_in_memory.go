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

type InMemoryExchangeRepository struct {
	lock         *sync.RWMutex
	ExchangeList map[string]internal.Exchange
}

func NewInMemoryExchangeRepository(exchangeList map[string]internal.Exchange) *InMemoryExchangeRepository {
	return &InMemoryExchangeRepository{
		lock:         &sync.RWMutex{},
		ExchangeList: exchangeList,
	}
}

func (r *InMemoryExchangeRepository) FindExchanges() []internal.Exchange {
	r.lock.Lock()
	defer r.lock.Unlock()

	result := make([]internal.Exchange, 0, len(r.ExchangeList))

	for _, value := range r.ExchangeList {
		result = append(result, value)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

func (r *InMemoryExchangeRepository) StoreExchange(exchange *internal.Exchange) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.ExchangeList[exchange.Name] = *exchange
}

func (r *InMemoryExchangeRepository) GetExchange(name string) (exchange *internal.Exchange, err errs.AppError) {
	r.lock.Lock()
	defer r.lock.Unlock()

	q, ok := r.ExchangeList[name]
	if ok {
		return &q, nil
	}

	return nil, errs.NewExchangeNotFoundError(fmt.Sprintf("Exchange '%s' not found", name))
}

func (r *InMemoryExchangeRepository) DeleteExchange(name string) (err errs.AppError) {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.ExchangeList[name]
	if ok {
		delete(r.ExchangeList, name)
		return err
	}

	return errs.NewExchangeNotFoundError(fmt.Sprintf("Exchange '%s' not found", name))
}
