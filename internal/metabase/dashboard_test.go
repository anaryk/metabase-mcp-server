package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListDashboards(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/dashboard", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode([]Dashboard{
			{ID: 1, Name: "Sales"},
		})
		require.NoError(t, err)
	})

	dashboards, err := client.ListDashboards()
	require.NoError(t, err)
	assert.Len(t, dashboards, 1)
}

func TestGetDashboard(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/dashboard/1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Dashboard{ID: 1, Name: "Sales", DashCards: []DashCard{{ID: 1}}})
		require.NoError(t, err)
	})

	dash, err := client.GetDashboard(1)
	require.NoError(t, err)
	assert.Equal(t, "Sales", dash.Name)
	assert.Len(t, dash.DashCards, 1)
}

func TestCreateDashboard(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Dashboard{ID: 10, Name: "New Dash"})
		require.NoError(t, err)
	})

	dash, err := client.CreateDashboard(&Dashboard{Name: "New Dash"})
	require.NoError(t, err)
	assert.Equal(t, 10, dash.ID)
}

func TestDeleteDashboard(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.DeleteDashboard(1)
	require.NoError(t, err)
}

func TestAddCardToDashboard(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/dashboard/1/cards", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		cardID := 5
		err := json.NewEncoder(w).Encode(DashCard{ID: 1, CardID: &cardID, Row: 0, Col: 0})
		require.NoError(t, err)
	})

	cardID := 5
	dc, err := client.AddCardToDashboard(1, &DashCard{CardID: &cardID, Row: 0, Col: 0, SizeX: 6, SizeY: 4})
	require.NoError(t, err)
	assert.Equal(t, 1, dc.ID)
}

func TestRemoveCardFromDashboard(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "10", r.URL.Query().Get("dashcardId"))
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.RemoveCardFromDashboard(1, 10)
	require.NoError(t, err)
}

func TestCopyDashboard(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/dashboard/1/copy", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Dashboard{ID: 20, Name: "Copy of Sales"})
		require.NoError(t, err)
	})

	dash, err := client.CopyDashboard(1, "Copy of Sales", nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 20, dash.ID)
}
