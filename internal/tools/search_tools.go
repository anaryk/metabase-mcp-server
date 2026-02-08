package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerSearchTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "search", "Search across all Metabase entities (cards, dashboards, collections, tables)",
		inputSchema(map[string]any{
			"query":  map[string]any{"type": "string", "description": "Search query string"},
			"models": map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Filter by model types: card, dashboard, collection, table, database, action"},
		}, []string{"query"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			query, _ := stringArg(args, "query")
			models := stringSliceArg(args, "models")
			logger.Debug().Str("query", query).Msg("searching")
			result, err := client.Search(query, models)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(result)
		})
}
