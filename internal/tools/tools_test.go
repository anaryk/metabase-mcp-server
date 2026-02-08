package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func setupTestServer(t *testing.T, handler http.HandlerFunc) (*mcp.Server, *mcp.ClientSession) {
	t.Helper()

	mbServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/user/current" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":1,"email":"test@test.com"}`))
			return
		}
		handler(w, r)
	}))
	t.Cleanup(mbServer.Close)

	logger := zerolog.Nop()
	client, err := metabase.NewClient(mbServer.URL, "test-api-key", "", "", logger)
	require.NoError(t, err)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "metabase-mcp-server",
		Version: "test",
	}, nil)
	RegisterAll(server, client, logger)

	sTransport, cTransport := mcp.NewInMemoryTransports()
	mcpClient := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "test",
	}, nil)

	ctx := context.Background()
	go func() {
		_ = server.Run(ctx, sTransport)
	}()

	session, err := mcpClient.Connect(ctx, cTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })

	return server, session
}

func TestListTools(t *testing.T) {
	_, session := setupTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()
	result, err := session.ListTools(ctx, nil)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(result.Tools), 45, "expected at least 45 tools")

	// Verify some key tools exist
	toolNames := make(map[string]bool)
	for _, tool := range result.Tools {
		toolNames[tool.Name] = true
	}
	expectedTools := []string{
		"list_cards", "get_card", "create_card", "update_card", "delete_card", "execute_card_query",
		"list_dashboards", "get_dashboard", "create_dashboard", "update_dashboard", "delete_dashboard",
		"add_card_to_dashboard", "remove_card_from_dashboard", "update_dashboard_cards", "copy_dashboard",
		"list_collections", "get_collection", "create_collection", "update_collection", "list_collection_items",
		"list_databases", "get_database", "get_database_metadata", "sync_database",
		"list_tables", "get_table", "get_table_metadata", "get_table_fks",
		"get_field", "get_field_values", "search_field_values",
		"execute_query", "export_query_results",
		"list_users", "get_user", "get_current_user",
		"list_permission_groups", "get_permission_group", "get_permissions_graph",
		"search",
		"list_alerts", "get_alert", "create_alert",
		"list_settings", "get_setting",
		"get_activity", "get_recent_views",
		"list_actions", "get_action",
		"list_timelines", "get_timeline",
		"invalidate_cache",
	}
	for _, name := range expectedTools {
		assert.True(t, toolNames[name], "missing tool: %s", name)
	}
}

func TestExecuteQuery_ReadOnly(t *testing.T) {
	_, session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/dataset" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(metabase.DatasetQueryResponse{
				Status:   "completed",
				RowCount: 1,
				Data: metabase.DatasetData{
					Cols: []metabase.DatasetCol{{Name: "count", DisplayName: "Count", BaseType: "type/Integer"}},
					Rows: [][]any{{42}},
				},
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()

	// Valid SELECT query should work
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "execute_query",
		Arguments: map[string]any{
			"database_id":  1,
			"query_type":   "native",
			"native_query": "SELECT COUNT(*) FROM users",
		},
	})
	require.NoError(t, err)
	assert.False(t, result.IsError)

	// INSERT should be blocked
	result, err = session.CallTool(ctx, &mcp.CallToolParams{
		Name: "execute_query",
		Arguments: map[string]any{
			"database_id":  1,
			"query_type":   "native",
			"native_query": "INSERT INTO users (name) VALUES ('test')",
		},
	})
	require.NoError(t, err)
	assert.True(t, result.IsError)
	text := result.Content[0].(*mcp.TextContent).Text
	assert.Contains(t, text, "blocked operation")

	// DROP should be blocked
	result, err = session.CallTool(ctx, &mcp.CallToolParams{
		Name: "execute_query",
		Arguments: map[string]any{
			"database_id":  1,
			"query_type":   "native",
			"native_query": "DROP TABLE users",
		},
	})
	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestListDatabases(t *testing.T) {
	_, session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/database" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []metabase.Database{
					{ID: 1, Name: "Sample Database", Engine: "h2"},
				},
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "list_databases",
	})
	require.NoError(t, err)
	assert.False(t, result.IsError)

	text := result.Content[0].(*mcp.TextContent).Text
	assert.Contains(t, text, "Sample Database")
}

func TestSearch(t *testing.T) {
	_, session := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/search" {
			assert.Equal(t, "revenue", r.URL.Query().Get("q"))
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(metabase.SearchResponse{
				Data:  []metabase.SearchResult{{ID: 1, Name: "Revenue Card", Model: "card"}},
				Total: 1,
			})
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "search",
		Arguments: map[string]any{
			"query": "revenue",
		},
	})
	require.NoError(t, err)
	assert.False(t, result.IsError)

	text := result.Content[0].(*mcp.TextContent).Text
	assert.Contains(t, text, "Revenue Card")
}

func TestStreamableHTTPTransport(t *testing.T) {
	// Set up mock Metabase backend
	mbServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/user/current" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":1,"email":"test@test.com"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(mbServer.Close)

	logger := zerolog.Nop()
	client, err := metabase.NewClient(mbServer.URL, "test-api-key", "", "", logger)
	require.NoError(t, err)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "metabase-mcp-server",
		Version: "test",
	}, nil)
	RegisterAll(server, client, logger)

	// Start an HTTP server with StreamableHTTPHandler
	handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		return server
	}, nil)
	httpServer := httptest.NewServer(handler)
	t.Cleanup(httpServer.Close)

	// Connect via StreamableClientTransport
	mcpClient := mcp.NewClient(&mcp.Implementation{
		Name:    "test-streamable-client",
		Version: "test",
	}, nil)

	ctx := context.Background()
	transport := &mcp.StreamableClientTransport{
		Endpoint:   httpServer.URL,
		HTTPClient: httpServer.Client(),
	}
	session, err := mcpClient.Connect(ctx, transport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })

	// Verify tools are registered via the Streamable HTTP transport
	result, err := session.ListTools(ctx, nil)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(result.Tools), 45, "expected at least 45 tools via streamable HTTP")
}
