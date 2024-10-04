/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
	"net/http"
)

func HandleQueueGet(queueRepository storage.QueueRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queueName := chi.URLParam(r, "queueName")
		queue, err := queueRepository.GetQueue(queueName)
		if err != nil {
			util.Respond(w, err, util.HttpStatusCodeFromAppError(err))
			return
		}

		util.Respond(w, queue, http.StatusOK)
	}
}
