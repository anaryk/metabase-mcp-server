package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerActivityTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "get_activity", "Get recent activity log",
		inputSchema(map[string]any{}, nil),
		func(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			logger.Debug().Msg("getting activity")
			activity, err := client.GetActivity()
			if err != nil {
				return errResult(err)
			}
			return marshalResult(activity)
		})

	addTool(server, "get_recent_views", "Get recently viewed items",
		inputSchema(map[string]any{}, nil),
		func(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			logger.Debug().Msg("getting recent views")
			items, err := client.GetRecentViews()
			if err != nil {
				return errResult(err)
			}
			return marshalResult(items)
		})
}
