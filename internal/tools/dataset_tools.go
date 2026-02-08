package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerDatasetTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "execute_query", "Execute a native SQL or MBQL query against a database. IMPORTANT: Only read-only (SELECT) queries are allowed - write operations are blocked.",
		inputSchema(map[string]any{
			"database_id":   map[string]any{"type": "number", "description": "The database ID to query"},
			"query_type":    map[string]any{"type": "string", "description": "Query type: 'native' for SQL or 'query' for MBQL", "enum": []string{"native", "query"}},
			"native_query":  map[string]any{"type": "string", "description": "SQL query string (for native type)"},
			"mbql_query":    map[string]any{"type": "object", "description": "MBQL query object (for query type)"},
			"template_tags": map[string]any{"type": "object", "description": "Template tags for parameterized native queries"},
		}, []string{"database_id", "query_type"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			dbID, err := intArg(args, "database_id")
			if err != nil {
				return errResult(err)
			}
			queryType, _ := stringArg(args, "query_type")

			dsReq := &metabase.DatasetQueryRequest{
				Database: dbID,
				Type:     queryType,
			}

			if queryType == "native" {
				sql, _ := stringArg(args, "native_query")
				if sql == "" {
					return errResult(fmt.Errorf("native_query is required for native query type"))
				}
				// Enforce read-only SQL
				if err := metabase.ValidateReadOnlySQL(sql); err != nil {
					logger.Warn().Str("query", sql).Msg("blocked write query attempt")
					return errResult(err)
				}
				dsReq.Native = &metabase.NativeQuery{
					Query:        sql,
					TemplateTags: mapArg(args, "template_tags"),
				}
			} else {
				mbql := mapArg(args, "mbql_query")
				if mbql == nil {
					return errResult(fmt.Errorf("mbql_query is required for query type"))
				}
				dsReq.Query = mbql
			}

			logger.Debug().Int("database_id", dbID).Str("type", queryType).Msg("executing query")
			result, err := client.ExecuteQuery(dsReq)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(result)
		})

	addTool(server, "export_query_results", "Export query results as CSV, JSON, or XLSX. Only read-only queries are allowed.",
		inputSchema(map[string]any{
			"database_id":   map[string]any{"type": "number", "description": "The database ID"},
			"query_type":    map[string]any{"type": "string", "description": "Query type: 'native' or 'query'", "enum": []string{"native", "query"}},
			"native_query":  map[string]any{"type": "string", "description": "SQL query (for native type)"},
			"mbql_query":    map[string]any{"type": "object", "description": "MBQL query (for query type)"},
			"export_format": map[string]any{"type": "string", "description": "Export format", "enum": []string{"csv", "json", "xlsx"}},
		}, []string{"database_id", "query_type", "export_format"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			dbID, err := intArg(args, "database_id")
			if err != nil {
				return errResult(err)
			}
			queryType, _ := stringArg(args, "query_type")
			format, _ := stringArg(args, "export_format")

			dsReq := &metabase.DatasetQueryRequest{
				Database: dbID,
				Type:     queryType,
			}

			if queryType == "native" {
				sql, _ := stringArg(args, "native_query")
				if sql == "" {
					return errResult(fmt.Errorf("native_query is required for native query type"))
				}
				if err := metabase.ValidateReadOnlySQL(sql); err != nil {
					logger.Warn().Str("query", sql).Msg("blocked write query attempt in export")
					return errResult(err)
				}
				dsReq.Native = &metabase.NativeQuery{Query: sql}
			} else {
				mbql := mapArg(args, "mbql_query")
				if mbql == nil {
					return errResult(fmt.Errorf("mbql_query is required for query type"))
				}
				dsReq.Query = mbql
			}

			logger.Debug().Int("database_id", dbID).Str("format", format).Msg("exporting query results")
			data, err := client.ExportQueryResults(dsReq, format)
			if err != nil {
				return errResult(err)
			}
			return textResult(string(data)), nil
		})
}
