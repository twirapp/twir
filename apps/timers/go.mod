module github.com/satont/tsuwari/apps/timers

go 1.20

require (
	github.com/go-co-op/gocron v1.17.1
	github.com/samber/do v1.6.0
	github.com/satont/tsuwari/libs/config v0.0.0
	github.com/satont/tsuwari/libs/gomodels v0.0.0
	github.com/satont/tsuwari/libs/grpc v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.24.0
	google.golang.org/grpc v1.52.3
	gorm.io/driver/postgres v1.4.7
	gorm.io/gorm v1.24.5
)

require (
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/jackc/pgx/v5 v5.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/protobuf v1.28.1
)

replace github.com/satont/tsuwari/libs/config => ../../libs/config

replace github.com/satont/tsuwari/libs/gomodels => ../../libs/gomodels

replace github.com/satont/tsuwari/libs/twitch => ../../libs/twitch

replace github.com/satont/tsuwari/libs/grpc => ../../libs/grpc
