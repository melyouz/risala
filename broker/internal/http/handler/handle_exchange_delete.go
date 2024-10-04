/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
	"net/http"
)

func HandleExchangeDelete(exchangeRepository storage.ExchangeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exchangeName := chi.URLParam(r, "exchangeName")
		err := exchangeRepository.DeleteExchange(exchangeName)
		if err != nil {
			util.Respond(w, err, util.HttpStatusCodeFromAppError(err))
			return
		}
		util.Respond(w, nil, http.StatusAccepted)
	}
}
