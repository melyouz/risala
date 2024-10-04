/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
	"net/http"
)

func HandleQueueFind(queueRepository storage.QueueRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queuesList := queueRepository.FindQueues()

		util.Respond(w, queuesList, http.StatusOK)
	}
}
