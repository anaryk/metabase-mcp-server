package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCards(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/card", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode([]Card{
			{ID: 1, Name: "Users Count"},
			{ID: 2, Name: "Revenue"},
		})
		require.NoError(t, err)
	})

	cards, err := client.ListCards()
	require.NoError(t, err)
	assert.Len(t, cards, 2)
}

func TestGetCard(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/card/1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Card{ID: 1, Name: "Users Count", Display: "scalar"})
		require.NoError(t, err)
	})

	card, err := client.GetCard(1)
	require.NoError(t, err)
	assert.Equal(t, "Users Count", card.Name)
}

func TestCreateCard(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/card", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(Card{ID: 10, Name: "New Card"})
		require.NoError(t, err)
	})

	card, err := client.CreateCard(&Card{Name: "New Card"})
	require.NoError(t, err)
	assert.Equal(t, 10, card.ID)
}

func TestUpdateCard(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/api/card/1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Card{ID: 1, Name: "Updated"})
		require.NoError(t, err)
	})

	card, err := client.UpdateCard(1, &Card{Name: "Updated"})
	require.NoError(t, err)
	assert.Equal(t, "Updated", card.Name)
}

func TestDeleteCard(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/card/1", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.DeleteCard(1)
	require.NoError(t, err)
}

func TestExecuteCardQuery(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/card/1/query", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(DatasetQueryResponse{
			Status:   "completed",
			RowCount: 1,
			Data:     DatasetData{Rows: [][]any{{42}}},
		})
		require.NoError(t, err)
	})

	result, err := client.ExecuteCardQuery(1, nil)
	require.NoError(t, err)
	assert.Equal(t, "completed", result.Status)
}
