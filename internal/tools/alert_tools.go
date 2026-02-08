package tools

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerAlertTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "list_alerts", "List all alerts",
		inputSchema(map[string]any{}, nil),
		func(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			logger.Debug().Msg("listing alerts")
			alerts, err := client.ListAlerts()
			if err != nil {
				return errResult(err)
			}
			return marshalResult(alerts)
		})

	addTool(server, "get_alert", "Get alert details by ID",
		inputSchema(map[string]any{
			"alert_id": map[string]any{"type": "number", "description": "The alert ID"},
		}, []string{"alert_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "alert_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("alert_id", id).Msg("getting alert")
			alert, err := client.GetAlert(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(alert)
		})

	addTool(server, "create_alert", "Create a new alert on a card",
		inputSchema(map[string]any{
			"card_id":          map[string]any{"type": "number", "description": "Card ID to alert on"},
			"alert_condition":  map[string]any{"type": "string", "description": "Alert condition: 'rows' or 'goal'"},
			"alert_above_goal": map[string]any{"type": "boolean", "description": "Alert when above goal (for goal condition)"},
			"alert_first_only": map[string]any{"type": "boolean", "description": "Only alert on first match"},
			"channels":         map[string]any{"type": "array", "description": "Notification channels"},
		}, []string{"card_id", "alert_condition"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			cardID, err := intArg(args, "card_id")
			if err != nil {
				return errResult(err)
			}
			condition, _ := stringArg(args, "alert_condition")
			alert := &metabase.Alert{
				CardID:         cardID,
				AlertCondition: condition,
				AlertAboveGoal: optionalBoolArg(args, "alert_above_goal"),
			}
			if v := optionalBoolArg(args, "alert_first_only"); v != nil {
				alert.AlertFirstOnly = *v
			}
			if channels, ok := args["channels"].([]any); ok {
				raw, _ := json.Marshal(channels)
				var ch []map[string]any
				_ = json.Unmarshal(raw, &ch)
				alert.Channels = ch
			}
			logger.Debug().Int("card_id", cardID).Str("condition", condition).Msg("creating alert")
			result, err := client.CreateAlert(alert)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(result)
		})
}
