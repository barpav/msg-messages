package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db      *sql.DB
	cfg     *config
	queries map[query]*sql.Stmt
}

type query interface {
	text() string
}

func queriesToPrepare() []query {
	return []query{
		queryCreateMessage{},
		queryCreateAttachment{},
		queryWriteUpdate{},
	}
}

func (s *Storage) Open() (err error) {
	s.cfg = &config{}
	s.cfg.Read()

	err = s.connectToDatabase()

	if err != nil {
		return err
	}

	return s.prepareQueries()
}

func (s *Storage) Close(ctx context.Context) (err error) {
	var closeErr error
	closed := make(chan struct{}, 1)

	go func() {
		for _, stmt := range s.queries {
			if stmt != nil {
				closeErr = errors.Join(err, stmt.Close())
			}
		}

		closeErr = errors.Join(err, s.db.Close())

		closed <- struct{}{}
	}()

	select {
	case <-closed:
		err = closeErr
	case <-ctx.Done():
		err = ctx.Err()
	}

	if err != nil {
		err = fmt.Errorf("failed to disconnect from database: %w", err)
	}

	return err
}

func (s *Storage) connectToDatabase() (err error) {
	dbAddress := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", s.cfg.user, s.cfg.password, s.cfg.host, s.cfg.port, s.cfg.database)

	s.db, err = sql.Open("pgx", dbAddress)

	if err == nil {
		err = s.db.Ping()
	}

	if err == nil {
		log.Info().Msg(fmt.Sprintf("Successfully connected to DB at %s", dbAddress))
	}

	return err
}

func (s *Storage) prepareQueries() (err error) {
	s.queries = make(map[query]*sql.Stmt)

	for _, q := range queriesToPrepare() {
		err = errors.Join(err, s.prepare(q))
	}

	return err
}

func (s *Storage) prepare(q query) (err error) {
	s.queries[q], err = s.db.Prepare(q.text())
	return err
}
