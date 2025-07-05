module github.com/satont/twir/apps/websockets

go 1.24.1

replace (
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/gomodels => ../../libs/gomodels
	github.com/satont/twir/libs/logger => ../../libs/logger
	github.com/satont/twir/libs/pubsub => ../../libs/pubsub
	github.com/satont/twir/libs/sentry => ../../libs/sentry
	github.com/satont/twir/libs/types => ../../libs/types
	github.com/twirapp/twir/libs/api => ../../libs/api
	github.com/twirapp/twir/libs/bus-core => ../../libs/bus-core
	github.com/twirapp/twir/libs/grpc => ../../libs/grpc
	github.com/twirapp/twir/libs/twitch => ../../libs/twitch
	github.com/twirapp/twir/libs/uptrace => ../../libs/uptrace
)

require (
	github.com/goccy/go-json v0.10.3
	github.com/google/uuid v1.6.0
	github.com/nicklaw5/helix/v2 v2.31.1
	github.com/olahol/melody v1.1.4
	github.com/prometheus/client_golang v1.22.0
	github.com/redis/go-redis/v9 v9.7.0
	github.com/samber/lo v1.49.1
	github.com/satont/twir/libs/config v0.0.0-20240126231400-72985ccc25a5
	github.com/satont/twir/libs/gomodels v0.0.0-20240225024146-742838c78cea
	github.com/satont/twir/libs/logger v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/twitch v0.0.0-20240129045902-98120ce2e77b
	github.com/satont/twir/libs/types v0.0.0-20241026093228-a5284254dcfe
	github.com/twirapp/twir/libs/api v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/bus-core v0.0.0-20240225024146-742838c78cea
	github.com/twirapp/twir/libs/grpc v0.0.0-20240126231400-72985ccc25a5
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.56.0
	go.uber.org/fx v1.23.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.69.4
	google.golang.org/protobuf v1.36.5
	gorm.io/gorm v1.25.12
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/getsentry/sentry-go v0.29.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nats-io/nats.go v1.37.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	github.com/samber/slog-common v0.17.1 // indirect
	github.com/samber/slog-multi v1.2.4 // indirect
	github.com/samber/slog-sentry/v2 v2.8.0 // indirect
	github.com/samber/slog-zerolog/v2 v2.7.0 // indirect
	github.com/satont/twir/libs/pubsub v0.0.0-00010101000000-000000000000 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/sdk v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	go.uber.org/dig v1.18.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250106144421-5f5ef82da422 // indirect
)
