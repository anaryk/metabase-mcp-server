package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerFieldTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "get_field", "Get field details by ID including type and visibility",
		inputSchema(map[string]any{
			"field_id": map[string]any{"type": "number", "description": "The field ID"},
		}, []string{"field_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "field_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("field_id", id).Msg("getting field")
			field, err := client.GetField(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(field)
		})

	addTool(server, "get_field_values", "Get distinct values for a field (useful for building filters)",
		inputSchema(map[string]any{
			"field_id": map[string]any{"type": "number", "description": "The field ID"},
		}, []string{"field_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "field_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("field_id", id).Msg("getting field values")
			fv, err := client.GetFieldValues(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(fv)
		})

	addTool(server, "search_field_values", "Search field values by prefix",
		inputSchema(map[string]any{
			"field_id": map[string]any{"type": "number", "description": "The field ID"},
			"query":    map[string]any{"type": "string", "description": "Search prefix"},
			"limit":    map[string]any{"type": "number", "description": "Maximum number of results"},
		}, []string{"field_id", "query"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "field_id")
			if err != nil {
				return errResult(err)
			}
			query, _ := stringArg(args, "query")
			limit := 0
			if l := optionalIntArg(args, "limit"); l != nil {
				limit = *l
			}
			logger.Debug().Int("field_id", id).Str("query", query).Msg("searching field values")
			fv, err := client.SearchFieldValues(id, query, limit)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(fv)
		})
}
