package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/search", r.URL.Path)
		assert.Equal(t, "revenue", r.URL.Query().Get("q"))
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(SearchResponse{
			Data:  []SearchResult{{ID: 1, Name: "Revenue Card", Model: "card"}},
			Total: 1,
		})
		require.NoError(t, err)
	})

	result, err := client.Search("revenue", nil)
	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Data, 1)
}
