package metabase

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

// Client provides methods to interact with the Metabase API.
type Client struct {
	httpClient  *resty.Client
	baseURL     string
	apiKey      string
	sessionAuth *sessionAuth
	logger      zerolog.Logger
}

// NewClient creates a new Metabase API client.
func NewClient(baseURL, apiKey, username, password string, logger zerolog.Logger) (*Client, error) {
	c := &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		logger:  logger,
	}

	httpClient := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(30 * time.Second)

	httpClient.OnBeforeRequest(func(_ *resty.Client, r *resty.Request) error {
		logger.Debug().
			Str("method", r.Method).
			Str("url", r.URL).
			Msg("metabase API request")
		return nil
	})

	httpClient.OnAfterResponse(func(_ *resty.Client, r *resty.Response) error {
		logger.Debug().
			Str("method", r.Request.Method).
			Str("url", r.Request.URL).
			Int("status", r.StatusCode()).
			Dur("duration", r.Time()).
			Msg("metabase API response")
		return nil
	})

	if apiKey != "" {
		httpClient.SetHeader("x-api-key", apiKey)
	} else if username != "" && password != "" {
		sa := newSessionAuth(baseURL, username, password, logger)
		c.sessionAuth = sa

		if err := sa.authenticate(); err != nil {
			return nil, fmt.Errorf("initial authentication failed: %w", err)
		}

		httpClient.OnBeforeRequest(func(_ *resty.Client, r *resty.Request) error {
			r.SetHeader("X-Metabase-Session", sa.getSessionID())
			return nil
		})

		httpClient.AddRetryCondition(func(r *resty.Response, _ error) bool {
			if r != nil && r.StatusCode() == 401 {
				logger.Warn().Msg("received 401, re-authenticating")
				if err := sa.authenticate(); err != nil {
					logger.Error().Err(err).Msg("re-authentication failed")
					return false
				}
				return true
			}
			return false
		})
		httpClient.SetRetryCount(1)
	} else {
		return nil, fmt.Errorf("either API key or username/password must be provided")
	}

	c.httpClient = httpClient

	if err := c.HealthCheck(); err != nil {
		return nil, fmt.Errorf("metabase health check failed: %w", err)
	}

	return c, nil
}

// HealthCheck validates that Metabase is reachable and credentials are valid.
func (c *Client) HealthCheck() error {
	_, err := c.GetCurrentUser()
	return err
}

// checkResponse checks if the API response indicates an error.
func checkResponse(resp *resty.Response) error {
	if resp.IsError() {
		return fmt.Errorf("metabase API error (status %d): %s", resp.StatusCode(), resp.String())
	}
	return nil
}
