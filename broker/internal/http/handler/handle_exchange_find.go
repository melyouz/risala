package handler

import (
	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
	"net/http"
)

func HandleExchangeFind(exchangeRepository storage.ExchangeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exchangesList := exchangeRepository.FindExchanges()

		util.Respond(w, exchangesList, http.StatusOK)
	}
}
