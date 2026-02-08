package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/metabase"
)

// RegisterAll registers all Metabase tools on the given MCP server.
func RegisterAll(server *mcp.Server, client *metabase.Client, logger zerolog.Logger) {
	registerCardTools(server, client, logger)
	registerDashboardTools(server, client, logger)
	registerCollectionTools(server, client, logger)
	registerDatabaseTools(server, client, logger)
	registerTableTools(server, client, logger)
	registerFieldTools(server, client, logger)
	registerDatasetTools(server, client, logger)
	registerUserTools(server, client, logger)
	registerPermissionTools(server, client, logger)
	registerSearchTools(server, client, logger)
	registerAlertTools(server, client, logger)
	registerSettingTools(server, client, logger)
	registerActivityTools(server, client, logger)
	registerActionTools(server, client, logger)
	registerTimelineTools(server, client, logger)
	registerCacheTools(server, client, logger)
}

// marshalResult marshals a value to JSON and returns it as a CallToolResult.
func marshalResult(v any) (*mcp.CallToolResult, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling result: %w", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil
}

// textResult creates a CallToolResult with a text message.
func textResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: msg},
		},
	}
}

// errResult creates a CallToolResult representing an error.
func errResult(err error) (*mcp.CallToolResult, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: err.Error()},
		},
		IsError: true,
	}, nil
}

// inputSchema creates a JSON schema object for tool input.
func inputSchema(properties map[string]any, required []string) json.RawMessage {
	schema := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	data, _ := json.Marshal(schema)
	return data
}

// parseArgs unmarshals the tool arguments from the request.
func parseArgs(req *mcp.CallToolRequest, v any) error {
	if req.Params.Arguments == nil {
		return nil
	}
	return json.Unmarshal(req.Params.Arguments, v)
}

// intArg extracts an integer argument from the parsed args map.
func intArg(args map[string]any, key string) (int, error) {
	v, ok := args[key]
	if !ok {
		return 0, fmt.Errorf("missing required argument: %s", key)
	}
	// JSON numbers are float64
	f, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("argument %s must be a number", key)
	}
	return int(f), nil
}

// optionalIntArg extracts an optional integer argument.
func optionalIntArg(args map[string]any, key string) *int {
	v, ok := args[key]
	if !ok || v == nil {
		return nil
	}
	f, ok := v.(float64)
	if !ok {
		return nil
	}
	i := int(f)
	return &i
}

// stringArg extracts a string argument.
func stringArg(args map[string]any, key string) (string, error) {
	v, ok := args[key]
	if !ok {
		return "", fmt.Errorf("missing required argument: %s", key)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("argument %s must be a string", key)
	}
	return s, nil
}

// optionalStringArg extracts an optional string argument.
func optionalStringArg(args map[string]any, key string) *string {
	v, ok := args[key]
	if !ok || v == nil {
		return nil
	}
	s, ok := v.(string)
	if !ok {
		return nil
	}
	return &s
}

// optionalBoolArg extracts an optional boolean argument.
func optionalBoolArg(args map[string]any, key string) *bool {
	v, ok := args[key]
	if !ok || v == nil {
		return nil
	}
	b, ok := v.(bool)
	if !ok {
		return nil
	}
	return &b
}

// stringSliceArg extracts an optional string slice from the args.
func stringSliceArg(args map[string]any, key string) []string {
	v, ok := args[key]
	if !ok || v == nil {
		return nil
	}
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(arr))
	for _, item := range arr {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

// mapArg extracts an optional map argument.
func mapArg(args map[string]any, key string) map[string]any {
	v, ok := args[key]
	if !ok || v == nil {
		return nil
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	return m
}

// addTool is a convenience wrapper to add a tool with a raw JSON input schema.
func addTool(server *mcp.Server, name, description string, schema json.RawMessage, handler func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	server.AddTool(
		&mcp.Tool{
			Name:        name,
			Description: description,
			InputSchema: schema,
		},
		handler,
	)
}
