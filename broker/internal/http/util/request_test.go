/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	type TestSubEntity struct {
		Content string `json:"content" validate:"required"`
	}

	type TestEntity struct {
		Name     string          `json:"name" validate:"required"`
		Children []TestSubEntity `json:"children" validate:"dive"`
	}

	t.Run("Decodes JSON body from HTTP request", func(t *testing.T) {
		requestBody, _ := json.Marshal(map[string]interface{}{
			"name": "testEntityName",
			"children": []TestSubEntity{
				{Content: "Test content 1"},
				{Content: "Test content 2"},
			},
		})
		request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))

		var entity TestEntity
		Decode(request, &entity)

		assert.Equal(t, "testEntityName", entity.Name)
		assert.Len(t, entity.Children, 2)
		assert.Equal(t, "Test content 1", entity.Children[0].Content)
		assert.Equal(t, "Test content 2", entity.Children[1].Content)
	})
}
