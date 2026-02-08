package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListAlerts(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode([]Alert{{ID: 1, CardID: 5}})
		require.NoError(t, err)
	})

	alerts, err := client.ListAlerts()
	require.NoError(t, err)
	assert.Len(t, alerts, 1)
}

func TestGetAlert(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/alert/1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Alert{ID: 1, CardID: 5, AlertCondition: "rows"})
		require.NoError(t, err)
	})

	alert, err := client.GetAlert(1)
	require.NoError(t, err)
	assert.Equal(t, "rows", alert.AlertCondition)
}

func TestCreateAlert(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Alert{ID: 10, CardID: 5})
		require.NoError(t, err)
	})

	alert, err := client.CreateAlert(&Alert{CardID: 5, AlertCondition: "rows"})
	require.NoError(t, err)
	assert.Equal(t, 10, alert.ID)
}
