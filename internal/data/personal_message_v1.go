package data

import (
	"context"
	"database/sql"

	"github.com/barpav/msg-messages/internal/rest/models"
)

type queryGetPersonalMessageV1 struct{}

func (q queryGetPersonalMessageV1) text() string {
	return `
	SELECT
		event_timestamp,
		CASE WHEN COALESCE(is_deleted, false) THEN '' ELSE sender END,
		CASE WHEN COALESCE(is_deleted, false) THEN '' ELSE receiver END,
		created,
		edited,
		COALESCE(is_read, false),
		COALESCE(message_text, ''),
		COALESCE(is_deleted, false)
	FROM messages
	WHERE id = $1
		AND (sender = $2 OR receiver = $2);
	`
}

type queryGetPersonalMessageAttachmentsV1 struct{}

func (q queryGetPersonalMessageAttachmentsV1) text() string {
	return `
	SELECT file_id
	FROM attachments
	WHERE message_id = $1;
	`
}

func (s *Storage) PersonalMessageV1(ctx context.Context, userId string, messageId int64) (*models.PersonalMessageV1, error) {
	row := s.queries[queryGetPersonalMessageV1{}].QueryRowContext(ctx, messageId, userId)
	err := row.Err()

	if err != nil {
		return nil, err
	}

	message := &models.PersonalMessageV1{Id: messageId, Files: make([]string, 0)}
	err = row.Scan(
		&message.Timestamp,
		&message.From,
		&message.To,
		&message.Created,
		&message.Edited,
		&message.Read,
		&message.Text,
		&message.Deleted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	var rows *sql.Rows
	rows, err = s.queries[queryGetPersonalMessageAttachmentsV1{}].QueryContext(ctx, messageId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var fileId string
	for rows.Next() {
		err = rows.Scan(&fileId)

		if err != nil {
			return nil, err
		}

		message.Files = append(message.Files, fileId)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return message, nil
}
