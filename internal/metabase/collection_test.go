package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCollections(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode([]Collection{
			{ID: 1, Name: "Our analytics"},
		})
		require.NoError(t, err)
	})

	cols, err := client.ListCollections("")
	require.NoError(t, err)
	assert.Len(t, cols, 1)
}

func TestGetCollection(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/collection/root", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Collection{ID: "root", Name: "Our analytics"})
		require.NoError(t, err)
	})

	col, err := client.GetCollection("root")
	require.NoError(t, err)
	assert.Equal(t, "Our analytics", col.Name)
}

func TestCreateCollection(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Collection{ID: 5, Name: "New Collection"})
		require.NoError(t, err)
	})

	col, err := client.CreateCollection(&Collection{Name: "New Collection"})
	require.NoError(t, err)
	assert.Equal(t, "New Collection", col.Name)
}

func TestListCollectionItems(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/collection/1/items", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"data": []CollectionItem{
				{ID: 1, Name: "Dashboard 1", Model: "dashboard"},
			},
		})
		require.NoError(t, err)
	})

	items, err := client.ListCollectionItems("1", nil)
	require.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "dashboard", items[0].Model)
}
