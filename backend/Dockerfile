# --- Build stage ---
FROM golang:1.26-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /api ./cmd/api

# --- Runtime stage (minimal, non-root) ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates && adduser -D -u 10001 appuser
USER appuser
COPY --from=builder /api /api

EXPOSE 8080
ENTRYPOINT ["/api"]
