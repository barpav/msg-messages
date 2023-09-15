package data

import (
	"context"

	"github.com/barpav/msg-messages/internal/rest/models"
)

type queryGetMessageUpdatesV1 struct{}

func (q queryGetMessageUpdatesV1) text() string {
	return `
	SELECT
		event_timestamp,
		message_id
	FROM updates
	WHERE user_id = $1 AND event_timestamp > $2
	ORDER BY event_timestamp ASC
	LIMIT $3;
	`
}

func (s *Storage) GetMessageUpdatesV1(ctx context.Context, userId string, after int64, limit int) (*models.MessageUpdatesV1, error) {
	rows, err := s.queries[queryGetMessageUpdatesV1{}].QueryContext(ctx, userId, after, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	updates := &models.MessageUpdatesV1{Messages: make([]*models.MessageUpdateInfoV1, 0, limit)}

	for rows.Next() {
		info := &models.MessageUpdateInfoV1{}
		err = rows.Scan(&info.Id, &info.Timestamp)

		if err != nil {
			return nil, err
		}

		updates.Messages = append(updates.Messages, info)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	updates.Total = len(updates.Messages)

	return updates, nil
}
