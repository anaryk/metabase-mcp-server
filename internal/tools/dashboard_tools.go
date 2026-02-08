package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerDashboardTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "list_dashboards", "List all dashboards",
		inputSchema(map[string]any{}, nil),
		func(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			logger.Debug().Msg("listing dashboards")
			dashboards, err := client.ListDashboards()
			if err != nil {
				return errResult(err)
			}
			return marshalResult(dashboards)
		})

	addTool(server, "get_dashboard", "Get a dashboard by ID including all cards and layout",
		inputSchema(map[string]any{
			"dashboard_id": map[string]any{"type": "number", "description": "The dashboard ID"},
		}, []string{"dashboard_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "dashboard_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("dashboard_id", id).Msg("getting dashboard")
			dash, err := client.GetDashboard(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(dash)
		})

	addTool(server, "create_dashboard", "Create a new dashboard",
		inputSchema(map[string]any{
			"name":          map[string]any{"type": "string", "description": "Dashboard name"},
			"description":   map[string]any{"type": "string", "description": "Dashboard description"},
			"collection_id": map[string]any{"type": "number", "description": "Collection ID"},
			"parameters":    map[string]any{"type": "array", "description": "Dashboard filter parameters"},
		}, []string{"name"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			name, _ := stringArg(args, "name")
			dash := &metabase.Dashboard{
				Name:         name,
				Description:  optionalStringArg(args, "description"),
				CollectionID: optionalIntArg(args, "collection_id"),
			}
			if p, ok := args["parameters"]; ok {
				if params, ok := p.([]any); ok {
					for _, param := range params {
						if m, ok := param.(map[string]any); ok {
							dash.Parameters = append(dash.Parameters, m)
						}
					}
				}
			}
			logger.Debug().Str("name", name).Msg("creating dashboard")
			result, err := client.CreateDashboard(dash)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(result)
		})

	addTool(server, "update_dashboard", "Update dashboard properties",
		inputSchema(map[string]any{
			"dashboard_id":  map[string]any{"type": "number", "description": "The dashboard ID to update"},
			"name":          map[string]any{"type": "string", "description": "New name"},
			"description":   map[string]any{"type": "string", "description": "New description"},
			"archived":      map[string]any{"type": "boolean", "description": "Whether to archive"},
			"collection_id": map[string]any{"type": "number", "description": "New collection ID"},
		}, []string{"dashboard_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "dashboard_id")
			if err != nil {
				return errResult(err)
			}
			dash := &metabase.Dashboard{
				Description:  optionalStringArg(args, "description"),
				Archived:     optionalBoolArg(args, "archived"),
				CollectionID: optionalIntArg(args, "collection_id"),
			}
			if n, ok := args["name"].(string); ok {
				dash.Name = n
			}
			logger.Debug().Int("dashboard_id", id).Msg("updating dashboard")
			result, err := client.UpdateDashboard(id, dash)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(result)
		})

	addTool(server, "delete_dashboard", "Delete a dashboard",
		inputSchema(map[string]any{
			"dashboard_id": map[string]any{"type": "number", "description": "The dashboard ID to delete"},
		}, []string{"dashboard_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "dashboard_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("dashboard_id", id).Msg("deleting dashboard")
			if err := client.DeleteDashboard(id); err != nil {
				return errResult(err)
			}
			return textResult("Dashboard deleted successfully"), nil
		})

	addTool(server, "add_card_to_dashboard", "Add a card to a dashboard with position and size",
		inputSchema(map[string]any{
			"dashboard_id":       map[string]any{"type": "number", "description": "Dashboard ID"},
			"card_id":            map[string]any{"type": "number", "description": "Card ID to add"},
			"row":                map[string]any{"type": "number", "description": "Row position (default: 0)"},
			"col":                map[string]any{"type": "number", "description": "Column position (default: 0)"},
			"size_x":             map[string]any{"type": "number", "description": "Width in grid units (default: 6)"},
			"size_y":             map[string]any{"type": "number", "description": "Height in grid units (default: 4)"},
			"series":             map[string]any{"type": "array", "description": "Series to overlay"},
			"parameter_mappings": map[string]any{"type": "array", "description": "Parameter mappings"},
		}, []string{"dashboard_id", "card_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			dashID, err := intArg(args, "dashboard_id")
			if err != nil {
				return errResult(err)
			}
			cardID, err := intArg(args, "card_id")
			if err != nil {
				return errResult(err)
			}
			dc := &metabase.DashCard{
				CardID: &cardID,
				SizeX:  6,
				SizeY:  4,
			}
			if v := optionalIntArg(args, "row"); v != nil {
				dc.Row = *v
			}
			if v := optionalIntArg(args, "col"); v != nil {
				dc.Col = *v
			}
			if v := optionalIntArg(args, "size_x"); v != nil {
				dc.SizeX = *v
			}
			if v := optionalIntArg(args, "size_y"); v != nil {
				dc.SizeY = *v
			}
			if series, ok := args["series"].([]any); ok {
				for _, s := range series {
					if m, ok := s.(map[string]any); ok {
						dc.Series = append(dc.Series, m)
					}
				}
			}
			if mappings, ok := args["parameter_mappings"].([]any); ok {
				for _, pm := range mappings {
					if m, ok := pm.(map[string]any); ok {
						dc.ParameterMappings = append(dc.ParameterMappings, m)
					}
				}
			}
			logger.Debug().Int("dashboard_id", dashID).Int("card_id", cardID).Msg("adding card to dashboard")
			result, err := client.AddCardToDashboard(dashID, dc)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(result)
		})

	addTool(server, "remove_card_from_dashboard", "Remove a card from a dashboard",
		inputSchema(map[string]any{
			"dashboard_id": map[string]any{"type": "number", "description": "Dashboard ID"},
			"dashcard_id":  map[string]any{"type": "number", "description": "Dashcard ID to remove"},
		}, []string{"dashboard_id", "dashcard_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			dashID, err := intArg(args, "dashboard_id")
			if err != nil {
				return errResult(err)
			}
			dcID, err := intArg(args, "dashcard_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("dashboard_id", dashID).Int("dashcard_id", dcID).Msg("removing card from dashboard")
			if err := client.RemoveCardFromDashboard(dashID, dcID); err != nil {
				return errResult(err)
			}
			return textResult("Card removed from dashboard successfully"), nil
		})

	addTool(server, "update_dashboard_cards", "Update layout/positions of cards on a dashboard",
		inputSchema(map[string]any{
			"dashboard_id": map[string]any{"type": "number", "description": "Dashboard ID"},
			"cards":        map[string]any{"type": "array", "description": "Array of dashcard objects with id, row, col, size_x, size_y"},
		}, []string{"dashboard_id", "cards"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			dashID, err := intArg(args, "dashboard_id")
			if err != nil {
				return errResult(err)
			}
			cardsRaw, ok := args["cards"].([]any)
			if !ok {
				return errResult(fmt.Errorf("cards must be an array"))
			}
			var cards []metabase.DashCard
			rawJSON, _ := json.Marshal(cardsRaw)
			if err := json.Unmarshal(rawJSON, &cards); err != nil {
				return errResult(err)
			}
			logger.Debug().Int("dashboard_id", dashID).Int("card_count", len(cards)).Msg("updating dashboard cards")
			if err := client.UpdateDashboardCards(dashID, cards); err != nil {
				return errResult(err)
			}
			return textResult("Dashboard cards updated successfully"), nil
		})

	addTool(server, "copy_dashboard", "Copy a dashboard to a new collection",
		inputSchema(map[string]any{
			"dashboard_id":  map[string]any{"type": "number", "description": "Dashboard ID to copy"},
			"name":          map[string]any{"type": "string", "description": "Name for the copy"},
			"description":   map[string]any{"type": "string", "description": "Description for the copy"},
			"collection_id": map[string]any{"type": "number", "description": "Target collection ID"},
		}, []string{"dashboard_id", "name"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "dashboard_id")
			if err != nil {
				return errResult(err)
			}
			name, _ := stringArg(args, "name")
			logger.Debug().Int("dashboard_id", id).Str("name", name).Msg("copying dashboard")
			result, err := client.CopyDashboard(id, name, optionalStringArg(args, "description"), optionalIntArg(args, "collection_id"))
			if err != nil {
				return errResult(err)
			}
			return marshalResult(result)
		})
}
