package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerActionTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "list_actions", "List actions for a model",
		inputSchema(map[string]any{
			"model_id": map[string]any{"type": "number", "description": "The model ID"},
		}, []string{"model_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "model_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("model_id", id).Msg("listing actions")
			actions, err := client.ListActions(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(actions)
		})

	addTool(server, "get_action", "Get action details by ID",
		inputSchema(map[string]any{
			"action_id": map[string]any{"type": "number", "description": "The action ID"},
		}, []string{"action_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "action_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("action_id", id).Msg("getting action")
			action, err := client.GetAction(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(action)
		})
}
