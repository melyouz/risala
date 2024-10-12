/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespond(t *testing.T) {
	type TestSubEntity struct {
		Content string `json:"content" validate:"required"`
	}

	type TestEntity struct {
		Name     string          `json:"name" validate:"required"`
		Children []TestSubEntity `json:"children" validate:"dive"`
	}

	t.Run("Responds to HTTP request", func(t *testing.T) {
		entity := TestEntity{
			Name: "testEntityName",
			Children: []TestSubEntity{
				{Content: "Test content 1"},
				{Content: "Test content 2"},
			},
		}

		response := httptest.NewRecorder()
		Respond(response, entity, http.StatusOK)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, "{\"name\":\"testEntityName\",\"children\":[{\"content\":\"Test content 1\"},{\"content\":\"Test content 2\"}]}", response.Body.String())
	})
}
