package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetField(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/field/1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Field{ID: 1, Name: "EMAIL", BaseType: "type/Text"})
		require.NoError(t, err)
	})

	field, err := client.GetField(1)
	require.NoError(t, err)
	assert.Equal(t, "EMAIL", field.Name)
}

func TestGetFieldValues(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/field/1/values", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(FieldValues{
			FieldID: 1,
			Values:  [][]any{{"alice@test.com"}, {"bob@test.com"}},
		})
		require.NoError(t, err)
	})

	fv, err := client.GetFieldValues(1)
	require.NoError(t, err)
	assert.Equal(t, 1, fv.FieldID)
	assert.Len(t, fv.Values, 2)
}

func TestSearchFieldValues(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/api/field/1/search/1")
		assert.Equal(t, "alice", r.URL.Query().Get("value"))
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(FieldValues{
			FieldID: 1,
			Values:  [][]any{{"alice@test.com"}},
		})
		require.NoError(t, err)
	})

	fv, err := client.SearchFieldValues(1, "alice", 10)
	require.NoError(t, err)
	assert.Len(t, fv.Values, 1)
}
