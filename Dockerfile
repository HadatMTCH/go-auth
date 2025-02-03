# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies and netcat for health check
RUN apk add --no-cache git make

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy dependency files first
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o /app/main main.go

# Final stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd

WORKDIR /app

# Copy goose binary from builder
COPY --from=builder /go/bin/goose /usr/local/bin/

# Copy other files
COPY --from=builder /app/main .
COPY --from=builder /app/.env.yml .
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations

COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["./main", "serve-http"]