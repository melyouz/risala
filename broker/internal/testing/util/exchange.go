/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package util

import (
	"github.com/melyouz/risala/broker/internal"
)

func NewTestExchangeWithoutBindings(name string) (queue *internal.Exchange) {
	return &internal.Exchange{
		Name:     name,
		Bindings: []*internal.Binding{},
	}
}

func NewTestExchangeWithBindings(name string, bindings []*internal.Binding) (queue *internal.Exchange) {
	return &internal.Exchange{
		Name:     name,
		Bindings: bindings,
	}
}
