package main

import (
	"context"
	"fmt"
	"net/http"
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
		Str("transport", cfg.Transport).
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

	// Set up context with signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	switch cfg.Transport {
	case "sse":
		return runSSE(ctx, server, cfg.Port, logger)
	default:
		return runStdio(ctx, server, logger)
	}
}

func runStdio(ctx context.Context, server *mcp.Server, logger zerolog.Logger) error {
	logger.Info().Msg("MCP server ready, listening on stdio")
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		logger.Error().Err(err).Msg("server error")
		return err
	}
	return nil
}

func runSSE(ctx context.Context, server *mcp.Server, port int, logger zerolog.Logger) error {
	handler := mcp.NewSSEHandler(func(_ *http.Request) *mcp.Server {
		return server
	}, nil)

	addr := fmt.Sprintf(":%d", port)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info().Str("addr", addr).Msg("MCP server ready, listening on SSE")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("SSE server error: %w", err)
	case <-ctx.Done():
		logger.Info().Msg("shutting down SSE server")
		return httpServer.Close()
	}
}
