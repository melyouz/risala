/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
)

func HandleExchangeMessagePublish(
	exchangeRepository storage.ExchangeRepository,
	queueRepository storage.QueueRepository,
	validate *validator.Validate,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var message internal.Message
		message.Id = uuid.New()
		util.Decode(r, &message)

		var vErrors validator.ValidationErrors
		if errors.As(validate.Struct(&message), &vErrors) {
			util.Respond(w, errs.NewValidationError(vErrors), http.StatusBadRequest)
			return
		}

		exchangeName := chi.URLParam(r, "exchangeName")
		exchange, exchangeErr := exchangeRepository.GetExchange(exchangeName)
		if exchangeErr != nil {
			util.Respond(w, exchangeErr, util.HttpStatusCodeFromAppError(exchangeErr))
			return
		}

		for _, binding := range exchange.Bindings {
			queue, queueErr := queueRepository.GetQueue(binding.Queue)
			if queueErr != nil {
				util.Respond(w, queueErr, util.HttpStatusCodeFromAppError(queueErr))
				return
			}

			publishErr := queue.Enqueue(&message)
			if publishErr != nil {
				util.Respond(w, publishErr, util.HttpStatusCodeFromAppError(publishErr))
				return
			}
		}

		util.Respond(w, nil, http.StatusOK)
	}
}
