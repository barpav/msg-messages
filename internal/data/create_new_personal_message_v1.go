package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/barpav/msg-messages/internal/rest/models"
)

type queryCreateMessage struct{}

func (q queryCreateMessage) text() string {
	return `
	INSERT INTO messages (sender, receiver,	created, message_text)
	VALUES ($1, $2, $3, NULLIF($4, ''))
	RETURNING id, event_timestamp;
	`
}

type queryCreateAttachment struct{}

func (q queryCreateAttachment) text() string {
	return `
	INSERT INTO attachments (message_id, file_id)
	VALUES ($1, $2);
	`
}

type queryWriteUpdate struct{}

func (q queryWriteUpdate) text() string {
	return `
	INSERT INTO updates (user_id, event_timestamp, message_id)
	VALUES ($1, $2, $3);
	`
}

func (s *Storage) CreateNewPersonalMessageV1(ctx context.Context, sender string, message *models.NewPersonalMessageV1) (id int64, timestamp int64, err error) {
	var tx *sql.Tx
	tx, err = s.db.BeginTx(ctx, nil)

	if err != nil {
		return 0, 0, err
	}

	defer tx.Rollback()

	row := tx.Stmt(s.queries[queryCreateMessage{}]).QueryRowContext(ctx, sender, message.To, time.Now().UTC(), message.Text)
	err = row.Scan(&id, &timestamp)

	if err != nil {
		return 0, 0, fmt.Errorf("failed to create new message: %w", err)
	}

	if len(message.Files) != 0 {
		for _, fileId := range message.Files {
			_, err = tx.Stmt(s.queries[queryCreateAttachment{}]).ExecContext(ctx, id, fileId)

			if err != nil {
				return 0, 0, fmt.Errorf("failed to create attachment with file id '%s': %w", fileId, err)
			}
		}
	}

	_, err = tx.Stmt(s.queries[queryWriteUpdate{}]).ExecContext(ctx, sender, timestamp, id)

	if err != nil {
		return 0, 0, fmt.Errorf("failed to write user '%s' update: %w", sender, err)
	}

	_, err = tx.Stmt(s.queries[queryWriteUpdate{}]).ExecContext(ctx, message.To, timestamp, id)

	if err != nil {
		return 0, 0, fmt.Errorf("failed to write user '%s' update: %w", message.To, err)
	}

	err = tx.Commit()

	if err != nil {
		return 0, 0, err
	}

	return id, timestamp, nil
}
