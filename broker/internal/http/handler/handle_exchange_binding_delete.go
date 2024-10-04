/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/melyouz/risala/broker/internal/errs"
	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
	"net/http"
)

func HandleExchangeBindingDelete(exchangeRepository storage.ExchangeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exchangeName := chi.URLParam(r, "exchangeName")
		exchange, exchangeErr := exchangeRepository.GetExchange(exchangeName)
		if exchangeErr != nil {
			util.Respond(w, exchangeErr, util.HttpStatusCodeFromAppError(exchangeErr))
			return
		}

		bindingIdParam := "bindingId"
		bindingId, uuidErr := uuid.Parse(chi.URLParam(r, bindingIdParam))
		if uuidErr != nil {
			paramErr := errs.NewParamInvalidError(uuidErr.Error())
			util.Respond(w, paramErr, util.HttpStatusCodeFromAppError(paramErr))
			return
		}

		err := exchange.RemoveBinding(bindingId)
		if err != nil {
			util.Respond(w, err, util.HttpStatusCodeFromAppError(err))
			return
		}

		exchangeRepository.StoreExchange(exchange)

		util.Respond(w, nil, http.StatusAccepted)
	}
}
