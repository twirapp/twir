module github.com/satont/twir/apps/bots

go 1.23.0

require (
	github.com/aidenwallis/go-ratelimiting v0.0.5
	github.com/go-redsync/redsync/v4 v4.11.0
	github.com/goccy/go-json v0.10.2
	github.com/google/uuid v1.6.0
	github.com/hibiken/asynq v0.24.1
	github.com/imroc/req/v3 v3.42.3
	github.com/lib/pq v1.10.9
	github.com/nicklaw5/helix/v2 v2.25.3
	github.com/prometheus/client_golang v1.18.0
	github.com/redis/go-redis/v9 v9.5.3
	github.com/samber/lo v1.39.0
	github.com/satont/twir/libs/config v0.0.0-20240126231400-72985ccc25a5
	github.com/satont/twir/libs/gomodels v0.0.0-20240225024146-742838c78cea
	github.com/satont/twir/libs/logger v0.0.0-20240126231400-72985ccc25a5
	github.com/satont/twir/libs/twitch v0.0.0-20240126231400-72985ccc25a5
	github.com/satont/twir/libs/types v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/utils v0.0.0-20240126231400-72985ccc25a5
	github.com/twirapp/twir/libs/bus-core v0.0.0-20240225024146-742838c78cea
	github.com/twirapp/twir/libs/grpc v0.0.0-20240126231400-72985ccc25a5
	github.com/twirapp/twir/libs/redis_keys v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/trace v1.24.0
	go.uber.org/fx v1.21.0
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.7.0
	gorm.io/driver/postgres v1.5.7
	gorm.io/gorm v1.25.9
)

require golang.org/x/net v0.22.0 // indirect

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/getsentry/sentry-go v0.26.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/pprof v0.0.0-20240227163752-401108e1b7e7 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
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
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/nats-io/nats.go v1.36.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/onsi/ginkgo/v2 v2.15.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.46.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/quic-go v0.41.0 // indirect
	github.com/refraction-networking/utls v1.6.1 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rs/zerolog v1.31.0 // indirect
	github.com/samber/slog-common v0.15.0 // indirect
	github.com/samber/slog-multi v1.0.2 // indirect
	github.com/samber/slog-sentry/v2 v2.4.0 // indirect
	github.com/samber/slog-zerolog/v2 v2.2.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.48.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/mock v0.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225 // indirect
	golang.org/x/mod v0.16.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/tools v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240304212257-790db918fca8 // indirect
	google.golang.org/grpc v1.62.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace (
	github.com/nicklaw5/helix/v2 => ../../libs/helix
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/gomodels => ../../libs/gomodels
	github.com/satont/twir/libs/logger => ../../libs/logger
	github.com/satont/twir/libs/twitch => ../../libs/twitch
	github.com/satont/twir/libs/types => ../../libs/types
	github.com/satont/twir/libs/utils => ../../libs/utils
	github.com/twirapp/twir/libs/baseapp => ../../libs/baseapp
	github.com/twirapp/twir/libs/bus-core => ../../libs/bus-core
	github.com/twirapp/twir/libs/grpc => ../../libs/grpc
	github.com/twirapp/twir/libs/redis_keys => ../../libs/redis_keys
	github.com/twirapp/twir/libs/uptrace => ../../libs/uptrace
)
