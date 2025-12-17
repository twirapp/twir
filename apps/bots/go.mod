module github.com/twirapp/twir/apps/bots

go 1.25.4

require (
	github.com/aidenwallis/go-ratelimiting v0.0.5
	github.com/alitto/pond/v2 v2.5.0
	github.com/dlclark/regexp2 v1.11.5
	github.com/goccy/go-json v0.10.5
	github.com/google/uuid v1.6.0
	github.com/hibiken/asynq v0.25.1
	github.com/lib/pq v1.10.9
	github.com/nicklaw5/helix/v2 v2.32.0
	github.com/prometheus/client_golang v1.23.2
	github.com/redis/go-redis/v9 v9.17.1
	github.com/samber/lo v1.52.0
	github.com/twirapp/batch-processor v0.0.1
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/trace v1.38.0
	go.uber.org/fx v1.24.0
	golang.org/x/sync v0.18.0
	gorm.io/gorm v1.31.1
)

require (
	cloud.google.com/go/translate v1.12.7
	github.com/lkretschmer/deepl-go v0.3.0
	github.com/raitonoberu/ytsearch v0.2.0
	github.com/twirapp/kv v0.5.1
	github.com/twirapp/twir/libs/baseapp v0.0.0-20251217140736-48670d138d86
	github.com/twirapp/twir/libs/bus-core v0.0.0-20251217140736-48670d138d86
	github.com/twirapp/twir/libs/cache v0.0.0-20251217140736-48670d138d86
	github.com/twirapp/twir/libs/config v0.0.0-20251217140736-48670d138d86
	github.com/twirapp/twir/libs/entities v0.0.0-20251217140736-48670d138d86
	github.com/twirapp/twir/libs/gomodels v0.0.0-20251217140736-48670d138d86
	github.com/twirapp/twir/libs/grpc v0.0.0-20251217140736-48670d138d86
	github.com/twirapp/twir/libs/logger v0.0.0-20251217140736-48670d138d86
	github.com/twirapp/twir/libs/redis_keys v0.0.0-20251217140736-48670d138d86
	github.com/twirapp/twir/libs/repositories v0.0.0-20251217140736-48670d138d86
	github.com/twirapp/twir/libs/twitch v0.0.0-20251217140736-48670d138d86
	github.com/twirapp/twir/libs/types v0.0.0-20251217140736-48670d138d86
	github.com/twirapp/twir/libs/utils v0.0.0-20251217140736-48670d138d86
	google.golang.org/api v0.257.0
)

require (
	cloud.google.com/go v0.123.0 // indirect
	cloud.google.com/go/auth v0.17.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	github.com/ClickHouse/ch-go v0.69.0 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.41.0 // indirect
	github.com/Masterminds/squirrel v1.5.4 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2 v2.0.2 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/exaring/otelpgx v0.9.3 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/getsentry/sentry-go v0.40.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/golang/mock v1.7.0-rc.1 // indirect
	github.com/gomodule/redigo v1.9.3 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.7 // indirect
	github.com/googleapis/gax-go/v2 v2.15.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.6 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/maypok86/otter/v2 v2.2.1 // indirect
	github.com/nats-io/nats.go v1.47.0 // indirect
	github.com/nats-io/nkeys v0.4.12 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/oklog/ulid/v2 v2.1.1 // indirect
	github.com/paulmach/orb v0.12.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/redis/go-redis/extra/rediscmd/v9 v9.17.1 // indirect
	github.com/redis/go-redis/extra/redisotel/v9 v9.17.1 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/samber/slog-common v0.19.0 // indirect
	github.com/samber/slog-multi v1.6.0 // indirect
	github.com/samber/slog-sentry/v2 v2.10.1 // indirect
	github.com/samber/slog-zerolog/v2 v2.9.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/twirapp/twir/libs/audit v0.0.0-20251217140736-48670d138d86 // indirect
	github.com/twirapp/twir/libs/pubsub v0.0.0-20251217140736-48670d138d86 // indirect
	github.com/twirapp/twir/libs/sentry v0.0.0-20251217140736-48670d138d86 // indirect
	github.com/twirapp/twir/libs/uptrace v0.0.0-20251217140736-48670d138d86 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelgorm v0.3.2 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelsql v0.3.2 // indirect
	github.com/uptrace/uptrace-go v1.38.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.63.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.63.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.63.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.14.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.38.0 // indirect
	go.opentelemetry.io/otel/log v0.14.0 // indirect
	go.opentelemetry.io/otel/sdk v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.14.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.38.0 // indirect
	go.opentelemetry.io/proto/otlp v1.9.0 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/oauth2 v0.33.0 // indirect
	google.golang.org/genproto v0.0.0-20250721164621-a45f3dfb1074 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251124214823-79d6a2a48846 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251124214823-79d6a2a48846 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	gorm.io/driver/postgres v1.6.0 // indirect
)

require (
	github.com/avito-tech/go-transaction-manager/trm/v2 v2.0.2
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.4 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/spf13/cast v1.9.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0
	golang.org/x/time v0.14.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
