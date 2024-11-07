/*
 * Copyright (c) 2024 Mohammadi El Youzghi and contributors.
 */

package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/melyouz/risala/broker/internal/errs"
	"github.com/melyouz/risala/broker/internal/http/util"
	"github.com/melyouz/risala/broker/internal/storage"
)

func HandleQueueMessageAck(queueRepository storage.QueueRepository) http.HandlerFunc {
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

		ackErr := queue.Ack(messageId)
		if ackErr != nil {
			util.Respond(w, ackErr, util.HttpStatusCodeFromAppError(ackErr))
			return
		}

		util.Respond(w, nil, http.StatusAccepted)
	}
}
