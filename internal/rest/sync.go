package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/barpav/msg-messages/internal/rest/models"
)

const mimeTypeMessageUpdatesV1 = "application/vnd.messageUpdates.v1+json"

// https://barpav.github.io/msg-api-spec/#/messages/get_messages
func (s *Service) syncMessages(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Accept") {
	case "", mimeTypeMessageUpdatesV1: // including if not specified
		s.getMessageUpdatesV1(w, r)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
}

func (s *Service) getMessageUpdatesV1(w http.ResponseWriter, r *http.Request) {
	after, limit, err := getMessageUpdatesV1Parameters(r)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var updates *models.MessageUpdatesV1
	updates, err = s.storage.MessageUpdatesV1(r.Context(), authenticatedUser(r), after, limit)

	if err == nil {
		w.Header().Set("Content-Type", mimeTypeMessageUpdatesV1)
		err = json.NewEncoder(w).Encode(updates)
	}

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to get message updates (v1)")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getMessageUpdatesV1Parameters(r *http.Request) (after int64, limit int, err error) {
	var param string
	param = r.URL.Query().Get("after")

	if param != "" {
		after, err = strconv.ParseInt(param, 10, 0)

		if err != nil {
			return 0, 0, errors.New("Parameter 'after' must be an integer type.")
		}
	}

	param = r.URL.Query().Get("limit")

	if param == "" {
		return after, 50, nil
	}

	limit, err = strconv.Atoi(param)

	if err != nil {
		return 0, 0, errors.New("Parameter 'limit' must be an integer type.")
	}

	const limitMin = 1
	const limitMax = 100

	if limit < limitMin || limit > limitMax {
		return 0, 0, fmt.Errorf("Invalid parameter 'limit': min %d, max %d.", limitMin, limitMax)
	}

	return after, limit, nil
}
