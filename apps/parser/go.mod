module github.com/satont/twir/apps/parser

go 1.24.1

replace (
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/crypto => ../../libs/crypto
	github.com/satont/twir/libs/gomodels => ../../libs/gomodels
	github.com/satont/twir/libs/logger => ../../libs/logger
	github.com/satont/twir/libs/pubsub => ../../libs/pubsub
	github.com/satont/twir/libs/sentry => ../../libs/sentry
	github.com/satont/twir/libs/twitch => ../../libs/twitch
	github.com/satont/twir/libs/types => ../../libs/types
	github.com/satont/twir/libs/utils => ../../libs/utils
	github.com/twirapp/twir/libs/api => ../../libs/api
	github.com/twirapp/twir/libs/baseapp => ../../libs/baseapp
	github.com/twirapp/twir/libs/bus-core => ../../libs/bus-core
	github.com/twirapp/twir/libs/cache => ../../libs/cache
	github.com/twirapp/twir/libs/grpc => ../../libs/grpc
	github.com/twirapp/twir/libs/integrations => ../../libs/integrations
	github.com/twirapp/twir/libs/repositories => ../../libs/repositories
	github.com/twirapp/twir/libs/uptrace => ../../libs/uptrace
)

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2 v2.0.0
	github.com/avito-tech/go-transaction-manager/trm/v2 v2.0.0
	github.com/exaring/otelpgx v0.9.3
	github.com/getsentry/sentry-go v0.34.1
	github.com/go-redis/redis_rate/v10 v10.0.1
	github.com/go-redsync/redsync/v4 v4.13.0
	github.com/google/uuid v1.6.0
	github.com/guregu/null v4.0.0+incompatible
	github.com/imroc/req/v3 v3.54.0
	github.com/jackc/pgx/v5 v5.7.5
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
	github.com/matoous/go-nanoid/v2 v2.1.0
	github.com/mazznoer/csscolorparser v0.1.6
	github.com/nats-io/nats.go v1.43.0
	github.com/nicklaw5/helix/v2 v2.31.1
	github.com/prometheus/client_golang v1.22.0
	github.com/redis/go-redis/v9 v9.11.0
	github.com/samber/lo v1.51.0
	github.com/satont/twir/libs/config v0.0.0-20250723210134-6e95e974f9e4
	github.com/satont/twir/libs/gomodels v0.0.0-20250723210134-6e95e974f9e4
	github.com/satont/twir/libs/twitch v0.0.0-20250723210134-6e95e974f9e4
	github.com/satont/twir/libs/types v0.0.0-20250723210134-6e95e974f9e4
	github.com/satori/go.uuid v1.2.0
	github.com/shkh/lastfm-go v0.0.0-20191215035245-89a801c244e0
	github.com/tidwall/gjson v1.18.0
	github.com/twirapp/twir/libs/baseapp v0.0.0-20250723210134-6e95e974f9e4
	github.com/twirapp/twir/libs/bus-core v0.0.0-20250723210134-6e95e974f9e4
	github.com/twirapp/twir/libs/cache v0.0.0-20250723210134-6e95e974f9e4
	github.com/twirapp/twir/libs/grpc v0.0.0-20250723210134-6e95e974f9e4
	github.com/twirapp/twir/libs/integrations v0.0.0-20250723210134-6e95e974f9e4
	github.com/twirapp/twir/libs/redis_keys v0.0.0-20250723210134-6e95e974f9e4
	github.com/twirapp/twir/libs/repositories v0.0.0-20250723210134-6e95e974f9e4
	github.com/twirapp/twir/libs/uptrace v0.0.0-20250723210134-6e95e974f9e4
	github.com/valyala/fasttemplate v1.2.2
	github.com/xhit/go-str2duration/v2 v2.1.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.62.0
	go.uber.org/zap v1.27.0
	golang.org/x/exp v0.0.0-20250718183923-645b1fa84792
	golang.org/x/sync v0.16.0
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.6
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.30.1
)

require (
	github.com/ClickHouse/ch-go v0.67.0 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.39.0 // indirect
	github.com/Khan/genqlient v0.8.1 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudflare/circl v1.6.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/google/pprof v0.0.0-20241210010833-40e02aabc2ad // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/icholy/digest v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/oklog/ulid/v2 v2.1.1 // indirect
	github.com/onsi/ginkgo/v2 v2.22.2 // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.65.0 // indirect
	github.com/prometheus/procfs v0.17.0 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.54.0 // indirect
	github.com/redis/go-redis/extra/rediscmd/v9 v9.11.0 // indirect
	github.com/redis/go-redis/extra/redisotel/v9 v9.11.0 // indirect
	github.com/refraction-networking/utls v1.8.0 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/samber/slog-common v0.19.0 // indirect
	github.com/samber/slog-multi v1.4.1 // indirect
	github.com/samber/slog-sentry/v2 v2.9.3 // indirect
	github.com/samber/slog-zerolog/v2 v2.7.3 // indirect
	github.com/satont/twir/libs/logger v0.0.0-20250723210134-6e95e974f9e4 // indirect
	github.com/satont/twir/libs/pubsub v0.0.0-20250723210134-6e95e974f9e4 // indirect
	github.com/satont/twir/libs/sentry v0.0.0-20250723210134-6e95e974f9e4 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelgorm v0.3.2 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelsql v0.3.2 // indirect
	github.com/uptrace/uptrace-go v1.37.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/vektah/gqlparser/v2 v2.5.30 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.62.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.62.0 // indirect
	go.opentelemetry.io/otel v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.37.0 // indirect
	go.opentelemetry.io/otel/log v0.13.0 // indirect
	go.opentelemetry.io/otel/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.13.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.37.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/fx v1.24.0 // indirect
	go.uber.org/mock v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/mod v0.26.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	golang.org/x/tools v0.35.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250721164621-a45f3dfb1074 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250721164621-a45f3dfb1074 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
