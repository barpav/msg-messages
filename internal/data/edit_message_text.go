package data

import (
	"context"
	"time"
)

type queryEditMessageText struct{}

func (q queryEditMessageText) text() string {
	return `
	WITH update_constraints AS (
		SELECT
			COALESCE(is_deleted, false) AS message_deleted,
			event_timestamp = $1 AS timestamp_match,
			message_text != NULLIF($2, '') AS text_modified
		FROM messages
		WHERE id = $3
		FOR UPDATE
	),
	update_try AS (
		UPDATE messages SET
			event_timestamp = nextval('timeline'),
			message_text = NULLIF($2, ''),
			edited = $4
		WHERE id = $3
			AND COALESCE(is_deleted, false) = false
			AND event_timestamp = $1
			AND message_text != NULLIF($2, '')
		RETURNING
			id AS id,
			event_timestamp AS new_timestamp,
			sender AS sender,
			receiver AS receiver
	)
	SELECT
		COALESCE(update_try.id, 0) AS id,
		COALESCE(update_try.new_timestamp, 0) AS new_timestamp,
		update_constraints.message_deleted AS message_deleted,
		update_constraints.timestamp_match AS timestamp_match,
		update_constraints.text_modified AS text_modified,
		COALESCE(update_try.sender, '') AS sender,
		COALESCE(update_try.receiver, '') AS receiver
	FROM update_constraints
		LEFT OUTER JOIN update_try
		ON true;
	`
}

func (s *Storage) EditMessageText(ctx context.Context, id, timestamp int64, text string) (newTimestamp int64, err error) {
	newTimestamp, err = s.modifyMessage(ctx, queryEditMessageText{}, timestamp, text, id, time.Now().UTC())
	return newTimestamp, err
}
