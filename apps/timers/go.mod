module github.com/satont/twir/apps/timers

go 1.21.0

replace (
	github.com/satont/twir/libs/config => ../../libs/config
	github.com/satont/twir/libs/grpc => ../../libs/grpc
	github.com/satont/twir/libs/logger => ../../libs/logger
)

require (
	github.com/getsentry/sentry-go v0.25.0
	github.com/satont/twir/libs/config v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/gomodels v0.0.0-20230816015129-6ac66282caa2
	github.com/satont/twir/libs/grpc v0.0.0-00010101000000-000000000000
	github.com/satont/twir/libs/logger v0.0.0-00010101000000-000000000000
	go.uber.org/fx v1.20.0
	google.golang.org/grpc v1.58.3
	google.golang.org/protobuf v1.31.0
	gorm.io/driver/postgres v1.5.3
	gorm.io/gorm v1.25.5
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/guregu/null v4.0.0+incompatible // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.17.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231012201019-e917dd12ba7a // indirect
)
