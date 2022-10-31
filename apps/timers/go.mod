module tsuwari/timers

go 1.19

require (
	github.com/go-co-op/gocron v1.17.1
	github.com/golang/protobuf v1.5.2
	github.com/satont/tsuwari/libs/nats v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.23.0
	gorm.io/driver/postgres v1.4.4
	gorm.io/gorm v1.24.0
	tsuwari/config v0.0.0
	tsuwari/models v0.0.0
	tsuwari/twitch v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/nats-io/jwt/v2 v2.3.0 // indirect
	github.com/satont/go-helix/v2 v2.7.14 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/time v0.0.0-20220922220347-f3bd1da661af // indirect
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
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/nats-io/nats.go v1.18.0
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/protobuf v1.28.1
)

replace github.com/satont/tsuwari/libs/nats => ../../libs/nats

replace tsuwari/models => ../../libs/gomodels

replace tsuwari/config => ../../libs/config

replace tsuwari/twitch => ../../libs/twitch
