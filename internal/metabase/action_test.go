package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListActions(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/action", r.URL.Path)
		assert.Equal(t, "1", r.URL.Query().Get("model-id"))
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode([]Action{
			{ID: 1, Name: "Create User", ModelID: 1},
		})
		require.NoError(t, err)
	})

	actions, err := client.ListActions(1)
	require.NoError(t, err)
	assert.Len(t, actions, 1)
}

func TestGetAction(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/action/1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Action{ID: 1, Name: "Create User", Type: "query"})
		require.NoError(t, err)
	})

	action, err := client.GetAction(1)
	require.NoError(t, err)
	assert.Equal(t, "Create User", action.Name)
}
