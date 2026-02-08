# Metabase MCP Server

A Model Context Protocol (MCP) server that provides AI agents with full access to a self-hosted Metabase instance. Built in Go, it enables creating dashboards, managing collections, running queries, and more -- all through a standardized MCP interface.

## Features

- **45 MCP tools** covering the complete Metabase API surface
- **Read-only SQL enforcement** -- write operations (INSERT, UPDATE, DELETE, DROP, etc.) are blocked at the server level
- **Two authentication modes** -- API key or session-based (username/password)
- **Structured JSON logging** via zerolog
- **stdio transport** -- standard MCP transport compatible with all MCP clients
- **Production-ready** -- CI/CD, Docker support, cross-platform binaries

## Available Tools

| Category | Tools | Description |
|---|---|---|
| Cards | 6 | List, get, create, update, delete, execute saved questions |
| Dashboards | 9 | Full dashboard management including card placement and copying |
| Collections | 5 | Manage collections and browse collection items |
| Databases | 4 | List databases, inspect metadata, trigger sync |
| Tables | 4 | List tables, get metadata, inspect foreign keys |
| Fields | 3 | Get field details, distinct values, search values |
| Dataset/Query | 2 | Execute SQL/MBQL queries, export results (read-only enforced) |
| Users | 3 | List users, get user details |
| Permissions | 3 | Inspect permission groups and graphs |
| Search | 1 | Search across all entity types |
| Alerts | 3 | List, get, create alerts |
| Settings | 2 | List and get Metabase settings |
| Activity | 2 | View activity log and recent views |
| Actions | 2 | List and get model actions |
| Timelines | 2 | List and get timelines with events |
| Cache | 1 | Invalidate Metabase cache |

## Installation

### From Source

```bash
go install github.com/anaryk/metabase-mcp-server/cmd/metabase-mcp-server@latest
```

### From Binary Releases

