package metabase

import (
	"fmt"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

// sessionResponse represents the response from POST /api/session.
type sessionResponse struct {
	ID string `json:"id"`
}

// sessionAuth manages Metabase session-based authentication.
type sessionAuth struct {
	baseURL    string
	username   string
	password   string
	sessionID  string
	mu         sync.RWMutex
	authClient *resty.Client // separate client to avoid middleware loops
	logger     zerolog.Logger
}

func newSessionAuth(baseURL, username, password string, logger zerolog.Logger) *sessionAuth {
	authClient := resty.New().
		SetBaseURL(baseURL)

	return &sessionAuth{
		baseURL:    baseURL,
		username:   username,
		password:   password,
		authClient: authClient,
		logger:     logger,
	}
}

// authenticate creates a new Metabase session.
func (s *sessionAuth) authenticate() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Info().Msg("authenticating with Metabase session")

	var result sessionResponse
	resp, err := s.authClient.R().
		SetBody(map[string]string{
			"username": s.username,
			"password": s.password,
		}).
		SetResult(&result).
		Post("/api/session")

	if err != nil {
		return fmt.Errorf("session auth request failed: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("session auth failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	s.sessionID = result.ID
	s.logger.Info().Msg("metabase session created successfully")
	return nil
}

// getSessionID returns the current session ID.
func (s *sessionAuth) getSessionID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessionID
}
