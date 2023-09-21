package data

import "context"

type queryDeleteMessageData struct{}

func (q queryDeleteMessageData) text() string {
	return `
	WITH update_constraints AS (
		SELECT
			COALESCE(is_deleted, false) AS message_deleted,
			event_timestamp = $1 AS timestamp_match
		FROM messages
		WHERE id = $2
		FOR UPDATE
	),
	update_try AS (
		UPDATE messages SET
			event_timestamp = nextval('timeline'),
			created = null,
			edited = null,
			is_read = null,
			message_text = null,
			is_deleted = true
		WHERE id = $2
			AND COALESCE(is_deleted, false) = false
			AND event_timestamp = $1
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
		true AS message_modified,
		COALESCE(update_try.sender, '') AS sender,
		COALESCE(update_try.receiver, '') AS receiver
	FROM update_constraints
		LEFT OUTER JOIN update_try
		ON true;
	`
}

func (s *Storage) DeleteMessageData(ctx context.Context, id, timestamp int64) (newTimestamp int64, err error) {
	newTimestamp, err = s.modifyMessage(ctx, queryDeleteMessageData{}, timestamp, id)
	return newTimestamp, err
}
