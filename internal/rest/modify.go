package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/barpav/msg-messages/internal/rest/models"
	"github.com/go-chi/chi/v5"
)

const mimeTypeEditedMessageTextV1 = "application/vnd.editedMessageText.v1+json"
const mimeTypeMessageReadMarkV1 = "application/vnd.messageReadMark.v1+json"

type ErrTimestampIsNotMatch interface {
	Error() string
	ImplementsTimestampIsNotMatchError()
}

type ErrMessageNotModified interface {
	Error() string
	ImplementsMessageNotModifiedError()
}

type ErrMessageDeleted interface {
	Error() string
	ImplementsMessageDeletedError()
}

// https://barpav.github.io/msg-api-spec/#/messages/patch_messages__id_
func (s *Service) modifyMessage(w http.ResponseWriter, r *http.Request) {
	mimeType := r.Header.Get("Content-Type")

	if mimeType != mimeTypeEditedMessageTextV1 && mimeType != mimeTypeMessageReadMarkV1 {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

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
	userId := authenticatedUser(r)
	var message *models.PersonalMessageV1

	message, err = s.storage.PersonalMessageV1(ctx, userId, id)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to receive personal message data (v1)")
		return
	}

	if message == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var newTimestamp int64

	switch mimeType {
	case mimeTypeEditedMessageTextV1:
		editedData := models.EditedMessageTextV1{}
		err = editedData.Deserialize(r.Body)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		if userId != message.From {
			http.Error(w, "Only sender of the message can edit its text.", 400)
			return
		}

		if editedData.Text == "" && len(message.Files) == 0 {
			http.Error(w, "Text in a message without attachments cannot be empty.", 400)
			return
		}

		newTimestamp, err = s.storage.EditMessageText(ctx, id, clientTimestamp, editedData.Text)
	case mimeTypeMessageReadMarkV1:
		editedData := models.MessageReadMarkV1{}
		err = editedData.Deserialize(r.Body)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		if userId != message.To {
			http.Error(w, "Only receiver of the message can mark it as read.", 400)
			return
		}

		newTimestamp, err = s.storage.SetMessageReadState(ctx, id, clientTimestamp, editedData.Read)
	}

	if err != nil {
		if _, ok := err.(ErrTimestampIsNotMatch); ok {
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}

		if _, ok := err.(ErrMessageNotModified); ok {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		if _, ok := err.(ErrMessageDeleted); ok {
			w.WriteHeader(http.StatusGone)
			return
		}

		logAndReturnErrorWithIssue(w, r, err, "Failed to modify message")
		return
	}

	w.Header()["ETag"] = []string{fmt.Sprintf("%d", newTimestamp)}
	w.WriteHeader(http.StatusOK)
}
