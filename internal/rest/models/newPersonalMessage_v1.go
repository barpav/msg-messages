package models

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

// Schema: newPersonalMessage.v1
type NewPersonalMessageV1 struct {
	To    string
	Text  string
	Files []string
}

func (m *NewPersonalMessageV1) Deserialize(data io.Reader) error {
	if json.NewDecoder(data).Decode(m) != nil {
		return errors.New("New message data violates 'newPersonalMessage.v1' schema.")
	}

	m.To = strings.TrimSpace(m.To)
	m.Text = strings.TrimSpace(m.Text)

	return m.validate()
}

func (m *NewPersonalMessageV1) validate() (err error) {
	if m.To == "" {
		err = errors.Join(err, errors.New("Message recipient must be specified."))
	}

	if m.Text == "" && len(m.Files) == 0 {
		err = errors.Join(err, errors.New("Message text or attached files must be specified."))
	}

	for _, fileId := range m.Files {
		if len(fileId) != 24 {
			err = errors.Join(err, errors.New("Attached file id must be 24 character long."))
			break
		}
	}

	return err
}
