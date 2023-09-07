package rest

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

func (s *Service) traceInternalServerError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err := fmt.Errorf("recovered from panic: %v", rec)
				logAndReturnErrorWithIssue(w, r, err, "Internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

type authenticatedUserId struct{}

func (s *Service) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userId string
		var err error
		sessionKey := sessionKey(r)

		if sessionKey != "" {
			userId, err = s.auth.ValidateSession(r.Context(), sessionKey, userIP(r), userAgent(r))
		}

		if err != nil {
			logAndReturnErrorWithIssue(w, r, err, "Session key authentication failed")
			return
		}

		if userId == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authenticatedUserId{}, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func logAndReturnErrorWithIssue(w http.ResponseWriter, r *http.Request, err error, logMsg string) {
	issue := requestId(r)

	log.Err(err).Msg(fmt.Sprintf("Issue %s: %s", issue, logMsg))

	w.Header()["issue"] = []string{issue} // lowercase - non-canonical (vendor) header
	w.WriteHeader(http.StatusInternalServerError)
}

func requestId(r *http.Request) string {
	id := r.Header.Get("request-id") // set by api-gateway

	if id != "" {
		return id
	}

	return "untraced"
}

func sessionKey(r *http.Request) string {
	parts := strings.Split(r.Header.Get("Authorization"), "Bearer")

	if len(parts) != 2 {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

func authenticatedUser(r *http.Request) (id string) {
	id, _ = r.Context().Value(authenticatedUserId{}).(string)
	return id
}

func userIP(r *http.Request) string {
	ip := r.Header.Get("remote-addr") // set by api-gateway

	if ip != "" {
		return ip
	}

	return "unknown"
}

func userAgent(r *http.Request) string {
	agent := r.Header.Get("User-Agent")

	if agent != "" {
		return agent
	}

	return "unknown"
}
