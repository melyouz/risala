package handler

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
	"net/http"
)

func HandleExchangeCreate(exchangeRepository storage.ExchangeRepository, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var exchange internal.Exchange
		util.Decode(r, &exchange)

		var vErrors validator.ValidationErrors
		if errors.As(validate.Struct(exchange), &vErrors) {
			util.Respond(w, errs.NewValidationError(vErrors), http.StatusBadRequest)
			return
		}

		if exchange.Bindings == nil {
			exchange.Bindings = []internal.Binding{}
		}

		existingExchange, _ := exchangeRepository.GetExchange(exchange.Name)
		if existingExchange != nil {
			existsErr := errs.NewExchangeExistsError(fmt.Sprintf("Exchange '%s' already exists", exchange.Name))
			util.Respond(w, existsErr, util.HttpStatusCodeFromAppError(existsErr))
			return
		}

		exchangeRepository.StoreExchange(&exchange)

		util.Respond(w, exchange, http.StatusCreated)
	}
}