Download the appropriate binary from the [Releases](https://github.com/anaryk/metabase-mcp-server/releases) page.

### Docker

```bash
docker build -t metabase-mcp-server .
```

## Configuration

The server accepts configuration via command-line flags or environment variables. Flags take precedence over environment variables.

| Flag | Environment Variable | Required | Description |
|---|---|---|---|
| `--metabase-url` | `METABASE_URL` | Yes | Metabase instance URL |
| `--api-key` | `METABASE_API_KEY` | One of auth | Metabase API key |
| `--username` | `METABASE_USERNAME` | One of auth | Username for session auth |
| `--password` | `METABASE_PASSWORD` | One of auth | Password for session auth |
| `--log-level` | `LOG_LEVEL` | No | Log level: debug, info, warn, error (default: info) |

Either an API key or a username/password pair is required.

### Generating a Metabase API Key

1. Log in to Metabase as an admin.
2. Go to **Admin** > **Settings** > **Authentication** > **API Keys**.
3. Click **Create API Key**.
4. Copy the generated key.

## Usage with Claude Desktop

Add the following to your Claude Desktop configuration file:

- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`

### Using a binary (recommended)

```json
{
  "mcpServers": {
    "metabase": {
      "command": "metabase-mcp-server",
      "args": [
        "--metabase-url", "http://your-metabase-instance:3000",
        "--api-key", "mb_your_api_key_here"
      ]
    }
  }
}
```

### Using environment variables

```json
{
  "mcpServers": {
    "metabase": {
      "command": "metabase-mcp-server",
      "env": {
        "METABASE_URL": "http://your-metabase-instance:3000",
        "METABASE_API_KEY": "mb_your_api_key_here"
      }
    }
  }
}
```

### Using session authentication

```json
{
  "mcpServers": {
    "metabase": {
      "command": "metabase-mcp-server",
      "args": [
        "--metabase-url", "http://your-metabase-instance:3000",
        "--username", "admin@example.com",
        "--password", "your_password"
      ]
    }
  }
}
```

### Using Docker

```json
{
  "mcpServers": {
    "metabase": {
      "command": "docker",
      "args": [
        "run", "--rm", "-i",
        "-e", "METABASE_URL=http://your-metabase-instance:3000",
        "-e", "METABASE_API_KEY=mb_your_api_key_here",
        "metabase-mcp-server"
      ]
    }
  }
}
```

## Usage with VS Code

Add the following to your VS Code settings (`.vscode/settings.json` in your workspace or in your global settings):

```json
{
  "mcp": {
    "servers": {
      "metabase": {
        "command": "metabase-mcp-server",
        "args": [
          "--metabase-url", "http://your-metabase-instance:3000",
          "--api-key", "mb_your_api_key_here"
        ]
      }
    }
  }
}
```

Alternatively, using environment variables:

```json
{
  "mcp": {
    "servers": {
      "metabase": {
        "command": "metabase-mcp-server",
        "env": {
          "METABASE_URL": "http://your-metabase-instance:3000",
          "METABASE_API_KEY": "mb_your_api_key_here"
        }
      }
    }
  }
}
```

Make sure the `metabase-mcp-server` binary is in your PATH, or provide the full path to the binary.

## Usage with MCP Remote (server deployment)

If you are running `metabase-mcp-server` on a remote server (e.g. a VPS, cloud instance, or within your internal network), you can expose it to local MCP clients using [mcp-remote](https://www.npmjs.com/package/mcp-remote). This is useful when you want a single centralized instance serving multiple users or machines.

### Server-side setup

1. Install the server on your remote machine and run it behind an SSE proxy. Since `metabase-mcp-server` uses stdio transport, you need a bridge that exposes it over HTTP with Server-Sent Events. Use [supergateway](https://github.com/supercorp-ai/supergateway) for this:

```bash
npx -y supergateway \
  --stdio "metabase-mcp-server --metabase-url http://your-metabase:3000 --api-key mb_your_key" \
  --port 8808
```

This starts an SSE MCP endpoint at `http://your-server:8808/sse`.

2. For production, run it as a systemd service or in Docker:

```bash
docker run --rm -p 8808:8808 \
  -e METABASE_URL=http://your-metabase:3000 \
  -e METABASE_API_KEY=mb_your_key \
  supergateway-metabase
```

Or create a `docker-compose.yml`:

```yaml
services:
  metabase-mcp:
    build: .
    environment:
      METABASE_URL: http://metabase:3000
      METABASE_API_KEY: mb_your_key
    # supergateway wraps stdio into HTTP
  supergateway:
    image: node:22-alpine
    command: npx -y supergateway --stdio "metabase-mcp-server" --port 8808 --host 0.0.0.0
    ports:
      - "8808:8808"
```

### Client-side: Claude Desktop

Use `mcp-remote` to connect Claude Desktop to the remote server:

```json
{
  "mcpServers": {
    "metabase": {
      "command": "npx",
      "args": [
        "-y", "mcp-remote",
        "http://your-server:8808/sse"
      ]
    }
  }
}
```

### Client-side: VS Code

```json
{
  "mcp": {
    "servers": {
      "metabase": {
        "command": "npx",
        "args": [
          "-y", "mcp-remote",
          "http://your-server:8808/sse"
        ]
      }
    }
  }
}
```

### Client-side: Claude Code (CLI)

```bash
claude mcp add metabase -- npx -y mcp-remote http://your-server:8808/sse
```

### Security considerations

- The HTTP endpoint does not include authentication by default. Place it behind a reverse proxy (nginx, Caddy, Traefik) with TLS and authentication if exposed to the internet.
- Use SSH tunneling as a simple alternative for private access:

```bash
# On your local machine, create an SSH tunnel
ssh -L 8808:localhost:8808 user@your-server

# Then connect to localhost in your MCP client config
# "http://localhost:8808/sse"
```

- For OAuth2/bearer token auth, `mcp-remote` supports passing headers:

```json
{
  "mcpServers": {
    "metabase": {
      "command": "npx",
      "args": [
        "-y", "mcp-remote",
        "https://your-server.example.com/sse",
        "--header", "Authorization: Bearer YOUR_TOKEN"
      ]
    }
  }
}
```

## Read-Only Safety

All SQL queries submitted through `execute_query` and `export_query_results` tools are validated before execution. The following SQL operations are blocked:

- INSERT, UPDATE, DELETE
- DROP, ALTER, CREATE, TRUNCATE
- GRANT, REVOKE
- EXEC, EXECUTE, MERGE, CALL

SQL comments are stripped before validation to prevent bypass attempts. Only SELECT queries and other read-only operations are permitted. MBQL queries are not subject to this restriction as they are constructed programmatically by Metabase.

## Development

### Prerequisites

- Go 1.25+
- golangci-lint v2 (for linting)

### Running Tests

```bash
go test -v -race ./...
```

### Linting

```bash
golangci-lint run
```

### Building

```bash
go build -o metabase-mcp-server ./cmd/metabase-mcp-server
```

## Project Structure

```
metabase-mcp-server/
  cmd/metabase-mcp-server/   -- Application entry point
  internal/
    config/                  -- Configuration parsing (flags + env vars)
    metabase/                -- Metabase API client library
    tools/                   -- MCP tool definitions and registration
  .github/workflows/         -- CI/CD pipelines
```

## License

MIT License. See [LICENSE](LICENSE) for details.
