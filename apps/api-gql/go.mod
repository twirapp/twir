module github.com/twirapp/twir/apps/api-gql

go 1.22.2

replace (
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/gomodels => ../../libs/gomodels
	github.com/satont/twir/libs/twitch => ../../libs/twitch
	github.com/twirapp/twir/libs/cache => ../../libs/cache
	github.com/twirapp/twir/libs/grpc => ../../libs/grpc
)

require (
	github.com/99designs/gqlgen v0.17.45
	github.com/alexedwards/scs/goredisstore v0.0.0-20240316134038-7e11d57e8885
	github.com/alexedwards/scs/v2 v2.8.0
	github.com/gin-contrib/cors v1.7.1
	github.com/gin-gonic/gin v1.9.1
	github.com/go-chi/httplog/v2 v2.0.9
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.0
	github.com/guregu/null v4.0.0+incompatible
	github.com/nicklaw5/helix/v2 v2.28.1
	github.com/redis/go-redis/v9 v9.5.1
	github.com/samber/lo v1.39.0
	github.com/satont/twir/libs/config v0.0.0-20240126231400-72985ccc25a5
	github.com/satont/twir/libs/gomodels v0.0.0-20240408183527-c93a8f2cd10e
	github.com/satont/twir/libs/twitch v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/cache v0.0.0-00010101000000-000000000000
	github.com/twirapp/twir/libs/grpc v0.0.0-20240126231400-72985ccc25a5
	github.com/vektah/gqlparser/v2 v2.5.11
	go.uber.org/fx v1.21.0
	gorm.io/driver/postgres v1.5.7
	gorm.io/gorm v1.25.9
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/bytedance/sonic v1.11.3 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.1 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-chi/chi/v5 v5.0.12 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.19.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/pprof v0.0.0-20240227163752-401108e1b7e7 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/imroc/req/v3 v3.42.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/ginkgo/v2 v2.15.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.0 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/quic-go v0.41.0 // indirect
	github.com/refraction-networking/utls v1.6.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/satont/twir/libs/types v0.0.0-20240408183527-c93a8f2cd10e // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sosodev/duration v1.2.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/urfave/cli/v2 v2.27.1 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.48.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/mock v0.4.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/arch v0.7.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225 // indirect
	golang.org/x/mod v0.16.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240304212257-790db918fca8 // indirect
	google.golang.org/grpc v1.62.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
