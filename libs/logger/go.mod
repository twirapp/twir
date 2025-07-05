module github.com/satont/twir/libs/logger

go 1.24.1

replace (
	github.com/satont/twir/libs/config => ../config
	github.com/satont/twir/libs/gomodels => ../gomodels
	github.com/satont/twir/libs/pubsub => ../pubsub
	github.com/twirapp/twir/libs/bus-core => ../bus-core
	github.com/twirapp/twir/libs/repositories => ../repositories
)

require (
	github.com/getsentry/sentry-go v0.29.1
	github.com/goccy/go-json v0.10.3
	github.com/google/uuid v1.6.0
	github.com/guregu/null v4.0.0+incompatible
	github.com/rs/zerolog v1.33.0
	github.com/samber/slog-common v0.17.1
	github.com/samber/slog-multi v1.2.4
	github.com/samber/slog-sentry/v2 v2.8.0
	github.com/samber/slog-zerolog/v2 v2.7.0
	github.com/satont/twir/libs/pubsub v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/bus-core v0.0.0-20240225024146-742838c78cea
	github.com/twirapp/twir/libs/repositories v0.0.0-00010101000000-000000000000
	go.uber.org/fx v1.23.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/nats-io/nats.go v1.37.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/samber/lo v1.49.1 // indirect
	github.com/satont/twir/libs/config v0.0.0-00010101000000-000000000000 // indirect
	go.uber.org/dig v1.18.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
)
