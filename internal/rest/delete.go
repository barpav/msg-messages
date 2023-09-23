package rest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/barpav/msg-messages/internal/rest/models"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// https://barpav.github.io/msg-api-spec/#/messages/delete_messages__id_
func (s *Service) deleteMessageData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var clientTimestamp int64
	clientTimestamp, err = strconv.ParseInt(r.Header.Get("If-Match"), 10, 0)

	if err != nil {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	ctx := r.Context()
	var message *models.PersonalMessageV1

	message, err = s.storage.PersonalMessageV1(ctx, authenticatedUser(r), id)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to receive personal message data (v1)")
		return
	}

	if message == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var newTimestamp int64
	newTimestamp, err = s.storage.DeleteMessageData(ctx, id, clientTimestamp)

	if err != nil {
		if _, ok := err.(ErrTimestampIsNotMatch); ok {
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}

		if _, ok := err.(ErrMessageDeleted); ok {
			w.WriteHeader(http.StatusGone)
			return
		}

		logAndReturnErrorWithIssue(w, r, err, "Failed to delete message data")
		return
	}

	if len(message.Files) != 0 {
		go func() {
			ctx := context.Background()

			for _, fileId := range message.Files {
				err = s.fileStats.SendUsage(ctx, fileId, false)

				if err != nil {
					log.Err(err).Msg(fmt.Sprintf("Failed to send unused file '%s' statistics.", fileId))
				}
			}
		}()
	}

	w.Header()["ETag"] = []string{fmt.Sprintf("%d", newTimestamp)}
}
