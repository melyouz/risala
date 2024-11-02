/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package util

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/melyouz/risala/broker/internal/errs"
)

const ApiV1BasePath = "/api/v1"

func AssertOk(t *testing.T, response *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, response.Code)
}

func AssertAccepted(t *testing.T, response *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusAccepted, response.Code)
}

func AssertCreated(t *testing.T, response *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusCreated, response.Code)
}

func AssertNotFound(t *testing.T, response *httptest.ResponseRecorder, expectedErrorCode string, expectedErrorMessage string) {
	assert.Equal(t, http.StatusNotFound, response.Code)
	jsonResponse := JSONItemResponse(response)
	assert.Equal(t, expectedErrorCode, jsonResponse["code"])
	assert.Equal(t, expectedErrorMessage, jsonResponse["message"])
}

func AssertValidationErrors(t *testing.T, response *httptest.ResponseRecorder, expectedErrors []errs.ValidationError) {
	assert.Equal(t, http.StatusBadRequest, response.Code)
	jsonResponse := JSONItemResponse(response)
	assert.Equal(t, "VALIDATION_ERROR", jsonResponse["code"])
	assert.Len(t, jsonResponse["errors"], len(expectedErrors))

	errors := jsonResponse["errors"].([]interface{})
	for i, expectedErr := range expectedErrors {
		assert.Equal(t, expectedErr.Field, errors[i].(map[string]interface{})["field"])
		assert.Equal(t, expectedErr.Message, errors[i].(map[string]interface{})["message"])
	}
}

func AssertConflict(t *testing.T, response *httptest.ResponseRecorder, expectedErrorCode string, expectedErrorMessage string) {
	assert.Equal(t, http.StatusConflict, response.Code)
	jsonResponse := JSONItemResponse(response)
	assert.Equal(t, expectedErrorCode, jsonResponse["code"])
	assert.Equal(t, expectedErrorMessage, jsonResponse["message"])
}

func JSONCollectionResponse(response *httptest.ResponseRecorder) (jsonResponse []map[string]interface{}) {
	_ = json.Unmarshal([]byte(response.Body.String()), &jsonResponse)

	return jsonResponse
}

func JSONItemResponse(response *httptest.ResponseRecorder) (jsonResponse map[string]interface{}) {
	_ = json.Unmarshal([]byte(response.Body.String()), &jsonResponse)

	return jsonResponse
}
