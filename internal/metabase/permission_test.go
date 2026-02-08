package metabase

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListPermissionGroups(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode([]PermissionGroup{
			{ID: 1, Name: "All Users"},
			{ID: 2, Name: "Administrators"},
		})
		require.NoError(t, err)
	})

	groups, err := client.ListPermissionGroups()
	require.NoError(t, err)
	assert.Len(t, groups, 2)
}

func TestGetPermissionGroup(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/permissions/group/1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(PermissionGroup{
			ID: 1, Name: "All Users",
			Members: []User{{ID: 1, Email: "admin@test.com"}},
		})
		require.NoError(t, err)
	})

	group, err := client.GetPermissionGroup(1)
	require.NoError(t, err)
	assert.Equal(t, "All Users", group.Name)
	assert.Len(t, group.Members, 1)
}

func TestGetPermissionsGraph(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"revision": 1,
			"groups":   map[string]any{},
		})
		require.NoError(t, err)
	})

	graph, err := client.GetPermissionsGraph()
	require.NoError(t, err)
	assert.Contains(t, graph, "revision")
}
