module github.com/satont/twir/libs/migrations

go 1.21

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
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/docker/cli v25.0.4+incompatible // indirect
	github.com/docker/docker v25.0.4+incompatible // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/tools v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240304212257-790db918fca8 // indirect
	google.golang.org/grpc v1.62.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
