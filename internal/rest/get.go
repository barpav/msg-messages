package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/barpav/msg-messages/internal/rest/models"
	"github.com/go-chi/chi/v5"
)

const mimeTypePersonalMessageV1 = "application/vnd.personalMessage.v1+json"

// https://barpav.github.io/msg-api-spec/#/messages/get_messages__id_
func (s *Service) getMessageData(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var message *models.PersonalMessageV1
	message, err = s.storage.PersonalMessageV1(r.Context(), authenticatedUser(r), id)

	if err == nil {
		if message == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", mimeTypePersonalMessageV1)
		err = json.NewEncoder(w).Encode(message)
	}

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to get message data (v1)")
		return
	}

	w.WriteHeader(http.StatusOK)
}
