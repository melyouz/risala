/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package action

import (
	"fmt"

	"github.com/melyouz/risala/consumer/internal/errs"
)

type ProductPublishedAction struct {
}

func (ProductPublishedAction) SupportedType() string {
	return "product.published"
}

func (action ProductPublishedAction) Handle(event Event) errs.AppError {
	fmt.Println("[ProductPublishedAction] Event handled: ", event)

	return nil
}
