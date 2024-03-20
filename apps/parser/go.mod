module github.com/satont/twir/apps/parser

go 1.21.5

replace (
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/gomodels => ../../libs/gomodels
	github.com/satont/twir/libs/types => ../../libs/types
	github.com/twirapp/twir/libs/grpc => ../../libs/grpc
	github.com/satont/twir/libs/twitch => ../../libs/twitch
	github.com/twirapp/twir/libs/integrations => ../../libs/integrations
	github.com/twirapp/twir/libs/bus-core => ../../libs/bus-core
)

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/getsentry/sentry-go v0.26.0
	github.com/goccy/go-json v0.10.2
	github.com/google/uuid v1.6.0
	github.com/guregu/null v4.0.0+incompatible
	github.com/hibiken/asynq v0.24.1
	github.com/imroc/req/v3 v3.42.3
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.9
	github.com/mazznoer/csscolorparser v0.1.3
	github.com/nats-io/nats.go v1.33.1
	github.com/nicklaw5/helix/v2 v2.25.3
	github.com/prometheus/client_golang v1.18.0
	github.com/redis/go-redis/v9 v9.4.0
	github.com/samber/lo v1.39.0
	github.com/satont/twir/libs/config v0.0.0-20240126231400-72985ccc25a5
	github.com/satont/twir/libs/gomodels v0.0.0-20240225024146-742838c78cea
	github.com/satont/twir/libs/twitch v0.0.0-20231203205548-e635accc6b72
	github.com/satont/twir/libs/types v0.0.0-20231203205548-e635accc6b72
	github.com/satori/go.uuid v1.2.0
	github.com/shkh/lastfm-go v0.0.0-20191215035245-89a801c244e0
	github.com/tidwall/gjson v1.17.0
	github.com/twirapp/twir/libs/bus-core v0.0.0-20240225024146-742838c78cea
	github.com/twirapp/twir/libs/grpc v0.0.0-20240126231400-72985ccc25a5
	github.com/twirapp/twir/libs/integrations v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/uptrace v0.0.0-20240225000519-a8b08004c468
	github.com/valyala/fasttemplate v1.2.2
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.48.0
	go.uber.org/zap v1.27.0
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225
	google.golang.org/grpc v1.62.0
	google.golang.org/protobuf v1.33.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.7
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/pprof v0.0.0-20240227163752-401108e1b7e7 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgx/v5 v5.5.2 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/onsi/ginkgo/v2 v2.15.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.46.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/quic-go v0.41.0 // indirect
	github.com/refraction-networking/utls v1.6.1 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/uptrace/uptrace-go v1.21.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.46.1 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/sdk v1.24.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.opentelemetry.io/proto/otlp v1.0.0 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/fx v1.20.1 // indirect
	go.uber.org/mock v0.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/mod v0.16.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/tools v0.19.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240304212257-790db918fca8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240304212257-790db918fca8 // indirect
)

