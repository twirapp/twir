module github.com/twirapp/twir/apps/api-gql

go 1.25.4

require (
	github.com/99designs/gqlgen v0.17.84
	github.com/Masterminds/squirrel v1.5.4
	github.com/aidenwallis/go-ratelimiting v0.0.5
	github.com/alexedwards/scs/goredisstore v0.0.0-20251002162104-209de6e426de
	github.com/alexedwards/scs/v2 v2.9.0
	github.com/avito-tech/go-transaction-manager/trm/v2 v2.0.2
	github.com/danielgtaylor/huma/v2 v2.34.1
	github.com/evanphx/json-patch/v5 v5.9.11
	github.com/gin-contrib/cors v1.7.6
	github.com/gin-gonic/gin v1.11.0
	github.com/go-playground/validator/v10 v10.28.0
	github.com/goccy/go-json v0.10.5
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/guregu/null v4.0.0+incompatible
	github.com/lib/pq v1.10.9
	github.com/matoous/go-nanoid/v2 v2.1.0
	github.com/minio/minio-go/v7 v7.0.97
	github.com/nats-io/nats.go v1.47.0
	github.com/nicklaw5/helix/v2 v2.32.0
	github.com/oklog/ulid/v2 v2.1.1
	github.com/raitonoberu/ytsearch v0.2.0
	github.com/ravilushqa/otelgqlgen v0.19.0
	github.com/redis/go-redis/v9 v9.17.1
	github.com/samber/lo v1.52.0
	github.com/samber/slog-gin v1.18.0
	github.com/shkh/lastfm-go v0.0.0-20191215035245-89a801c244e0
	github.com/twirapp/kv v0.5.1
	github.com/twirapp/twir/apps/parser v0.0.0-20251127124349-67ad7fa0003f
	github.com/twirapp/twir/libs/audit v0.0.0-20260102031833-cfb924c9b7eb
	github.com/twirapp/twir/libs/baseapp v0.0.0-20251127124349-67ad7fa0003f
	github.com/twirapp/twir/libs/bus-core v0.0.0-20260102031833-cfb924c9b7eb
	github.com/twirapp/twir/libs/cache v0.0.0-20260102031833-cfb924c9b7eb
	github.com/twirapp/twir/libs/config v0.0.0-20251201102513-6f706c3cc7e1
	github.com/twirapp/twir/libs/crypto v0.0.0-20251127124349-67ad7fa0003f
	github.com/twirapp/twir/libs/entities v0.0.0-20260102031833-cfb924c9b7eb
	github.com/twirapp/twir/libs/gomodels v0.0.0-20251127124349-67ad7fa0003f
	github.com/twirapp/twir/libs/integrations v0.0.0-20251127124349-67ad7fa0003f
	github.com/twirapp/twir/libs/logger v0.0.0-20251127124349-67ad7fa0003f
	github.com/twirapp/twir/libs/pubsub v0.0.0-20251127124349-67ad7fa0003f
	github.com/twirapp/twir/libs/redis_keys v0.0.0-20251127124349-67ad7fa0003f
	github.com/twirapp/twir/libs/repositories v0.0.0-20260102031833-cfb924c9b7eb
	github.com/twirapp/twir/libs/twitch v0.0.0-20251127124349-67ad7fa0003f
	github.com/twirapp/twir/libs/types v0.0.0-20251127124349-67ad7fa0003f
	github.com/twirapp/twir/libs/utils v0.0.0-20251127124349-67ad7fa0003f
	github.com/vektah/gqlparser/v2 v2.5.31
	github.com/vikstrous/dataloadgen v0.0.10
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.63.0
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/trace v1.38.0
	go.uber.org/fx v1.24.0
	golang.org/x/sync v0.18.0
	gorm.io/gorm v1.31.1
)

require (
	github.com/ClickHouse/ch-go v0.69.0 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.41.0 // indirect
	github.com/Khan/genqlient v0.8.1 // indirect
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2 v2.0.2 // indirect
	github.com/bytedance/gopkg v0.1.3 // indirect
	github.com/bytedance/sonic v1.14.2 // indirect
	github.com/bytedance/sonic/loader v0.4.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/exaring/otelpgx v0.9.3 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.11 // indirect
	github.com/getsentry/sentry-go v0.40.0 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-redis/redis_rate/v10 v10.0.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/goccy/go-yaml v1.18.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/gomodule/redigo v1.9.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/icholy/digest v1.1.0 // indirect
	github.com/imroc/req/v3 v3.54.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.6 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/klauspost/crc32 v1.3.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/maypok86/otter/v2 v2.2.1 // indirect
	github.com/minio/crc64nvme v1.1.1 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nats-io/nkeys v0.4.12 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/paulmach/orb v0.12.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.57.1 // indirect
	github.com/redis/go-redis/extra/rediscmd/v9 v9.17.1 // indirect
	github.com/redis/go-redis/extra/redisotel/v9 v9.17.1 // indirect
	github.com/refraction-networking/utls v1.8.1 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/samber/slog-common v0.19.0 // indirect
	github.com/samber/slog-multi v1.6.0 // indirect
	github.com/samber/slog-sentry/v2 v2.10.1 // indirect
	github.com/samber/slog-zerolog/v2 v2.9.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/tinylib/msgp v1.5.0 // indirect
	github.com/twirapp/twir/libs/sentry v0.0.0-20251127124349-67ad7fa0003f // indirect
	github.com/twirapp/twir/libs/uptrace v0.0.0-20251127124349-67ad7fa0003f // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.1 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelgorm v0.3.2 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelsql v0.3.2 // indirect
	github.com/uptrace/uptrace-go v1.38.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib v1.38.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.63.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.63.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.14.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.38.0 // indirect
	go.opentelemetry.io/otel/log v0.14.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.14.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.38.0 // indirect
	go.opentelemetry.io/proto/otlp v1.9.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/mock v0.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/arch v0.23.0 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251124214823-79d6a2a48846 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251124214823-79d6a2a48846 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.6.0 // indirect
)
