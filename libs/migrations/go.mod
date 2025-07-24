module github.com/satont/twir/libs/migrations

go 1.24.1

replace github.cim/satont/twir/libs/config => ./../config

replace github.cim/satont/twir/libs/crypto => ./../crypto

replace github.com/twirapp/twir/libs/grpc => ./../grpc

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/pressly/goose/v3 v3.24.3
	github.com/satont/twir/libs/config v0.0.0-20250723210134-6e95e974f9e4
	github.com/satont/twir/libs/crypto v0.0.0-20250723210134-6e95e974f9e4
)

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.37.2 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/containerd/continuity v0.4.5 // indirect
	github.com/docker/cli v28.2.2+incompatible // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/jackc/pgx/v5 v5.7.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/opencontainers/runc v1.1.14 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/exp v0.0.0-20250606033433-dcc06ee1d476 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	golang.org/x/tools v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/grpc v1.72.2 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
