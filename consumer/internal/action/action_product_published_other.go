/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package action

import (
	"fmt"

	"github.com/melyouz/risala/consumer/internal"
	"github.com/melyouz/risala/consumer/internal/errs"
)

type ProductPublishedOtherAction struct {
}

func (ProductPublishedOtherAction) SupportedType() string {
	return "product.published"
}

func (action ProductPublishedOtherAction) Handle(event internal.Event) errs.AppError {
	fmt.Println("[ProductPublishedOtherAction] Event handled: ", event)

	return nil
}
