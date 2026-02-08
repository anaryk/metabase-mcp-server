package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerTimelineTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "list_timelines", "List all timelines with optional collection filter",
		inputSchema(map[string]any{
			"collection_id": map[string]any{"type": "number", "description": "Optional collection ID filter"},
		}, nil),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			_ = parseArgs(req, &args)
			colID := optionalIntArg(args, "collection_id")
			logger.Debug().Msg("listing timelines")
			timelines, err := client.ListTimelines(colID)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(timelines)
		})

	addTool(server, "get_timeline", "Get timeline by ID with events",
		inputSchema(map[string]any{
			"timeline_id": map[string]any{"type": "number", "description": "The timeline ID"},
		}, []string{"timeline_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "timeline_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("timeline_id", id).Msg("getting timeline")
			tl, err := client.GetTimeline(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(tl)
		})
}
