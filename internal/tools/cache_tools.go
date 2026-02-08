package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerCacheTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "invalidate_cache", "Invalidate the Metabase cache",
		inputSchema(map[string]any{}, nil),
		func(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			logger.Debug().Msg("invalidating cache")
			if err := client.InvalidateCache(); err != nil {
				return errResult(err)
			}
			return textResult("Cache invalidated successfully"), nil
		})
}
