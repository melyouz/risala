/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package action

import (
	"fmt"

	"github.com/melyouz/risala/consumer/internal/errs"
)

type ProductCreatedAction struct {
}

func (ProductCreatedAction) SupportedType() string {
	return "product.created"
}

func (action ProductCreatedAction) Handle(event Event) errs.AppError {
	fmt.Println("[ProductCreatedAction] Event handled: ", event)

	return nil
}
