package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteQuery(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/dataset", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var body DatasetQueryRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, 1, body.Database)
		assert.Equal(t, "native", body.Type)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(DatasetQueryResponse{
			Status:   "completed",
			RowCount: 2,
			Data: DatasetData{
				Cols: []DatasetCol{{Name: "ID", DisplayName: "ID", BaseType: "type/Integer"}},
				Rows: [][]any{{1}, {2}},
			},
		})
		require.NoError(t, err)
	})

	result, err := client.ExecuteQuery(&DatasetQueryRequest{
		Database: 1,
		Type:     "native",
		Native:   &NativeQuery{Query: "SELECT id FROM users"},
	})
	require.NoError(t, err)
	assert.Equal(t, "completed", result.Status)
	assert.Equal(t, 2, result.RowCount)
}

func TestExportQueryResults(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/dataset/csv", r.URL.Path)
		w.Header().Set("Content-Type", "text/csv")
		_, _ = w.Write([]byte("ID\n1\n2\n"))
	})

	data, err := client.ExportQueryResults(&DatasetQueryRequest{
		Database: 1,
		Type:     "native",
		Native:   &NativeQuery{Query: "SELECT id FROM users"},
	}, "csv")
	require.NoError(t, err)
	assert.Contains(t, string(data), "ID")
}
