package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListTables(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/database/1/metadata/tables", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode([]Table{
			{ID: 1, Name: "USERS"},
			{ID: 2, Name: "ORDERS"},
		})
		require.NoError(t, err)
	})

	tables, err := client.ListTables(1)
	require.NoError(t, err)
	assert.Len(t, tables, 2)
}

func TestGetTable(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/table/1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Table{ID: 1, Name: "USERS", DBID: 1})
		require.NoError(t, err)
	})

	tbl, err := client.GetTable(1)
	require.NoError(t, err)
	assert.Equal(t, "USERS", tbl.Name)
}

func TestGetTableMetadata(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/table/1/query_metadata", r.URL.Path)
		assert.Equal(t, "true", r.URL.Query().Get("include_hidden_fields"))
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Table{
			ID: 1, Name: "USERS",
			Fields: []Field{{ID: 1, Name: "ID"}, {ID: 2, Name: "NAME"}},
		})
		require.NoError(t, err)
	})

	tbl, err := client.GetTableMetadata(1)
	require.NoError(t, err)
	assert.Len(t, tbl.Fields, 2)
}

func TestGetTableForeignKeys(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/table/1/fks", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode([]ForeignKey{
			{Relationship: "one-to-many"},
		})
		require.NoError(t, err)
	})

	fks, err := client.GetTableForeignKeys(1)
	require.NoError(t, err)
	assert.Len(t, fks, 1)
}
