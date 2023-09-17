package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/barpav/msg-messages/internal/rest/models"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Service struct {
	Shutdown  chan struct{}
	cfg       *config
	server    *http.Server
	auth      Authenticator
	storage   Storage
	fileStats FileStats
}

type Authenticator interface {
	ValidateSession(ctx context.Context, key, ip, agent string) (userId string, err error)
}

type Storage interface {
	CreateNewPersonalMessageV1(ctx context.Context, sender string, data *models.NewPersonalMessageV1) (id int64, timestamp int64, err error)
	MessageUpdatesV1(ctx context.Context, userId string, after int64, limit int) (*models.MessageUpdatesV1, error)
	PersonalMessageV1(ctx context.Context, userId string, messageId int64) (*models.PersonalMessageV1, error)
}

type FileStats interface {
	SendUsage(ctx context.Context, fileId string, inUse bool) error
}

func (s *Service) Start(auth Authenticator, storage Storage, fileStats FileStats) {
	s.cfg = &config{}
	s.cfg.Read()

	s.auth, s.storage, s.fileStats = auth, storage, fileStats

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.port),
		Handler: s.operations(),
	}

	s.Shutdown = make(chan struct{}, 1)

	go func() {
		err := s.server.ListenAndServe()

		if err != http.ErrServerClosed {
			log.Err(err).Msg("HTTP server crashed.")
		}

		s.Shutdown <- struct{}{}
	}()
}

func (s *Service) Stop(ctx context.Context) (err error) {
	err = s.server.Shutdown(ctx)

	if err != nil {
		err = fmt.Errorf("failed to stop HTTP service: %w", err)
	}

	return err
}

// Specification: https://barpav.github.io/msg-api-spec/#/messages
func (s *Service) operations() *chi.Mux {
	ops := chi.NewRouter()

	ops.Use(s.traceInternalServerError)
	ops.Use(s.authenticate)

	// Public endpoint is the concern of the api gateway
	ops.Post("/", s.sendNewMessage)
	ops.Get("/", s.syncMessages)
	ops.Get("/{id}", s.getMessageData)
	ops.Put("/{id}", s.editMessage)
	ops.Patch("/{id}", s.markMessageAsRead)
	ops.Delete("/{id}", s.deleteMessageData)

	return ops
}
