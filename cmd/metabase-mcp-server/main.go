package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/zerolog"

	"github.com/anaryk/metabase-mcp-server/internal/config"
	"github.com/anaryk/metabase-mcp-server/internal/metabase"
	"github.com/anaryk/metabase-mcp-server/internal/tools"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load(os.Args[1:])
	if err != nil {
		return err
	}

	// Configure zerolog to write to stderr (stdout is reserved for MCP stdio transport)
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	logger := zerolog.New(os.Stderr).
		With().
		Timestamp().
		Str("component", "metabase-mcp-server").
		Logger().
		Level(level)

	logger.Info().
		Str("version", version).
		Str("metabase_url", cfg.MetabaseURL).
		Msg("starting metabase MCP server")

	// Create Metabase API client
	client, err := metabase.NewClient(cfg.MetabaseURL, cfg.APIKey, cfg.Username, cfg.Password, logger)
	if err != nil {
		return fmt.Errorf("creating metabase client: %w", err)
	}

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "metabase-mcp-server",
		Version: version,
	}, &mcp.ServerOptions{
		Instructions: "Metabase MCP Server provides tools to interact with a Metabase instance. " +
			"You can manage dashboards, cards (saved questions), collections, run queries, and more. " +
			"SQL queries are restricted to read-only (SELECT) operations for safety.",
	})

	// Register all tools
	tools.RegisterAll(server, client, logger)

	logger.Info().Msg("MCP server ready, listening on stdio")

	// Set up context with signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Run MCP server over stdio
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		logger.Error().Err(err).Msg("server error")
		return err
	}

	return nil
}
