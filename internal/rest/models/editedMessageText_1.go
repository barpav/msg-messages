package models

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

// Schema: editedMessageText.v1
type EditedMessageTextV1 struct {
	Text string
}

func (m *EditedMessageTextV1) Deserialize(data io.Reader) error {
	if json.NewDecoder(data).Decode(m) != nil {
		return errors.New("New message data violates 'editedMessageText.v1' schema.")
	}

	m.Text = strings.TrimSpace(m.Text)

	return nil
}
