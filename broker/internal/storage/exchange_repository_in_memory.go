/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package storage

import (
	"fmt"
	"slices"
	"sort"
	"sync"

	"github.com/google/uuid"

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
	if !ok {
		return nil, errs.NewExchangeNotFoundError(fmt.Sprintf("Exchange '%s' not found", name))
	}

	return &q, nil
}

func (r *InMemoryExchangeRepository) DeleteExchange(name string) (err errs.AppError) {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.ExchangeList[name]
	if !ok {
		return errs.NewExchangeNotFoundError(fmt.Sprintf("Exchange '%s' not found", name))
	}

	delete(r.ExchangeList, name)
	return nil
}

func (r *InMemoryExchangeRepository) AddBinding(name string, binding internal.Binding) (err errs.AppError) {
	r.lock.Lock()
	defer r.lock.Unlock()

	exchange, ok := r.ExchangeList[name]
	if !ok {
		return errs.NewExchangeNotFoundError(fmt.Sprintf("Exchange '%s' not found", name))
	}

	bindingErr := validateBindingDoesNotExist(exchange, binding)
	if bindingErr != nil {
		return bindingErr
	}

	exchange.Bindings = append(exchange.Bindings, binding)
	r.ExchangeList[name] = exchange

	return nil
}

func (r *InMemoryExchangeRepository) DeleteBinding(name string, bindingId uuid.UUID) (err errs.AppError) {
	r.lock.Lock()
	defer r.lock.Unlock()

	exchange, ok := r.ExchangeList[name]
	if !ok {
		return errs.NewExchangeNotFoundError(fmt.Sprintf("Exchange '%s' not found", name))
	}

	for i, binding := range exchange.Bindings {
		if binding.Id == bindingId {
			exchange.Bindings = slices.Delete(exchange.Bindings, i, i+1)
			r.ExchangeList[name] = exchange
			return nil
		}
	}

	return errs.NewBindingNotFoundError(fmt.Sprintf("Binding '%s' not found", bindingId))
}

func validateBindingDoesNotExist(exchange internal.Exchange, binding internal.Binding) errs.AppError {
	for _, v := range exchange.Bindings {
		if v.Queue == binding.Queue {
			return errs.NewBindingExistsError(fmt.Sprintf("Binding to Queue '%s' already exists", binding.Queue))
		}
	}

	return nil
}
