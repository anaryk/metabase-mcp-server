FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /metabase-mcp-server ./cmd/metabase-mcp-server

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

COPY --from=builder /metabase-mcp-server /usr/local/bin/metabase-mcp-server

ENTRYPOINT ["metabase-mcp-server"]
