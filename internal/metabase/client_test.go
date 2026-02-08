package metabase

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	logger := zerolog.Nop()
	client, err := NewClient(server.URL, "test-api-key", "", "", logger)
	require.NoError(t, err)

	return server, client
}

func TestNewClient_APIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-key", r.Header.Get("x-api-key"))
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
