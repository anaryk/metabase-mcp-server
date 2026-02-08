package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListUsers(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"data": []User{{ID: 1, Email: "admin@test.com"}},
		})
		require.NoError(t, err)
	})

	users, err := client.ListUsers()
	require.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestGetUser(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/user/1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(User{ID: 1, Email: "admin@test.com"})
		require.NoError(t, err)
	})

	user, err := client.GetUser(1)
	require.NoError(t, err)
	assert.Equal(t, "admin@test.com", user.Email)
}

func TestGetCurrentUser(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// newTestServer returns a valid User for /api/user/current
	user, err := client.GetCurrentUser()
	require.NoError(t, err)
	assert.Equal(t, "test@test.com", user.Email)
}
