module github.com/hawful70/shop-email-service

go 1.25.0

require (
	github.com/hawful70/platform-events v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/segmentio/kafka-go v0.4.46
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
)

replace github.com/hawful70/platform-events => ../platform-events
