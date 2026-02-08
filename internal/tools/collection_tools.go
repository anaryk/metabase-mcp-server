package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerCollectionTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "list_collections", "List all collections",
		inputSchema(map[string]any{
			"namespace": map[string]any{"type": "string", "description": "Optional namespace filter"},
		}, nil),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			_ = parseArgs(req, &args)
			ns := ""
			if s := optionalStringArg(args, "namespace"); s != nil {
				ns = *s
			}
			logger.Debug().Msg("listing collections")
			collections, err := client.ListCollections(ns)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(collections)
		})

	addTool(server, "get_collection", "Get collection details by ID",
		inputSchema(map[string]any{
			"collection_id": map[string]any{"type": "string", "description": "Collection ID (number or 'root')"},
		}, []string{"collection_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id := fmt.Sprintf("%v", args["collection_id"])
			logger.Debug().Str("collection_id", id).Msg("getting collection")
			col, err := client.GetCollection(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(col)
		})

	addTool(server, "create_collection", "Create a new collection",
		inputSchema(map[string]any{
			"name":        map[string]any{"type": "string", "description": "Collection name"},
			"description": map[string]any{"type": "string", "description": "Collection description"},
			"parent_id":   map[string]any{"type": "number", "description": "Parent collection ID"},
			"color":       map[string]any{"type": "string", "description": "Collection color (hex)"},
		}, []string{"name"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			name, _ := stringArg(args, "name")
			col := &metabase.Collection{
				Name:        name,
				Description: optionalStringArg(args, "description"),
				ParentID:    optionalIntArg(args, "parent_id"),
				Color:       optionalStringArg(args, "color"),
			}
			logger.Debug().Str("name", name).Msg("creating collection")
			result, err := client.CreateCollection(col)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(result)
		})

	addTool(server, "update_collection", "Update a collection",
		inputSchema(map[string]any{
			"collection_id": map[string]any{"type": "number", "description": "Collection ID to update"},
			"name":          map[string]any{"type": "string", "description": "New name"},
			"description":   map[string]any{"type": "string", "description": "New description"},
			"color":         map[string]any{"type": "string", "description": "New color"},
			"archived":      map[string]any{"type": "boolean", "description": "Whether to archive"},
		}, []string{"collection_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "collection_id")
			if err != nil {
				return errResult(err)
			}
			col := &metabase.Collection{
				Description: optionalStringArg(args, "description"),
				Color:       optionalStringArg(args, "color"),
				Archived:    optionalBoolArg(args, "archived"),
			}
			if n, ok := args["name"].(string); ok {
				col.Name = n
			}
			logger.Debug().Int("collection_id", id).Msg("updating collection")
			result, err := client.UpdateCollection(id, col)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(result)
		})

	addTool(server, "list_collection_items", "List items in a collection with optional model filter",
		inputSchema(map[string]any{
			"collection_id": map[string]any{"type": "string", "description": "Collection ID (number or 'root')"},
			"models":        map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Filter by model types: card, dashboard, collection, etc."},
		}, []string{"collection_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id := fmt.Sprintf("%v", args["collection_id"])
			models := stringSliceArg(args, "models")
			logger.Debug().Str("collection_id", id).Msg("listing collection items")
			items, err := client.ListCollectionItems(id, models)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(items)
		})
}
