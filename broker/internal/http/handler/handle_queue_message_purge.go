/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
)

func HandleQueueMessagePurge(queueRepository storage.QueueRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		queueName := chi.URLParam(r, "queueName")
		queue, queueErr := queueRepository.GetQueue(queueName)
		if queueErr != nil {
			util.Respond(w, queueErr, util.HttpStatusCodeFromAppError(queueErr))
			return
		}

		purgeErr := queue.Purge()
		if purgeErr != nil {
			util.Respond(w, queueErr, util.HttpStatusCodeFromAppError(purgeErr))
			return
		}

		util.Respond(w, nil, http.StatusNoContent)
	}
}
