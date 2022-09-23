module tsuwari/timers

go 1.19

require (
	github.com/go-co-op/gocron v1.17.0
	github.com/go-redis/redis/v9 v9.0.0-beta.2
	github.com/golang/protobuf v1.5.0
	github.com/nicklaw5/helix v1.25.0
	github.com/satont/tsuwari/nats/bots v0.0.0-00010101000000-000000000000
	github.com/satont/tsuwari/nats/parser v0.0.0-00010101000000-000000000000
	github.com/satont/tsuwari/nats/timers v0.0.0-00010101000000-000000000000
	gorm.io/driver/postgres v1.3.10
	gorm.io/gorm v1.23.10
	tsuwari/config v0.0.0
	tsuwari/models v0.0.0
	tsuwari/twitch v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang-jwt/jwt v3.2.1+incompatible // indirect
	github.com/nats-io/nats-server/v2 v2.9.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.23.0 // indirect
)

require (
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/nats-io/nats.go v1.17.0
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	golang.org/x/crypto v0.0.0-20220919173607-35f4265a4bc0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.28.1
)

replace tsuwari/models => ../../libs/gomodels

replace tsuwari/config => ../../libs/config

replace github.com/satont/tsuwari/nats/timers => ../../libs/nats/timers

replace github.com/satont/tsuwari/nats/parser => ../../libs/nats/parser

replace github.com/satont/tsuwari/nats/bots => ../../libs/nats/bots

replace tsuwari/twitch => ../../libs/twitch
