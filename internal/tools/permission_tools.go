package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerPermissionTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "list_permission_groups", "List all permission groups",
		inputSchema(map[string]any{}, nil),
		func(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			logger.Debug().Msg("listing permission groups")
			groups, err := client.ListPermissionGroups()
			if err != nil {
				return errResult(err)
			}
			return marshalResult(groups)
		})

	addTool(server, "get_permission_group", "Get permission group details with members",
		inputSchema(map[string]any{
			"group_id": map[string]any{"type": "number", "description": "The permission group ID"},
		}, []string{"group_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "group_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("group_id", id).Msg("getting permission group")
			group, err := client.GetPermissionGroup(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(group)
		})

	addTool(server, "get_permissions_graph", "Get the full permissions graph showing all group permissions",
		inputSchema(map[string]any{}, nil),
		func(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			logger.Debug().Msg("getting permissions graph")
			graph, err := client.GetPermissionsGraph()
			if err != nil {
				return errResult(err)
			}
			return marshalResult(graph)
		})
}
