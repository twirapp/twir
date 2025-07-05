module github.com/satont/twir/apps/timers

go 1.24.1

replace (
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/logger => ../../libs/logger
	github.com/satont/twir/libs/pubsub => ../../libs/pubsub
	github.com/twirapp/twir/libs/bus-core => ../../libs/bus-core
	github.com/twirapp/twir/libs/grpc => ../../libs/grpc
	github.com/twirapp/twir/libs/redis_keys => ../../libs/redis_keys
	github.com/twirapp/twir/libs/uptrace => ../../libs/uptrace
)

require (
	github.com/google/uuid v1.6.0
	github.com/redis/go-redis/v9 v9.7.0
	github.com/satont/twir/libs/config v0.0.0-20231203205548-e635accc6b72
	github.com/satont/twir/libs/gomodels v0.0.0-20231203205548-e635accc6b72
	github.com/satont/twir/libs/logger v0.0.0-20231203205548-e635accc6b72
	github.com/twirapp/twir/libs/bus-core v0.0.0-20240225024146-742838c78cea
	github.com/twirapp/twir/libs/grpc v0.0.0-20231203205548-e635accc6b72
	github.com/twirapp/twir/libs/redis_keys v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/uptrace v0.0.0-00010101000000-000000000000
	go.temporal.io/sdk v1.25.1
	go.uber.org/fx v1.23.0
	gorm.io/gorm v1.25.12
)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/getsentry/sentry-go v0.29.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gogo/status v1.1.1 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.23.0 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/nats-io/nats.go v1.37.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	github.com/samber/lo v1.49.1 // indirect
	github.com/samber/slog-common v0.17.1 // indirect
	github.com/samber/slog-multi v1.2.4 // indirect
	github.com/samber/slog-sentry/v2 v2.8.0 // indirect
	github.com/samber/slog-zerolog/v2 v2.7.0 // indirect
	github.com/satont/twir/libs/pubsub v0.0.0-00010101000000-000000000000 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/uptrace/uptrace-go v1.31.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.56.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.56.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.31.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.31.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.31.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.31.0 // indirect
	go.opentelemetry.io/otel/log v0.7.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/sdk v1.36.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.7.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.31.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	go.temporal.io/api v1.24.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.18.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	golang.org/x/time v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20240401170217-c3f982113cda // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250106144421-5f5ef82da422 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250106144421-5f5ef82da422 // indirect
	google.golang.org/grpc v1.69.4 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
