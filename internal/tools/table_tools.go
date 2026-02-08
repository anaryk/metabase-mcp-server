package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerTableTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "list_tables", "List all tables for a database",
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
			logger.Debug().Int("database_id", id).Msg("listing tables")
			tables, err := client.ListTables(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(tables)
		})

	addTool(server, "get_table", "Get table details by ID",
		inputSchema(map[string]any{
			"table_id": map[string]any{"type": "number", "description": "The table ID"},
		}, []string{"table_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "table_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("table_id", id).Msg("getting table")
			table, err := client.GetTable(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(table)
		})

	addTool(server, "get_table_metadata", "Get table metadata with all fields and foreign keys. Essential for understanding table structure before building queries.",
		inputSchema(map[string]any{
			"table_id": map[string]any{"type": "number", "description": "The table ID"},
		}, []string{"table_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "table_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("table_id", id).Msg("getting table metadata")
			table, err := client.GetTableMetadata(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(table)
		})

	addTool(server, "get_table_fks", "Get foreign key relationships for a table",
		inputSchema(map[string]any{
			"table_id": map[string]any{"type": "number", "description": "The table ID"},
		}, []string{"table_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "table_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("table_id", id).Msg("getting table foreign keys")
			fks, err := client.GetTableForeignKeys(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(fks)
		})
}
