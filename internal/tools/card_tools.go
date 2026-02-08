package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

func registerCardTools(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	addTool(server, "list_cards", "List all saved questions/cards in Metabase",
		inputSchema(map[string]any{}, nil),
		func(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			logger.Debug().Msg("listing cards")
			cards, err := client.ListCards()
			if err != nil {
				return errResult(err)
			}
			return marshalResult(cards)
		})

	addTool(server, "get_card", "Get a saved question/card by ID",
		inputSchema(map[string]any{
			"card_id": map[string]any{"type": "number", "description": "The card ID"},
		}, []string{"card_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "card_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("card_id", id).Msg("getting card")
			card, err := client.GetCard(id)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(card)
		})

	addTool(server, "create_card", "Create a new saved question/card",
		inputSchema(map[string]any{
			"name":                   map[string]any{"type": "string", "description": "Card name"},
			"dataset_query":          map[string]any{"type": "object", "description": "The query definition (native or MBQL)"},
			"display":                map[string]any{"type": "string", "description": "Display type (table, bar, line, pie, scalar, etc.)"},
			"collection_id":          map[string]any{"type": "number", "description": "Collection ID to put the card in"},
			"description":            map[string]any{"type": "string", "description": "Card description"},
			"visualization_settings": map[string]any{"type": "object", "description": "Visualization settings"},
		}, []string{"name", "dataset_query", "display"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			name, _ := stringArg(args, "name")
			card := &metabase.Card{
				Name:                  name,
				DatasetQuery:          mapArg(args, "dataset_query"),
				Display:               args["display"].(string),
				CollectionID:          optionalIntArg(args, "collection_id"),
				Description:           optionalStringArg(args, "description"),
				VisualizationSettings: mapArg(args, "visualization_settings"),
			}
			logger.Debug().Str("name", name).Msg("creating card")
			result, err := client.CreateCard(card)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(result)
		})

	addTool(server, "update_card", "Update an existing saved question/card",
		inputSchema(map[string]any{
			"card_id":                map[string]any{"type": "number", "description": "The card ID to update"},
			"name":                   map[string]any{"type": "string", "description": "New name"},
			"description":            map[string]any{"type": "string", "description": "New description"},
			"dataset_query":          map[string]any{"type": "object", "description": "New query definition"},
			"display":                map[string]any{"type": "string", "description": "New display type"},
			"archived":               map[string]any{"type": "boolean", "description": "Whether to archive the card"},
			"collection_id":          map[string]any{"type": "number", "description": "New collection ID"},
			"visualization_settings": map[string]any{"type": "object", "description": "New visualization settings"},
			"enable_embedding":       map[string]any{"type": "boolean", "description": "Enable embedding"},
			"embedding_params":       map[string]any{"type": "object", "description": "Embedding parameters"},
		}, []string{"card_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "card_id")
			if err != nil {
				return errResult(err)
			}
			card := &metabase.Card{
				Description:           optionalStringArg(args, "description"),
				DatasetQuery:          mapArg(args, "dataset_query"),
				Archived:              optionalBoolArg(args, "archived"),
				CollectionID:          optionalIntArg(args, "collection_id"),
				VisualizationSettings: mapArg(args, "visualization_settings"),
				EnableEmbedding:       optionalBoolArg(args, "enable_embedding"),
				EmbeddingParams:       mapArg(args, "embedding_params"),
			}
			if n, ok := args["name"].(string); ok {
				card.Name = n
			}
			if d, ok := args["display"].(string); ok {
				card.Display = d
			}
			logger.Debug().Int("card_id", id).Msg("updating card")
			result, err := client.UpdateCard(id, card)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(result)
		})

	addTool(server, "delete_card", "Delete (archive) a saved question/card",
		inputSchema(map[string]any{
			"card_id": map[string]any{"type": "number", "description": "The card ID to delete"},
		}, []string{"card_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "card_id")
			if err != nil {
				return errResult(err)
			}
			logger.Debug().Int("card_id", id).Msg("deleting card")
			if err := client.DeleteCard(id); err != nil {
				return errResult(err)
			}
			return textResult("Card deleted successfully"), nil
		})

	addTool(server, "execute_card_query", "Run a saved question's query and return results",
		inputSchema(map[string]any{
			"card_id":    map[string]any{"type": "number", "description": "The card ID"},
			"parameters": map[string]any{"type": "object", "description": "Optional query parameters"},
		}, []string{"card_id"}),
		func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args map[string]any
			if err := parseArgs(req, &args); err != nil {
				return errResult(err)
			}
			id, err := intArg(args, "card_id")
			if err != nil {
				return errResult(err)
			}
			params := mapArg(args, "parameters")
			logger.Debug().Int("card_id", id).Msg("executing card query")
			result, err := client.ExecuteCardQuery(id, params)
			if err != nil {
				return errResult(err)
			}
			return marshalResult(result)
		})
}
