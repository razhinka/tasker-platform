# Multi-stage Dockerfile for Go Microservice
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/notifier ./cmd/notifier

# Minimal production image (~15 MB)
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/notifier /app/notifier

ENTRYPOINT ["/app/notifier"]
