/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/melyouz/risala/broker/internal"
	"github.com/melyouz/risala/broker/internal/errs"
	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
)

func HandleQueueMessageNack(queueRepository storage.QueueRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queueName := chi.URLParam(r, "queueName")
		messageIdParamName := "messageId"
		messageId, uuidErr := uuid.Parse(chi.URLParam(r, messageIdParamName))
		if uuidErr != nil {
			paramErr := errs.NewParamInvalidError(messageIdParamName, uuidErr.Error())
			util.Respond(w, paramErr, util.HttpStatusCodeFromAppError(paramErr))
			return
		}

		queue, queueErr := queueRepository.GetQueue(queueName)
		if queueErr != nil {
			util.Respond(w, queueErr, util.HttpStatusCodeFromAppError(queueErr))
			return
		}

		message, nackErr := queue.Nack(messageId)
		if nackErr != nil {
			util.Respond(w, nackErr, util.HttpStatusCodeFromAppError(nackErr))
			return
		}

		deadLetterQueue, deadLetterQueueErr := queueRepository.GetQueue(internal.DeadLetterQueueName)
		if deadLetterQueueErr != nil {
			util.Respond(w, deadLetterQueueErr, util.HttpStatusCodeFromAppError(deadLetterQueueErr))
			return
		}

		enqueueErr := deadLetterQueue.Enqueue(message)
		if enqueueErr != nil {
			util.Respond(w, enqueueErr, util.HttpStatusCodeFromAppError(enqueueErr))
			return
		}

		util.Respond(w, nil, http.StatusNoContent)
	}
}
