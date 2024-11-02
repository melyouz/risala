/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
)

func HandleQueueMessageConsume(queueRepository storage.QueueRepository) http.HandlerFunc {
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

		result := make([]*internal.Message, 0)

		for i := 0; i < limit; i++ {
			message := queue.Dequeue()
			if message == nil {
				break
			}

			_ = queue.Ack(message.Id)

			result = append(result, message)
		}

		util.Respond(w, result, http.StatusOK)
	}
}
