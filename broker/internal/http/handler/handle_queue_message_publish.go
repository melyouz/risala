/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
)

func HandleQueueMessagePublish(queueRepository storage.QueueRepository, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var message internal.Message
		util.Decode(r, &message)

		var vErrors validator.ValidationErrors
		if errors.As(validate.Struct(message), &vErrors) {
			util.Respond(w, errs.NewValidationError(vErrors), http.StatusBadRequest)
			return
		}

		queueName := chi.URLParam(r, "queueName")
		queue, err := queueRepository.GetQueue(queueName)
		if err != nil {
			util.Respond(w, err, util.HttpStatusCodeFromAppError(err))
			return
		}

		queue.PublishMessage(message)
		queueRepository.StoreQueue(queue)

		util.Respond(w, nil, http.StatusOK)
	}
}
