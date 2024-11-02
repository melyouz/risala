/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/storage"
	"github.com/melyouz/risala/broker/internal/testing/util"
)

func setupExchangeFindTest(t *testing.T, exchanges map[string]*internal.Exchange) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	exchangeRepository := storage.NewInMemoryExchangeRepository(exchanges)
	request := httptest.NewRequest(http.MethodGet, util.ApiV1BasePath+"/exchanges", nil)
	response := httptest.NewRecorder()

	HandleExchangeFind(exchangeRepository)(response, request)

	return response, request
}

func TestHandleExchangesFind(t *testing.T) {
	t.Parallel()
	t.Run("Returns list when exchanges exist", func(t *testing.T) {
		t.Parallel()

		exchanges := map[string]*internal.Exchange{
			"app.internal": util.NewTestExchangeWithoutBindings("app.internal"),
			"app.external": util.NewTestExchangeWithoutBindings("app.external"),
		}

		response, _ := setupExchangeFindTest(t, exchanges)

		util.AssertOk(t, response)
		jsonResponse := util.JSONCollectionResponse(response)
		assert.Len(t, jsonResponse, 2)
		assert.Equal(t, "app.external", jsonResponse[0]["name"])
		assert.Len(t, jsonResponse[0]["bindings"], 0)
		assert.Equal(t, "app.internal", jsonResponse[1]["name"])
		assert.Len(t, jsonResponse[1]["bindings"], 0)
	})

	t.Run("Returns empty list when no exchanges", func(t *testing.T) {
		t.Parallel()

		exchanges := map[string]*internal.Exchange{}

		response, _ := setupExchangeFindTest(t, exchanges)

		util.AssertOk(t, response)
		jsonResponse := util.JSONCollectionResponse(response)
		assert.Len(t, jsonResponse, 0)
	})
}
