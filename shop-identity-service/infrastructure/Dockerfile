FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o identity-service ./cmd/identity-service

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/identity-service /usr/local/bin/identity-service
COPY --from=builder /app/.env /app/.env

EXPOSE 8081 9091

ENTRYPOINT ["/usr/local/bin/identity-service"]
