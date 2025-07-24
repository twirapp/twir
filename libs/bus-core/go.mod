module github.com/twirapp/twir/libs/bus-core

go 1.24.1

replace github.com/satont/twir/libs/config => ../config

require (
	github.com/goccy/go-json v0.10.5
	github.com/google/uuid v1.6.0
	github.com/nats-io/nats.go v1.43.0
	github.com/satont/twir/libs/config v0.0.0-20250723210134-6e95e974f9e4
	github.com/twirapp/twir/libs/repositories v0.0.0-20250723210134-6e95e974f9e4
	go.opentelemetry.io/otel v1.37.0
	go.opentelemetry.io/otel/trace v1.37.0
)

require (
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.37.0 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
)
