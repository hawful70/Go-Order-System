FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o identity-service ./cmd/identity-service

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/identity-service /usr/local/bin/identity-service
COPY --from=builder /app/.env /app/.env
EXPOSE 8081 9091
ENTRYPOINT ["/usr/local/bin/identity-service"]