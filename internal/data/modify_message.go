package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type ErrMessageDeleted struct{}
type ErrTimestampIsNotMatch struct{}
type ErrMessageNotModified struct{}

func (s *Storage) modifyMessage(ctx context.Context, id, timestamp int64, q query, arg any) (newTimestamp int64, err error) {
	var (
		tx                                       *sql.Tx
		messageDeleted, timestampMatch, modified bool
		sender, receiver                         string
	)

	tx, err = s.db.BeginTx(ctx, nil)

	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	err = tx.Stmt(s.queries[q]).QueryRowContext(ctx, timestamp, arg, id, time.Now().UTC()).Scan(
		&newTimestamp, &messageDeleted, &timestampMatch, &modified, &sender, &receiver,
	)

	switch {
	case err != nil:
		if err == sql.ErrNoRows {
			return 0, errors.New("message not found")
		}
		return 0, err
	case messageDeleted:
		return 0, &ErrMessageDeleted{}
	case !timestampMatch:
		return 0, &ErrTimestampIsNotMatch{}
	case !modified:
		return 0, &ErrMessageNotModified{}
	case newTimestamp == 0:
		return 0, errors.New("failed to modify message")
	}

	_, err = tx.Stmt(s.queries[queryWriteUpdate{}]).ExecContext(ctx, sender, newTimestamp, id)

	if err != nil {
		return 0, fmt.Errorf("failed to write user '%s' update: %w", sender, err)
	}

	_, err = tx.Stmt(s.queries[queryWriteUpdate{}]).ExecContext(ctx, receiver, newTimestamp, id)

	if err != nil {
		return 0, fmt.Errorf("failed to write user '%s' update: %w", receiver, err)
	}

	err = tx.Commit()

	if err != nil {
		return 0, err
	}

	return newTimestamp, nil
}

func (e *ErrMessageDeleted) Error() string {
	return "message deleted"
}

func (e *ErrMessageDeleted) ImplementsMessageDeletedError() {
}

func (e *ErrTimestampIsNotMatch) Error() string {
	return "message timestamp is not match"
}

func (e *ErrTimestampIsNotMatch) ImplementsTimestampIsNotMatchError() {
}

func (e *ErrMessageNotModified) Error() string {
	return "message has not been modified"
}

func (e *ErrMessageNotModified) ImplementsMessageNotModifiedError() {
}
