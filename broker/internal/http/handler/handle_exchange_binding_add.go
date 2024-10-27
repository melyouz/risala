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

func HandleExchangeBindingAdd(exchangeRepository storage.ExchangeRepository, queueRepository storage.QueueRepository, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var binding internal.Binding
		binding.Id = uuid.New()
		util.Decode(r, &binding)

		var vErrors validator.ValidationErrors
		if errors.As(validate.Struct(binding), &vErrors) {
			util.Respond(w, errs.NewValidationError(vErrors), http.StatusBadRequest)
			return
		}

		_, queueErr := queueRepository.GetQueue(binding.Queue)
		if queueErr != nil {
			util.Respond(w, queueErr, util.HttpStatusCodeFromAppError(queueErr))
			return
		}

		exchangeName := chi.URLParam(r, "exchangeName")
		exchange, exchangeErr := exchangeRepository.GetExchange(exchangeName)
		if exchangeErr != nil {
			util.Respond(w, exchangeErr, util.HttpStatusCodeFromAppError(exchangeErr))
			return
		}

		bindErr := exchange.Bind(&binding)
		if bindErr != nil {
			util.Respond(w, bindErr, util.HttpStatusCodeFromAppError(bindErr))
			return
		}

		util.Respond(w, binding, http.StatusCreated)
	}
}
