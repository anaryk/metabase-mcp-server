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

func TestSessionAuth_Authenticate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/session" && r.Method == http.MethodPost {
			var body map[string]string
			err := json.NewDecoder(r.Body).Decode(&body)
			require.NoError(t, err)
			assert.Equal(t, "admin@test.com", body["username"])
			assert.Equal(t, "password123", body["password"])

			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(sessionResponse{ID: "test-session-id"})
			require.NoError(t, err)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	logger := zerolog.Nop()
	sa := newSessionAuth(server.URL, "admin@test.com", "password123", logger)

	err := sa.authenticate()
	require.NoError(t, err)
	assert.Equal(t, "test-session-id", sa.getSessionID())
}

func TestSessionAuth_AuthenticateFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"errors":{"password":"did not match"}}`))
	}))
	defer server.Close()

	logger := zerolog.Nop()
	sa := newSessionAuth(server.URL, "admin@test.com", "wrong", logger)

	err := sa.authenticate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "401")
}
