FROM golang:1.25-alpine AS builder
WORKDIR /workspace

# Copy module metadata first for better docker caching
COPY shop-identity-service/go.mod shop-identity-service/go.sum ./shop-identity-service/
COPY platform-events ./platform-events

WORKDIR /workspace/shop-identity-service
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY shop-identity-service/. .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o identity-service ./cmd/identity-service

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /workspace/shop-identity-service/identity-service /usr/local/bin/identity-service
COPY --from=builder /workspace/shop-identity-service/.env /app/.env
EXPOSE 8081 9091
ENTRYPOINT ["/usr/local/bin/identity-service"]
