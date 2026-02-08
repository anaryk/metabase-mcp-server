package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListSettings(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode([]Setting{
			{Key: "site-name", Value: "My Metabase"},
		})
		require.NoError(t, err)
	})

	settings, err := client.ListSettings()
	require.NoError(t, err)
	assert.Len(t, settings, 1)
}

func TestGetSetting(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/setting/site-name", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`"My Metabase"`))
	})

	val, err := client.GetSetting("site-name")
	require.NoError(t, err)
	assert.Equal(t, "My Metabase", val)
}
