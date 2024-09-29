package handler

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
	"net/http"
)

func HandleExchangeBindingAdd(exchangeRepository storage.ExchangeRepository, queueRepository storage.QueueRepository, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		exchange, exchangeErr := exchangeRepository.GetExchange(name)
		if exchangeErr != nil {
			util.Respond(w, exchangeErr, util.HttpStatusCodeFromAppError(exchangeErr))
			return
		}

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

		bindingErr := validateBindingDoesNotExist(exchange, binding)
		if bindingErr != nil {
			util.Respond(w, bindingErr, util.HttpStatusCodeFromAppError(bindingErr))
			return
		}

		exchange.AddBinding(binding)
		exchangeRepository.StoreExchange(exchange)

		util.Respond(w, binding, http.StatusCreated)
	}
}

func validateBindingDoesNotExist(exchange *internal.Exchange, binding internal.Binding) errs.AppError {
	for _, v := range exchange.Bindings {
		if v.Queue == binding.Queue {
			return errs.NewBindingExistsError(fmt.Sprintf("Binding to Queue '%s' already exists", binding.Queue))
		}
	}

	return nil
}
