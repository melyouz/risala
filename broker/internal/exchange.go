/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package internal

type Exchange struct {
	Name     string    `json:"name" validate:"required"`
	Bindings []Binding `json:"bindings"`
}
