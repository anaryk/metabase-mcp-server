package config

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the Metabase MCP server.
type Config struct {
	MetabaseURL string
	APIKey      string
	Username    string
	Password    string
	LogLevel    string
	Transport   string
	Port        int
}

// Load parses configuration from command-line flags and environment variables.
// Flags take precedence over environment variables.
func Load(args []string) (*Config, error) {
	fs := flag.NewFlagSet("metabase-mcp-server", flag.ContinueOnError)

	var cfg Config
	fs.StringVar(&cfg.MetabaseURL, "metabase-url", "", "Metabase instance URL")
	fs.StringVar(&cfg.APIKey, "api-key", "", "Metabase API key")
	fs.StringVar(&cfg.Username, "username", "", "Metabase username")
	fs.StringVar(&cfg.Password, "password", "", "Metabase password")
	fs.StringVar(&cfg.LogLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	fs.StringVar(&cfg.Transport, "transport", "stdio", "Transport type: stdio or sse")
	fs.IntVar(&cfg.Port, "port", 8808, "Port for SSE transport")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	// Environment variables as fallback
	if cfg.MetabaseURL == "" {
		cfg.MetabaseURL = os.Getenv("METABASE_URL")
	}
	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv("METABASE_API_KEY")
	}
	if cfg.Username == "" {
		cfg.Username = os.Getenv("METABASE_USERNAME")
	}
	if cfg.Password == "" {
		cfg.Password = os.Getenv("METABASE_PASSWORD")
	}
	if cfg.LogLevel == "info" {
		if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
			cfg.LogLevel = envLevel
		}
	}
	if cfg.Transport == "stdio" {
		if envTransport := os.Getenv("TRANSPORT"); envTransport != "" {
			cfg.Transport = envTransport
		}
	}
	if cfg.Port == 8808 {
		if envPort := os.Getenv("PORT"); envPort != "" {
			if p, err := strconv.Atoi(envPort); err == nil {
				cfg.Port = p
			}
		}
	}

	cfg.MetabaseURL = strings.TrimRight(cfg.MetabaseURL, "/")

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.MetabaseURL == "" {
		return errors.New("metabase URL is required (--metabase-url or METABASE_URL)")
	}
	if c.APIKey == "" && (c.Username == "" || c.Password == "") {
		return errors.New("either API key (--api-key or METABASE_API_KEY) or username/password (--username/--password or METABASE_USERNAME/METABASE_PASSWORD) is required")
	}
	if c.Transport != "stdio" && c.Transport != "sse" {
		return errors.New("transport must be 'stdio' or 'sse'")
	}
	return nil
}
