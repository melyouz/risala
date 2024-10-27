/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package internal

import (
	"fmt"
	"slices"
	"sync"

	"github.com/google/uuid"

	"github.com/melyouz/risala/broker/internal/errs"
)

type Exchange struct {
	sync.RWMutex
	Name     string     `json:"name" validate:"required"`
	Bindings []*Binding `json:"bindings"`
}

func (e *Exchange) Bind(binding *Binding) (err errs.AppError) {
	e.Lock()
	defer e.Unlock()

	bindingErr := validateBindingDoesNotExist(e, binding)
	if bindingErr != nil {
		return bindingErr
	}

	e.Bindings = append(e.Bindings, binding)

	return nil
}

func (e *Exchange) Unbind(bindingId uuid.UUID) (err errs.AppError) {
	e.Lock()
	defer e.Unlock()

	for i, binding := range e.Bindings {
		if binding.Id == bindingId {
			e.Bindings = slices.Delete(e.Bindings, i, i+1)
			return nil
		}
	}

	return errs.NewBindingNotFoundError(fmt.Sprintf("Binding '%s' not found", bindingId))
}

func validateBindingDoesNotExist(exchange *Exchange, binding *Binding) errs.AppError {
	for _, v := range exchange.Bindings {
		if v.Queue == binding.Queue {
			return errs.NewBindingExistsError(fmt.Sprintf("Binding to Queue '%s' already exists", binding.Queue))
		}
	}

	return nil
}
