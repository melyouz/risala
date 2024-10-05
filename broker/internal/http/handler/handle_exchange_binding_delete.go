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
		bindingIdParam := "bindingId"
		bindingId, uuidErr := uuid.Parse(chi.URLParam(r, bindingIdParam))
		if uuidErr != nil {
			paramErr := errs.NewParamInvalidError(bindingIdParam, uuidErr.Error())
			util.Respond(w, paramErr, util.HttpStatusCodeFromAppError(paramErr))
			return
		}

		exchangeName := chi.URLParam(r, "exchangeName")
		bindingErr := exchangeRepository.DeleteBinding(exchangeName, bindingId)
		if bindingErr != nil {
			util.Respond(w, bindingErr, util.HttpStatusCodeFromAppError(bindingErr))
			return
		}

		util.Respond(w, nil, http.StatusAccepted)
	}
}
