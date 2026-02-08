package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerDatabaseTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "list_databases", "List all connected databases",
		inputSchema(map[string]any{}, nil),
		func(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			logger.Debug().Msg("listing databases")
			dbs, err := client.ListDatabases()
			if err != nil {
				return errResult(err)
			}
			return marshalResult(dbs)
		})

	addTool(server, "get_database", "Get database details by ID",
		inputSchema(map[string]any{
			"database_id": map[string]any{"type": "number", "description": "The database ID"},
		}, []string{"database_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "database_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("database_id", id).Msg("getting database")
			db, err := client.GetDatabase(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(db)
		})

	addTool(server, "get_database_metadata", "Get full metadata for a database including tables, fields, and types. Essential for understanding the schema before building queries.",
		inputSchema(map[string]any{
			"database_id": map[string]any{"type": "number", "description": "The database ID"},
		}, []string{"database_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "database_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("database_id", id).Msg("getting database metadata")
			db, err := client.GetDatabaseMetadata(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(db)
		})

	addTool(server, "sync_database", "Trigger a schema sync for a database",
		inputSchema(map[string]any{
			"database_id": map[string]any{"type": "number", "description": "The database ID to sync"},
		}, []string{"database_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "database_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("database_id", id).Msg("syncing database")
			if err := client.SyncDatabase(id); err != nil {
				return errResult(err)
			}
			return textResult("Database sync triggered successfully"), nil
		})
}
