module github.com/satont/twir/libs/migrations

go 1.24.1

replace github.cim/satont/twir/libs/config => ./../config

replace github.cim/satont/twir/libs/crypto => ./../crypto

replace github.com/twirapp/twir/libs/grpc => ./../grpc

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/pressly/goose/v3 v3.18.0
	github.com/satont/twir/libs/config v0.0.0-20231203205548-e635accc6b72
	github.com/satont/twir/libs/crypto v0.0.0-20231203205548-e635accc6b72
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/docker/cli v25.0.4+incompatible // indirect
	github.com/docker/docker v25.0.4+incompatible // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.1 // indirect
	github.com/jackc/pgx/v5 v5.7.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.opentelemetry.io/otel/trace v1.33.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/tools v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241104194629-dd2ea8efbc28 // indirect
	google.golang.org/grpc v1.67.1 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)
