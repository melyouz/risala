/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package internal

import (
	"fmt"
	"slices"

	"github.com/google/uuid"

	"github.com/melyouz/risala/broker/internal/errs"
)

type Exchange struct {
	Name     string    `json:"name" validate:"required"`
	Bindings []Binding `json:"bindings"`
}

func (e *Exchange) AddBinding(binding Binding) {
	e.Bindings = append(e.Bindings, binding)
}

func (e *Exchange) RemoveBinding(id uuid.UUID) errs.AppError {
	for i, binding := range e.Bindings {
		if binding.Id == id {
			e.Bindings = slices.Delete(e.Bindings, i, i+1)
			return nil
		}
	}

	return errs.NewBindingNotFoundError(fmt.Sprintf("Binding '%s' not found", id))
}
