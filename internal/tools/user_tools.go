package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerUserTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "list_users", "List all Metabase users",
		inputSchema(map[string]any{}, nil),
		func(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			logger.Debug().Msg("listing users")
			users, err := client.ListUsers()
			if err != nil {
				return errResult(err)
			}
			return marshalResult(users)
		})

	addTool(server, "get_user", "Get a user by ID",
		inputSchema(map[string]any{
			"user_id": map[string]any{"type": "number", "description": "The user ID"},
		}, []string{"user_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "user_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("user_id", id).Msg("getting user")
			user, err := client.GetUser(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(user)
		})

	addTool(server, "get_current_user", "Get the currently authenticated user",
		inputSchema(map[string]any{}, nil),
		func(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			logger.Debug().Msg("getting current user")
			user, err := client.GetCurrentUser()
			if err != nil {
				return errResult(err)
			}
			return marshalResult(user)
		})
}
