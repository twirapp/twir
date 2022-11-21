module github.com/satont/tsuwari/apps/api

go 1.19

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/getsentry/sentry-go v0.15.0
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.11.1
	github.com/go-redis/redis/v9 v9.0.0-rc.1
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/gofiber/contrib/fiberzap v0.0.0-20221120131424-7696474b1d10
	github.com/gofiber/fiber/v2 v2.40.0
	github.com/gofiber/storage/redis v0.0.0-20221120160944-6c0e70cefb0d
	github.com/golang/protobuf v1.5.2
	github.com/guregu/null v4.0.0+incompatible
	github.com/imroc/req/v3 v3.25.0
	github.com/nats-io/nats.go v1.20.0
	github.com/samber/lo v1.35.0
	github.com/satont/tsuwari/libs/integrations/spotify v0.0.0-20221119012024-2c6809904863
	github.com/satont/tsuwari/libs/nats v0.0.0-20221119012024-2c6809904863
	github.com/satont/tsuwari/libs/types v0.0.0-00010101000000-000000000000
	github.com/satori/go.uuid v1.2.0
	go.uber.org/zap v1.23.0
	google.golang.org/protobuf v1.28.1
	gorm.io/driver/postgres v1.4.5
	gorm.io/gorm v1.24.2
	tsuwari/config v0.0.0-00010101000000-000000000000
	tsuwari/models v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/pprof v0.0.0-20221118152302-e6195bd50e26 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/lucas-clemente/quic-go v0.31.0 // indirect
	github.com/marten-seemann/qpack v0.3.0 // indirect
	github.com/marten-seemann/qtls-go1-16 v0.1.5 // indirect
	github.com/marten-seemann/qtls-go1-17 v0.1.2 // indirect
	github.com/marten-seemann/qtls-go1-18 v0.1.3 // indirect
	github.com/marten-seemann/qtls-go1-19 v0.1.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/nats-io/jwt/v2 v2.3.0 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/onsi/ginkgo/v2 v2.5.1 // indirect
	github.com/rivo/uniseg v0.4.3 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/exp v0.0.0-20221114191408-850992195362 // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/tools v0.3.0 // indirect
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.15.12 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/satont/go-helix/v2 v2.7.21
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.41.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/crypto v0.3.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	tsuwari/twitch v0.0.0-00010101000000-000000000000
)

replace tsuwari/config => ../../libs/config

replace tsuwari/models => ../../libs/gomodels

replace tsuwari/twitch => ../../libs/twitch

replace github.com/satont/tsuwari/libs/nats => ../../libs/nats

replace github.com/satont/tsuwari/libs/integrations/spotify => ../../libs/integrations/spotify

replace github.com/satont/tsuwari/libs/types => ../../libs/types
