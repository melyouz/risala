/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
)

func HandleQueueMessagePeek(queueRepository storage.QueueRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 1 {
			limit = 1
		}

		queueName := chi.URLParam(r, "queueName")
		queue, queueErr := queueRepository.GetQueue(queueName)
		if queueErr != nil {
			util.Respond(w, queueErr, util.HttpStatusCodeFromAppError(queueErr))
			return
		}

		messages, err := queue.Peek(limit)
		if err != nil {
			util.Respond(w, err, util.HttpStatusCodeFromAppError(err))
			return
		}

		util.Respond(w, messages, http.StatusOK)
	}
}
