/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/melyouz/risala/broker/internal/errs"
	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
)

func HandleExchangeBindingDelete(exchangeRepository storage.ExchangeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bindingIdParamName := "bindingId"
		bindingId, uuidErr := uuid.Parse(chi.URLParam(r, bindingIdParamName))
		if uuidErr != nil {
			paramErr := errs.NewParamInvalidError(bindingIdParamName, uuidErr.Error())
			util.Respond(w, paramErr, util.HttpStatusCodeFromAppError(paramErr))
			return
		}

		exchangeName := chi.URLParam(r, "exchangeName")
		exchange, exchangeErr := exchangeRepository.GetExchange(exchangeName)
		if exchangeErr != nil {
			util.Respond(w, exchangeErr, util.HttpStatusCodeFromAppError(exchangeErr))
			return
		}

		bindingErr := exchange.Unbind(bindingId)
		if bindingErr != nil {
			util.Respond(w, bindingErr, util.HttpStatusCodeFromAppError(bindingErr))
			return
		}

		util.Respond(w, nil, http.StatusAccepted)
	}
}
