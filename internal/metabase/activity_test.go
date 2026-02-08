package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetActivity(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/activity", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode([]ActivityItem{
			{ID: 1, Topic: "card-create"},
		})
		require.NoError(t, err)
	})

	items, err := client.GetActivity()
	require.NoError(t, err)
	assert.Len(t, items, 1)
}

func TestGetRecentViews(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/activity/recent_views", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode([]RecentItem{
			{ID: 1, Name: "Sales Dashboard", Model: "dashboard"},
		})
		require.NoError(t, err)
	})

	items, err := client.GetRecentViews()
	require.NoError(t, err)
	assert.Len(t, items, 1)
}
