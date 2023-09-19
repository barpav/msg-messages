package data

import "context"

type querySetMessageReadState struct{}

func (q querySetMessageReadState) text() string {
	return `
	WITH update_constraints AS (
		SELECT
			COALESCE(is_deleted, false) AS message_deleted,
			event_timestamp = $1 AS timestamp_match,
			COALESCE(is_read, false) != $2 AS state_modified
		FROM messages
		WHERE id = $3
		FOR UPDATE
	),
	update_try AS (
		UPDATE messages SET
			event_timestamp = nextval('timeline'),
			is_read = NULLIF($2, false),
			edited = $4
		WHERE id = $3
			AND COALESCE(is_deleted, false) = false
			AND event_timestamp = $1
			AND COALESCE(is_read, false) != $2
		RETURNING
			event_timestamp AS new_timestamp,
			sender AS sender,
			receiver AS receiver
	)
	SELECT
		COALESCE(update_try.new_timestamp, 0) AS new_timestamp,
		update_constraints.message_deleted AS message_deleted,
		update_constraints.timestamp_match AS timestamp_match,
		update_constraints.state_modified AS state_modified,
		COALESCE(update_try.sender, '') AS sender,
		COALESCE(update_try.receiver, '') AS receiver
	FROM update_constraints
		LEFT OUTER JOIN update_try
		ON true;
	`
}

func (s *Storage) SetMessageReadState(ctx context.Context, id, timestamp int64, read bool) (newTimestamp int64, err error) {
	newTimestamp, err = s.modifyMessage(ctx, id, timestamp, querySetMessageReadState{}, read)
	return newTimestamp, err
}
