package models

import (
	"encoding/json"
	"errors"
	"io"
)

// Schema: messageReadMark.v1
type MessageReadMarkV1 struct {
	Read bool
}

func (m *MessageReadMarkV1) Deserialize(data io.Reader) error {
	if json.NewDecoder(data).Decode(m) != nil {
		return errors.New("New message data violates 'messageReadMark.v1' schema.")
	}
	return nil
}
