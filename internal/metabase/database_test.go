package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListDatabases(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/database", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"data": []Database{
				{ID: 1, Name: "H2", Engine: "h2"},
				{ID: 2, Name: "Postgres", Engine: "postgres"},
			},
		})
		require.NoError(t, err)
	})

	dbs, err := client.ListDatabases()
	require.NoError(t, err)
	assert.Len(t, dbs, 2)
	assert.Equal(t, "H2", dbs[0].Name)
}

func TestGetDatabase(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/database/1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Database{ID: 1, Name: "H2", Engine: "h2"})
		require.NoError(t, err)
	})

	db, err := client.GetDatabase(1)
	require.NoError(t, err)
	assert.Equal(t, 1, db.ID)
	assert.Equal(t, "H2", db.Name)
}

func TestGetDatabaseMetadata(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/database/1/metadata", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Database{
			ID: 1, Name: "H2",
			Tables: []Table{{ID: 1, Name: "USERS"}},
		})
		require.NoError(t, err)
	})

	db, err := client.GetDatabaseMetadata(1)
	require.NoError(t, err)
	assert.Len(t, db.Tables, 1)
	assert.Equal(t, "USERS", db.Tables[0].Name)
}

func TestSyncDatabase(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/database/1/sync_schema", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
	})

	err := client.SyncDatabase(1)
	require.NoError(t, err)
}
