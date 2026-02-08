package metabase

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/user/current" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(User{ID: 1, Email: "test@test.com"})
			return
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)

	logger := zerolog.Nop()
	client, err := NewClient(server.URL, "test-api-key", "", "", logger)
	require.NoError(t, err)

	return server, client
}

func TestNewClient_APIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-key", r.Header.Get("x-api-key"))
		if r.URL.Path == "/api/user/current" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(User{ID: 1, Email: "test@test.com"})
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger := zerolog.Nop()
	client, err := NewClient(server.URL, "test-key", "", "", logger)
	require.NoError(t, err)

	resp, err := client.httpClient.R().Get("/api/database")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestNewClient_Session(t *testing.T) {
	sessionCreated := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/session" {
			sessionCreated = true
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(sessionResponse{ID: "sess-123"})
			require.NoError(t, err)
			return
		}
		if r.URL.Path == "/api/user/current" {
			assert.Equal(t, "sess-123", r.Header.Get("X-Metabase-Session"))
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(User{ID: 1, Email: "test@test.com"})
			return
		}
		assert.Equal(t, "sess-123", r.Header.Get("X-Metabase-Session"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger := zerolog.Nop()
	client, err := NewClient(server.URL, "", "admin@test.com", "pass", logger)
	require.NoError(t, err)
	assert.True(t, sessionCreated)

	resp, err := client.httpClient.R().Get("/api/database")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestNewClient_NoAuth(t *testing.T) {
	logger := zerolog.Nop()
	_, err := NewClient("http://localhost:3000", "", "", "", logger)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "either API key or username/password")
}

func TestHealthCheck_Success(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	err := client.HealthCheck()
	require.NoError(t, err)
}

func TestHealthCheck_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/user/current" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger := zerolog.Nop()
	// Build client manually without NewClient to skip the startup health check
	client := &Client{
		baseURL: server.URL,
		apiKey:  "test-api-key",
		logger:  logger,
	}
	client.httpClient = resty.New().SetBaseURL(server.URL).SetHeader("x-api-key", "test-api-key")

	err := client.HealthCheck()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestNewClient_HealthCheckFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/user/current" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"Unauthenticated"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger := zerolog.Nop()
	_, err := NewClient(server.URL, "bad-api-key", "", "", logger)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "metabase health check failed")
}
