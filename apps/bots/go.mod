module github.com/satont/twir/apps/bots

go 1.21

require (
	github.com/aidenwallis/go-ratelimiting v0.0.1
	github.com/gempir/go-twitch-irc/v3 v3.2.0
	github.com/getsentry/sentry-go v0.25.0
	github.com/goccy/go-json v0.10.2
	github.com/google/uuid v1.3.1
	github.com/imroc/req/v3 v3.42.1
	github.com/lib/pq v1.10.9
	github.com/nicklaw5/helix/v2 v2.25.2
	github.com/prometheus/client_golang v1.17.0
	github.com/redis/go-redis/v9 v9.2.1
	github.com/samber/lo v1.38.1
	github.com/satont/twir/libs/config v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/gomodels v0.0.0-20221114143619-e5e207524b96
	github.com/satont/twir/libs/gopool v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/grpc v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/logger v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/pubsub v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/twitch v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/types v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/utils v0.0.0-00010101000000-000000000000
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.8.4
	go.uber.org/fx v1.20.1
	go.uber.org/zap v1.26.0
	google.golang.org/grpc v1.58.3
	google.golang.org/protobuf v1.31.0
	gorm.io/driver/postgres v1.5.3
	gorm.io/gorm v1.25.5
)

require golang.org/x/net v0.17.0 // indirect

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gaukas/godicttls v0.0.4 // indirect
	github.com/go-redis/redis/v9 v9.0.0-rc.2 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/pprof v0.0.0-20230901174712-0191c66da455 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/onsi/ginkgo/v2 v2.12.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.4.1-0.20230718164431-9a2bf3000d16 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/qtls-go1-20 v0.3.3 // indirect
	github.com/quic-go/quic-go v0.38.1 // indirect
	github.com/refraction-networking/utls v1.5.3 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.17.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/mod v0.13.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/tools v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231012201019-e917dd12ba7a // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/gomodels => ../../libs/gomodels
	github.com/satont/twir/libs/gopool => ../../libs/gopool
	github.com/satont/twir/libs/grpc => ../../libs/grpc
	github.com/satont/twir/libs/logger => ../../libs/logger
	github.com/satont/twir/libs/pubsub => ../../libs/pubsub
	github.com/satont/twir/libs/twitch => ../../libs/twitch
	github.com/satont/twir/libs/types => ../../libs/types
	github.com/satont/twir/libs/utils => ../../libs/utils
)
