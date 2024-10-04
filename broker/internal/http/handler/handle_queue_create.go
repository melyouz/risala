/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
	"net/http"
)

func HandleQueueCreate(queueRepository storage.QueueRepository, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var queue internal.Queue
		util.Decode(r, &queue)

		var vErrors validator.ValidationErrors
		if errors.As(validate.Struct(queue), &vErrors) {
			util.Respond(w, errs.NewValidationError(vErrors), http.StatusBadRequest)
			return
		}

		queueRepository.StoreQueue(&queue)

		util.Respond(w, queue, http.StatusCreated)
	}
}
