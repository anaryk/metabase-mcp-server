package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListTimelines(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode([]Timeline{
			{ID: 1, Name: "Releases"},
		})
		require.NoError(t, err)
	})

	timelines, err := client.ListTimelines(nil)
	require.NoError(t, err)
	assert.Len(t, timelines, 1)
}

func TestGetTimeline(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/timeline/1", r.URL.Path)
		assert.Equal(t, "events", r.URL.Query().Get("include"))
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Timeline{
			ID: 1, Name: "Releases",
			Events: []TimelineEvent{{ID: 1, Name: "v1.0"}},
		})
		require.NoError(t, err)
	})

	tl, err := client.GetTimeline(1)
	require.NoError(t, err)
	assert.Equal(t, "Releases", tl.Name)
	assert.Len(t, tl.Events, 1)
}
