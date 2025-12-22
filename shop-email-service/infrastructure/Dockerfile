FROM golang:1.25-alpine AS builder
WORKDIR /workspace

COPY shop-email-service/go.mod shop-email-service/go.sum ./shop-email-service/
COPY platform-events ./platform-events

WORKDIR /workspace/shop-email-service
RUN go mod download

COPY shop-email-service/. .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o email-service ./cmd/email-service

FROM gcr.io/distroless/base-debian12
COPY --from=builder /workspace/shop-email-service/email-service /usr/local/bin/email-service
ENTRYPOINT ["/usr/local/bin/email-service"]
