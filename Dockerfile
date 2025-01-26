# Build stage
FROM golang:1.23.3-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy dependency files first
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

# Build binary using root main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main main.go

# Final stage
FROM alpine:3.18

WORKDIR /app

# Copy built binary and configs
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations

# Create log directory
RUN mkdir -p /var/log/app
VOLUME /var/log/app

EXPOSE 8080

# Run with serve-http command
CMD ["./main", "serve-http"]