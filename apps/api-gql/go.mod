module github.com/twirapp/twir/apps/api-gql

go 1.23.0

toolchain go1.23.1

replace (
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/gomodels => ../../libs/gomodels
	github.com/satont/twir/libs/logger => ../../libs/logger
	github.com/satont/twir/libs/sentry => ../../libs/sentry
	github.com/satont/twir/libs/twitch => ../../libs/twitch
	github.com/twirapp/twir/libs/baseapp => ../../libs/baseapp
	github.com/twirapp/twir/libs/bus-core => ../../libs/bus-core
	github.com/twirapp/twir/libs/cache => ../../libs/cache
	github.com/twirapp/twir/libs/grpc => ../../libs/grpc
	github.com/twirapp/twir/libs/redis_keys => ../../libs/redis_keys
	github.com/twirapp/twir/libs/uptrace => ../../libs/uptrace
	github.com/twirapp/twir/libs/pubsub => ../../libs/pubsub
)

require (
	github.com/99designs/gqlgen v0.17.45
	github.com/Masterminds/squirrel v1.5.4
	github.com/alexedwards/scs/goredisstore v0.0.0-20240316134038-7e11d57e8885
	github.com/alexedwards/scs/v2 v2.8.0
	github.com/evanphx/json-patch/v5 v5.9.0
	github.com/gin-contrib/cors v1.7.1
	github.com/gin-gonic/gin v1.9.1
	github.com/go-chi/httplog/v2 v2.0.9
	github.com/go-playground/validator/v10 v10.19.0
	github.com/goccy/go-json v0.10.2
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.1
	github.com/guregu/null v4.0.0+incompatible
	github.com/kr/pretty v0.3.1
	github.com/lib/pq v1.10.9
	github.com/minio/minio-go/v7 v7.0.69
	github.com/nats-io/nats.go v1.36.0
	github.com/nicklaw5/helix/v2 v2.28.1
	github.com/raitonoberu/ytsearch v0.2.0
	github.com/ravilushqa/otelgqlgen v0.15.0
	github.com/redis/go-redis/v9 v9.5.3
	github.com/samber/lo v1.39.0
	github.com/samber/slog-gin v1.11.0
	github.com/satont/twir/libs/config v0.0.0-20240126231400-72985ccc25a5
	github.com/satont/twir/libs/crypto v0.0.0-20240923063656-34c54be7623e
	github.com/satont/twir/libs/gomodels v0.0.0-20240208154120-fc098a9e20a2
	github.com/satont/twir/libs/logger v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/pubsub v0.0.0-20240923063656-34c54be7623e
	github.com/satont/twir/libs/twitch v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/types v0.0.0-20240411125706-789748d8a577
	github.com/twirapp/twir/libs/baseapp v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/bus-core v0.0.0-20240225024146-742838c78cea
	github.com/twirapp/twir/libs/cache v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/grpc v0.0.0-20240126231400-72985ccc25a5
	github.com/twirapp/twir/libs/integrations v0.0.0-20240923063656-34c54be7623e
	github.com/twirapp/twir/libs/redis_keys v0.0.0-00010101000000-000000000000
	github.com/vektah/gqlparser/v2 v2.5.11
	github.com/vikstrous/dataloadgen v0.0.6
	go.uber.org/fx v1.21.0
	gorm.io/gorm v1.25.9
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/bytedance/sonic v1.11.3 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.1 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/getsentry/sentry-go v0.26.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-chi/chi/v5 v5.0.10 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/pprof v0.0.0-20240227163752-401108e1b7e7 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/imroc/req/v3 v3.42.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgx/v5 v5.5.2 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/onsi/ginkgo/v2 v2.15.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/quic-go v0.41.0 // indirect
	github.com/refraction-networking/utls v1.6.1 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/rs/xid v1.5.0 // indirect
	github.com/rs/zerolog v1.31.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/samber/slog-common v0.15.0 // indirect
	github.com/samber/slog-multi v1.0.2 // indirect
	github.com/samber/slog-sentry/v2 v2.4.0 // indirect
	github.com/samber/slog-zerolog/v2 v2.2.0 // indirect
	github.com/satont/twir/libs/sentry v0.0.0-00010101000000-000000000000 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/shkh/lastfm-go v0.0.0-20191215035245-89a801c244e0 // indirect
	github.com/sosodev/duration v1.2.0 // indirect
	github.com/twirapp/twir/libs/uptrace v0.0.0-00010101000000-000000000000 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/uptrace/uptrace-go v1.21.0 // indirect
	github.com/urfave/cli/v2 v2.27.1 // indirect
	github.com/xrash/smetrics v0.0.0-20231213231151-1d8dd44e695e // indirect
	go.opentelemetry.io/contrib v1.21.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.48.0 // indirect
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
	go.uber.org/mock v0.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/arch v0.7.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225 // indirect
	golang.org/x/mod v0.16.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.19.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240304212257-790db918fca8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240304212257-790db918fca8 // indirect
	google.golang.org/grpc v1.62.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.5.7 // indirect
)
