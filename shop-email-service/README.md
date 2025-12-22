# shop-email-service

Simple Go microservice that listens for `user_created` events on Kafka (via the Confluent REST Proxy) and sends welcome emails. Today it logs the outgoing email but the `Mailer` interface allows plugging in real providers (SMTP, SES, etc.).

## Run locally

```bash
cp .env.example .env
make run # (see note below if you create a Makefile)
```

The service expects Kafka REST proxy at `KAFKA_REST_URL` (default `http://localhost:8082`) and consumes topic `user_created` with group `email-service`.
