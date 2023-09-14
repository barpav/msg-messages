package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/barpav/msg-messages/internal/rest/models"
	"github.com/rs/zerolog/log"
)

// https://barpav.github.io/msg-api-spec/#/messages/post_messages
func (s *Service) sendNewMessage(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Content-Type") {
	case "application/vnd.newPersonalMessage.v1+json":
		s.sendPersonalMessageV1(w, r)
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
}

func (s *Service) sendPersonalMessageV1(w http.ResponseWriter, r *http.Request) {
	message := models.NewPersonalMessageV1{}
	err := message.Deserialize(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var id, timestamp int64
	id, timestamp, err = s.storage.CreateNewPersonalMessageV1(r.Context(), authenticatedUser(r), &message)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to send new personal message (v1)")
		return
	}

	if len(message.Files) != 0 {
		go func() {
			ctx := context.Background()

			for _, fileId := range message.Files {
				err = s.fileStats.SendUsage(ctx, fileId, true)

				if err != nil {
					log.Err(err).Msg(fmt.Sprintf("Failed to send used file '%s' statistics.", fileId))
				}
			}
		}()
	}

	w.Header().Set("Location", fmt.Sprintf("/%d", id))
	w.Header().Set("ETag", fmt.Sprintf("%d", timestamp))
	w.WriteHeader(http.StatusCreated)
}
