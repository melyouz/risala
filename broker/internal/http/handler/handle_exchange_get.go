/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
)

func HandleExchangeGet(exchangeRepository storage.ExchangeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exchangeName := chi.URLParam(r, "exchangeName")
		exchange, err := exchangeRepository.GetExchange(exchangeName)
		if err != nil {
			util.Respond(w, err, util.HttpStatusCodeFromAppError(err))
			return
		}

		util.Respond(w, exchange, http.StatusOK)
	}
}
