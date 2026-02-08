package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerSettingTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "list_settings", "List all Metabase settings (admin only)",
		inputSchema(map[string]any{}, nil),
		func(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			logger.Debug().Msg("listing settings")
			settings, err := client.ListSettings()
			if err != nil {
				return errResult(err)
			}
			return marshalResult(settings)
		})

	addTool(server, "get_setting", "Get a specific Metabase setting value",
		inputSchema(map[string]any{
			"key": map[string]any{"type": "string", "description": "Setting key (e.g. 'site-name', 'admin-email')"},
		}, []string{"key"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			key, _ := stringArg(args, "key")
			logger.Debug().Str("key", key).Msg("getting setting")
			val, err := client.GetSetting(key)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(val)
		})
}
