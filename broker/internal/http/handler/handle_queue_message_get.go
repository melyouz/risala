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

func HandleQueueMessageGet(queueRepository storage.QueueRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		countParam := "count"
		messagesCount, _ := strconv.Atoi(r.URL.Query().Get(countParam))
		if messagesCount < 1 {
			messagesCount = 1
		}

		queueName := chi.URLParam(r, "queueName")
		queue, queueErr := queueRepository.GetQueue(queueName)
		if queueErr != nil {
			util.Respond(w, queueErr, util.HttpStatusCodeFromAppError(queueErr))
			return
		}

		messages, err := queue.Seek(messagesCount)
		if err != nil {
			util.Respond(w, err, util.HttpStatusCodeFromAppError(err))
			return
		}

		util.Respond(w, messages, http.StatusOK)
	}
}
