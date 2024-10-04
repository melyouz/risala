/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"net/http"

	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
)

func HandleExchangeFind(exchangeRepository storage.ExchangeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exchangesList := exchangeRepository.FindExchanges()

		util.Respond(w, exchangesList, http.StatusOK)
	}
}
